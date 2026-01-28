package firebaseauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrEmailExists         = errors.New("email exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserDisabled        = errors.New("user disabled")
	ErrEmailNotFound       = errors.New("email not found")
	ErrOperationNotAllowed = errors.New("operation not allowed")
	ErrWeakPassword        = errors.New("weak password")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrTooManyAttempts     = errors.New("too many attempts")
)

type Client struct {
	apiKey string
	client *http.Client
}

type AuthResult struct {
	IDToken      string
	RefreshToken string
	LocalID      string
	Email        string
	DisplayName  string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: strings.TrimSpace(apiKey),
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SignUp(email, password, displayName string) (AuthResult, error) {
	payload := map[string]any{
		"email":             email,
		"password":          password,
		"displayName":       displayName,
		"returnSecureToken": true,
	}
	var resp signInResponse
	if err := c.postJSON("accounts:signUp", payload, &resp); err != nil {
		return AuthResult{}, err
	}
	return AuthResult{
		IDToken:      resp.IDToken,
		RefreshToken: resp.RefreshToken,
		LocalID:      resp.LocalID,
		Email:        resp.Email,
		DisplayName:  resp.DisplayName,
	}, nil
}

func (c *Client) Login(email, password string) (AuthResult, error) {
	payload := map[string]any{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}
	var resp signInResponse
	if err := c.postJSON("accounts:signInWithPassword", payload, &resp); err != nil {
		return AuthResult{}, err
	}
	return AuthResult{
		IDToken:      resp.IDToken,
		RefreshToken: resp.RefreshToken,
		LocalID:      resp.LocalID,
		Email:        resp.Email,
		DisplayName:  resp.DisplayName,
	}, nil
}

func (c *Client) SendEmailVerification(idToken string) error {
	payload := map[string]any{
		"requestType": "VERIFY_EMAIL",
		"idToken":     idToken,
	}
	var resp map[string]any
	return c.postJSON("accounts:sendOobCode", payload, &resp)
}

func (c *Client) UpdateEmail(idToken, email string) (AuthResult, error) {
	payload := map[string]any{
		"idToken":           idToken,
		"email":             email,
		"returnSecureToken": true,
	}
	var resp signInResponse
	if err := c.postJSON("accounts:update", payload, &resp); err != nil {
		return AuthResult{}, err
	}
	return AuthResult{
		IDToken:      resp.IDToken,
		RefreshToken: resp.RefreshToken,
		LocalID:      resp.LocalID,
		Email:        resp.Email,
		DisplayName:  resp.DisplayName,
	}, nil
}
func (c *Client) Refresh(refreshToken string) (AuthResult, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	endpoint := "https://securetoken.googleapis.com/v1/token?key=" + url.QueryEscape(c.apiKey)
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return AuthResult{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return AuthResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return AuthResult{}, decodeFirebaseError(resp)
	}
	var payload refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return AuthResult{}, err
	}
	return AuthResult{
		IDToken:      payload.IDToken,
		RefreshToken: payload.RefreshToken,
		LocalID:      payload.UserID,
	}, nil
}

func (c *Client) postJSON(path string, payload any, out any) error {
	if c.apiKey == "" {
		return errors.New("firebase api key is empty")
	}
	endpoint := "https://identitytoolkit.googleapis.com/v1/" + path + "?key=" + url.QueryEscape(c.apiKey)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return decodeFirebaseError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

type signInResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	LocalID      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
}

type refreshResponse struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

type firebaseError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func decodeFirebaseError(resp *http.Response) error {
	var payload firebaseError
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return errors.New("firebase request failed")
	}
	msg := payload.Error.Message
	// Firebaseのエラーメッセージは "ERROR_CODE : detail" の形式の場合がある
	if idx := strings.Index(msg, " : "); idx > 0 {
		msg = msg[:idx]
	}
	switch msg {
	case "EMAIL_EXISTS":
		return ErrEmailExists
	case "INVALID_ID_TOKEN":
		return ErrInvalidCredentials
	case "TOKEN_EXPIRED":
		return ErrInvalidCredentials
	case "INVALID_PASSWORD":
		return ErrInvalidCredentials
	case "EMAIL_NOT_FOUND":
		return ErrEmailNotFound
	case "USER_DISABLED":
		return ErrUserDisabled
	case "OPERATION_NOT_ALLOWED":
		return ErrOperationNotAllowed
	case "WEAK_PASSWORD":
		return ErrWeakPassword
	case "INVALID_EMAIL":
		return ErrInvalidEmail
	case "TOO_MANY_ATTEMPTS_TRY_LATER":
		return ErrTooManyAttempts
	default:
		if payload.Error.Message == "" {
			return errors.New("firebase request failed")
		}
		return errors.New(payload.Error.Message)
	}
}
