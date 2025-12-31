package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type Service struct {
	users  repository.UserRepository
	tokens *TokenStore
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         domain.User
}

func NewService(users repository.UserRepository) *Service {
	return &Service{
		users:  users,
		tokens: NewTokenStore(),
	}
}

func (s *Service) SeedUser(id, email, username, password string) error {
	hashed, err := hashPassword(password)
	if err != nil {
		return err
	}
	user := domain.User{
		ID:           id,
		Email:        email,
		Username:     username,
		PasswordHash: hashed,
	}
	if err := s.users.Create(user); err != nil {
		return err
	}
	return nil
}

func (s *Service) Signup(email, password, username string) (AuthResult, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return AuthResult{}, err
	}
	user := domain.User{
		ID:           newID(),
		Email:        email,
		Username:     username,
		PasswordHash: hashed,
	}
	if err := s.users.Create(user); err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return AuthResult{}, ErrUserExists
		}
		return AuthResult{}, err
	}

	accessToken, refreshToken := s.tokens.Issue(user.ID)
	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *Service) Login(email, password string) (AuthResult, error) {
	user, ok := s.users.FindByEmail(email)
	if !ok {
		return AuthResult{}, ErrInvalidCredentials
	}
	if !verifyPassword(user.PasswordHash, password) {
		return AuthResult{}, ErrInvalidCredentials
	}
	accessToken, refreshToken := s.tokens.Issue(user.ID)
	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *Service) Refresh(refreshToken string) (string, string, error) {
	accessToken, nextRefreshToken, ok := s.tokens.Rotate(refreshToken)
	if !ok {
		return "", "", ErrInvalidCredentials
	}
	return accessToken, nextRefreshToken, nil
}

func (s *Service) Logout(refreshToken string) {
	s.tokens.Revoke(refreshToken)
}

type TokenStore struct {
	mu      sync.Mutex
	refresh map[string]string
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		refresh: make(map[string]string),
	}
}

func (s *TokenStore) Issue(userID string) (string, string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	accessToken := newToken()
	refreshToken := newToken()
	s.refresh[refreshToken] = userID
	return accessToken, refreshToken
}

func (s *TokenStore) Rotate(refreshToken string) (string, string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID, ok := s.refresh[refreshToken]
	if !ok {
		return "", "", false
	}
	delete(s.refresh, refreshToken)
	newAccess := newToken()
	newRefresh := newToken()
	s.refresh[newRefresh] = userID
	return newAccess, newRefresh, true
}

func (s *TokenStore) Revoke(refreshToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.refresh, refreshToken)
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	encodedSalt := base64.RawURLEncoding.EncodeToString(salt)
	sum := sha256.Sum256([]byte(encodedSalt + ":" + password))
	return fmt.Sprintf("%s:%x", encodedSalt, sum[:]), nil
}

func verifyPassword(stored, password string) bool {
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return false
	}
	salt := parts[0]
	sum := sha256.Sum256([]byte(salt + ":" + password))
	return stored == fmt.Sprintf("%s:%x", salt, sum[:])
}

func newToken() string {
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return "token"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}

func newID() string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return "user_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}
