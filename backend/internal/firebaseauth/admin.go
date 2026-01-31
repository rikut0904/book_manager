package firebaseauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// AdminClient はFirebase Admin SDKのラッパーです
type AdminClient struct {
	authClient *auth.Client
}

// AdminCredentials はFirebase Admin SDKの認証情報です
type AdminCredentials struct {
	ProjectID   string
	ClientEmail string
	PrivateKey  string
}

// NewAdminClient はFirebase Admin SDKクライアントを初期化します
func NewAdminClient(ctx context.Context, creds AdminCredentials) (*AdminClient, error) {
	privateKey := normalizePrivateKey(creds.PrivateKey)
	credentialsJSON, err := json.Marshal(map[string]string{
		"type":          "service_account",
		"project_id":    creds.ProjectID,
		"client_email":  creds.ClientEmail,
		"private_key":   privateKey,
		"token_uri":     "https://oauth2.googleapis.com/token",
	})
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(credentialsJSON)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &AdminClient{authClient: authClient}, nil
}

// RevokeRefreshTokens は指定したユーザーのすべてのリフレッシュトークンを失効させます
func (c *AdminClient) RevokeRefreshTokens(ctx context.Context, uid string) error {
	return c.authClient.RevokeRefreshTokens(ctx, uid)
}

// GetUser はユーザー情報を取得します
func (c *AdminClient) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return c.authClient.GetUser(ctx, uid)
}

// DeleteUser はFirebaseからユーザーを削除します
func (c *AdminClient) DeleteUser(ctx context.Context, uid string) error {
	return c.authClient.DeleteUser(ctx, uid)
}

// normalizePrivateKey は秘密鍵を正規化します
// 以下の形式をサポート:
// 1. Base64エンコードされた秘密鍵（推奨）
// 2. リテラル \n を含む秘密鍵（後方互換性）
// 3. 実際の改行を含む秘密鍵
func normalizePrivateKey(key string) string {
	key = strings.TrimSpace(key)

	// リテラル \n を含む場合（後方互換性）- 最初にチェック
	if strings.Contains(key, "\\n") {
		log.Println("INFO: Firebase private key contains literal \\n, converting to newlines")
		return strings.ReplaceAll(key, "\\n", "\n")
	}

	// 既にPEM形式の場合はそのまま返す
	if strings.HasPrefix(key, "-----BEGIN") {
		return key
	}

	// Base64エンコードされている場合
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err == nil && strings.HasPrefix(string(decoded), "-----BEGIN") {
		log.Println("INFO: Firebase private key decoded from Base64")
		return string(decoded)
	}

	// どの形式にも該当しない場合は元の値を返す（エラーは初期化時に発生する）
	log.Printf("WARNING: Firebase private key format not recognized")
	return key
}
