package firebaseauth

import (
	"context"
	"encoding/json"
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
	// .envファイルで \n がリテラル文字列として保存されている場合に対応
	creds.PrivateKey = strings.ReplaceAll(creds.PrivateKey, "\\n", "\n")
	credentialsJSON, err := json.Marshal(map[string]string{
		"type":          "service_account",
		"project_id":    creds.ProjectID,
		"client_email":  creds.ClientEmail,
		"private_key":   creds.PrivateKey,
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
