package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"book_manager/backend/internal/admininvitations"
	"book_manager/backend/internal/adminusers"
	"book_manager/backend/internal/config"
	"book_manager/backend/internal/firebaseauth"
	"book_manager/backend/internal/users"
	"book_manager/backend/internal/validation"
)

func (h *Handler) AuthSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		UserID      string `json:"userId"`
		DisplayName string `json:"displayName"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if req.Email == "" {
		badRequest(w, "email_required")
		return
	}
	if req.Password == "" {
		badRequest(w, "password_required")
		return
	}
	if !validation.IsValidEmail(req.Email) {
		badRequest(w, "invalid_email")
		return
	}
	if len(req.Password) < config.PasswordMinLength {
		badRequest(w, "password_too_short")
		return
	}
	normalizedUserID := strings.TrimSpace(req.UserID)
	if normalizedUserID == "" {
		badRequest(w, "user_id_required")
		return
	}
	if len(normalizedUserID) < config.UserIDMinLength {
		badRequest(w, "user_id_too_short")
		return
	}
	if len(normalizedUserID) > config.UserIDMaxLength {
		badRequest(w, "user_id_too_long")
		return
	}
	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" {
		displayName = normalizedUserID // 表示名が未指定の場合はユーザーIDと同じにする
	}
	if len(displayName) > config.DisplayNameMaxLength {
		badRequest(w, "display_name_too_long")
		return
	}
	if h.firebaseClient == nil {
		internalError(w)
		return
	}
	if h.users.IsUserIDTaken(normalizedUserID) {
		conflict(w, "user_id_exists")
		return
	}
	if h.adminUsers != nil && h.adminUsers.IsAdmin(normalizedUserID) {
		conflict(w, "user_id_reserved")
		return
	}
	result, err := h.firebaseClient.SignUp(req.Email, req.Password, normalizedUserID)
	if err != nil {
		switch {
		case errors.Is(err, firebaseauth.ErrEmailExists):
			conflict(w, "email_exists")
		case errors.Is(err, firebaseauth.ErrWeakPassword):
			badRequest(w, "weak_password")
		case errors.Is(err, firebaseauth.ErrInvalidEmail):
			badRequest(w, "invalid_email")
		case errors.Is(err, firebaseauth.ErrTooManyAttempts):
			http.Error(w, "too_many_attempts", http.StatusTooManyRequests)
		default:
			internalError(w)
		}
		return
	}
	if err := h.firebaseClient.SendEmailVerification(result.IDToken); err != nil {
		internalError(w)
		return
	}
	user, err := h.users.Create(result.LocalID, req.Email, normalizedUserID, displayName)
	if err != nil {
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken":  result.IDToken,
		"refreshToken": result.RefreshToken,
		"user": map[string]string{
			"id":          user.ID,
			"email":       user.Email,
			"userId":      user.UserID,
			"displayName": user.DisplayName,
		},
		"emailVerified": false,
	})
}

func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if req.Email == "" || req.Password == "" {
		badRequest(w, "email and password are required")
		return
	}
	if !validation.IsValidEmail(req.Email) {
		badRequest(w, "email is invalid")
		return
	}
	if h.firebaseClient == nil {
		internalError(w)
		return
	}
	result, err := h.firebaseClient.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, firebaseauth.ErrInvalidCredentials) || errors.Is(err, firebaseauth.ErrEmailNotFound) {
			unauthorized(w)
			return
		}
		internalError(w)
		return
	}
	if h.firebaseVerifier == nil {
		internalError(w)
		return
	}
	info, err := h.firebaseVerifier.VerifyIDToken(r.Context(), result.IDToken)
	if err != nil {
		log.Printf("WARNING: ID token verification failed after successful login: %v", err)
		unauthorized(w)
		return
	}
	user, ok := h.users.Get(result.LocalID)
	if !ok {
		userID := strings.TrimSpace(result.DisplayName)
		if userID == "" {
			userID = result.LocalID // Firebase DisplayNameが空の場合、LocalIDをデフォルトとして使用
		}
		// 管理者として予約されているUserIDを一般ユーザーが使用することを防ぐ
		if h.adminUsers != nil && h.adminUsers.IsAdmin(userID) {
			conflict(w, "user_id_reserved")
			return
		}
		created, err := h.users.Create(result.LocalID, result.Email, userID, userID)
		if err != nil {
			internalError(w)
			return
		}
		user = created
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken":  result.IDToken,
		"refreshToken": result.RefreshToken,
		"user": map[string]string{
			"id":          user.ID,
			"email":       user.Email,
			"userId":      user.UserID,
			"displayName": user.DisplayName,
		},
		"emailVerified": info.EmailVerified,
	})
}

func (h *Handler) AuthRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if req.RefreshToken == "" {
		badRequest(w, "refreshToken is required")
		return
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		badRequest(w, "refreshToken is required")
		return
	}
	if h.firebaseClient == nil {
		internalError(w)
		return
	}
	result, err := h.firebaseClient.Refresh(req.RefreshToken)
	if err != nil {
		if errors.Is(err, firebaseauth.ErrInvalidCredentials) {
			unauthorized(w)
			return
		}
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken":  result.IDToken,
		"refreshToken": result.RefreshToken,
	})
}

func (h *Handler) AuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		badRequest(w, "refreshToken is required")
		return
	}
	// リフレッシュトークンからユーザーIDを取得
	if h.firebaseClient == nil {
		internalError(w)
		return
	}
	result, err := h.firebaseClient.Refresh(req.RefreshToken)
	if err != nil {
		if errors.Is(err, firebaseauth.ErrInvalidCredentials) {
			// トークンが既に無効な場合は成功として扱う
			writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
			return
		}
		internalError(w)
		return
	}
	// Firebase Admin SDKでリフレッシュトークンを失効
	// 注: 失効に失敗してもログアウト自体は成功として扱う
	// 理由: クライアント側でトークンを削除すれば実質的にログアウトとなり、
	// ユーザー体験を優先する（Admin SDKが利用不可でも動作させる）
	if h.firebaseAdmin != nil {
		if err := h.firebaseAdmin.RevokeRefreshTokens(r.Context(), result.LocalID); err != nil {
			log.Printf("WARNING: failed to revoke refresh tokens for user %s: %v", result.LocalID, err)
		}
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) AuthResendVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		badRequest(w, "refreshToken is required")
		return
	}
	if h.firebaseClient == nil {
		internalError(w)
		return
	}
	result, err := h.firebaseClient.Refresh(req.RefreshToken)
	if err != nil {
		if errors.Is(err, firebaseauth.ErrInvalidCredentials) {
			unauthorized(w)
			return
		}
		internalError(w)
		return
	}
	if err := h.firebaseClient.SendEmailVerification(result.IDToken); err != nil {
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) AuthUpdateEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		methodNotAllowed(w, http.MethodPatch)
		return
	}
	var req struct {
		Email        string `json:"email"`
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	email := strings.TrimSpace(req.Email)
	if email == "" {
		badRequest(w, "email_required")
		return
	}
	if !validation.IsValidEmail(email) {
		badRequest(w, "invalid_email")
		return
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		badRequest(w, "refreshToken is required")
		return
	}
	if h.firebaseClient == nil {
		internalError(w)
		return
	}
	refreshed, err := h.firebaseClient.Refresh(req.RefreshToken)
	if err != nil {
		if errors.Is(err, firebaseauth.ErrInvalidCredentials) {
			unauthorized(w)
			return
		}
		internalError(w)
		return
	}
	updated, err := h.firebaseClient.UpdateEmail(refreshed.IDToken, email)
	if err != nil {
		if errors.Is(err, firebaseauth.ErrEmailExists) {
			conflict(w, "email_exists")
			return
		}
		internalError(w)
		return
	}
	if err := h.firebaseClient.SendEmailVerification(updated.IDToken); err != nil {
		internalError(w)
		return
	}
	user, err := h.users.UpdateProfile(updated.LocalID, nil, &email)
	if err != nil {
		if errors.Is(err, users.ErrEmailExists) {
			conflict(w, "email_exists")
			return
		}
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken":  updated.IDToken,
		"refreshToken": updated.RefreshToken,
		"user": map[string]string{
			"id":          user.ID,
			"email":       user.Email,
			"userId":      user.UserID,
			"displayName": user.DisplayName,
		},
		"emailVerified": false,
	})
}

func (h *Handler) AuthStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	if h.firebaseVerifier == nil {
		internalError(w)
		return
	}
	token := validation.BearerToken(r.Header.Get("Authorization"))
	if token == "" {
		unauthorized(w)
		return
	}
	info, err := h.firebaseVerifier.VerifyIDToken(r.Context(), token)
	if err != nil {
		unauthorized(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"id":            info.UserID,
			"email":         info.Email,
			"emailVerified": info.EmailVerified,
		},
	})
}

func (h *Handler) AuthSignupAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		Token       string `json:"token"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	token := strings.TrimSpace(req.Token)
	if token == "" {
		badRequest(w, "token_required")
		return
	}
	email := strings.TrimSpace(req.Email)
	if email == "" {
		badRequest(w, "email_required")
		return
	}
	if !validation.IsValidEmail(email) {
		badRequest(w, "invalid_email")
		return
	}
	password := req.Password
	if password == "" {
		badRequest(w, "password_required")
		return
	}
	if len(password) < config.PasswordMinLength {
		badRequest(w, "password_too_short")
		return
	}
	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" {
		displayName = email
	}
	if len(displayName) > config.DisplayNameMaxLength {
		badRequest(w, "display_name_too_long")
		return
	}
	invitation, err := h.adminInvitations.ValidateAndCheckEmail(token, email)
	if err != nil {
		if errors.Is(err, admininvitations.ErrInvitationNotFound) {
			badRequest(w, "invalid_token")
			return
		}
		if errors.Is(err, admininvitations.ErrInvitationExpired) {
			badRequest(w, "token_expired")
			return
		}
		if errors.Is(err, admininvitations.ErrInvitationUsed) {
			badRequest(w, "token_already_used")
			return
		}
		if errors.Is(err, admininvitations.ErrEmailMismatch) {
			badRequest(w, "email_mismatch")
			return
		}
		internalError(w)
		return
	}
	if h.firebaseClient == nil {
		internalError(w)
		return
	}
	result, err := h.firebaseClient.SignUp(email, password, invitation.UserID)
	if err != nil {
		if errors.Is(err, firebaseauth.ErrEmailExists) {
			conflict(w, "email_exists")
			return
		}
		internalError(w)
		return
	}
	if err := h.firebaseClient.SendEmailVerification(result.IDToken); err != nil {
		log.Printf("failed to send email verification: %v", err)
	}
	user, err := h.users.Create(result.LocalID, email, invitation.UserID, displayName)
	if err != nil {
		internalError(w)
		return
	}
	if err := h.adminInvitations.MarkUsed(invitation.ID, result.LocalID); err != nil {
		log.Printf("failed to mark invitation as used: %v", err)
	}
	if h.adminUsers != nil {
		if err := h.adminUsers.Add(invitation.UserID, invitation.CreatedBy); err != nil && !errors.Is(err, adminusers.ErrAlreadyAdmin) {
			log.Printf("failed to add admin user: %v", err)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken":  result.IDToken,
		"refreshToken": result.RefreshToken,
		"user": map[string]string{
			"id":          user.ID,
			"email":       user.Email,
			"userId":      user.UserID,
			"displayName": user.DisplayName,
		},
		"emailVerified": false,
		"isAdmin":       true,
	})
}
