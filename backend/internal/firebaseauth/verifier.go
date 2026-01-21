package firebaseauth

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const certsURL = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrUnauthorized = errors.New("unauthorized")
)

type Verifier struct {
	projectID string
	client    *http.Client

	mu        sync.Mutex
	certs     map[string]*rsa.PublicKey
	expiresAt time.Time
}

type AuthInfo struct {
	UserID string
	Email  string
	Name   string
	EmailVerified bool
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type tokenClaims struct {
	Aud    string `json:"aud"`
	Iss    string `json:"iss"`
	Sub    string `json:"sub"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	EmailVerified bool `json:"email_verified"`
}

func NewVerifier(projectID string) *Verifier {
	return &Verifier{
		projectID: strings.TrimSpace(projectID),
		client:    http.DefaultClient,
		certs:     map[string]*rsa.PublicKey{},
	}
}

func (v *Verifier) VerifyIDToken(ctx context.Context, token string) (AuthInfo, error) {
	if v.projectID == "" {
		return AuthInfo{}, fmt.Errorf("%w: firebase project id is empty", ErrUnauthorized)
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return AuthInfo{}, ErrInvalidToken
	}

	header, err := decodeSegment[tokenHeader](parts[0])
	if err != nil {
		return AuthInfo{}, ErrInvalidToken
	}
	if header.Alg != "RS256" || header.Kid == "" {
		return AuthInfo{}, ErrInvalidToken
	}

	claims, err := decodeSegment[tokenClaims](parts[1])
	if err != nil {
		return AuthInfo{}, ErrInvalidToken
	}
	if claims.Sub == "" {
		return AuthInfo{}, ErrInvalidToken
	}
	if claims.Aud != v.projectID {
		return AuthInfo{}, ErrUnauthorized
	}
	expectedIssuer := fmt.Sprintf("https://securetoken.google.com/%s", v.projectID)
	if claims.Iss != expectedIssuer {
		return AuthInfo{}, ErrUnauthorized
	}
	now := time.Now().Unix()
	if claims.Exp == 0 || now > claims.Exp {
		return AuthInfo{}, ErrUnauthorized
	}

	key, err := v.keyFor(ctx, header.Kid)
	if err != nil {
		return AuthInfo{}, ErrUnauthorized
	}

	if err := verifySignature(key, parts[0], parts[1], parts[2]); err != nil {
		return AuthInfo{}, ErrUnauthorized
	}

	userID := claims.UserID
	if userID == "" {
		userID = claims.Sub
	}
	return AuthInfo{
		UserID: userID,
		Email:  claims.Email,
		Name:   claims.Name,
		EmailVerified: claims.EmailVerified,
	}, nil
}

func (v *Verifier) keyFor(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if time.Now().After(v.expiresAt) || len(v.certs) == 0 {
		if err := v.refreshCertsLocked(ctx); err != nil {
			return nil, err
		}
	}
	key, ok := v.certs[kid]
	if !ok {
		if err := v.refreshCertsLocked(ctx); err != nil {
			return nil, err
		}
		key, ok = v.certs[kid]
	}
	if !ok {
		return nil, ErrUnauthorized
	}
	return key, nil
}

func (v *Verifier) refreshCertsLocked(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, certsURL, nil)
	if err != nil {
		return err
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("certs fetch failed: %s", resp.Status)
	}

	var payload map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}

	certs := make(map[string]*rsa.PublicKey, len(payload))
	for kid, pemCert := range payload {
		block, _ := pem.Decode([]byte(pemCert))
		if block == nil {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		publicKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			continue
		}
		certs[kid] = publicKey
	}
	if len(certs) == 0 {
		return errors.New("no valid certs")
	}
	v.certs = certs
	v.expiresAt = time.Now().Add(parseMaxAge(resp.Header.Get("Cache-Control")))
	return nil
}

func parseMaxAge(cacheControl string) time.Duration {
	if cacheControl == "" {
		return time.Minute * 5
	}
	for _, part := range strings.Split(cacheControl, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "max-age=") {
			value := strings.TrimPrefix(part, "max-age=")
			if seconds, err := strconv.Atoi(value); err == nil && seconds > 0 {
				return time.Duration(seconds) * time.Second
			}
		}
	}
	return time.Minute * 5
}

func verifySignature(key *rsa.PublicKey, header, payload, signature string) error {
	rawSig, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	data := []byte(header + "." + payload)
	hash := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], rawSig)
}

func decodeSegment[T any](segment string) (T, error) {
	var out T
	data, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}
	return out, nil
}
