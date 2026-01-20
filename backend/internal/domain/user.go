package domain

type User struct {
	ID           string // 内部システムID
	Email        string
	UserID       string // ユーザーID（ログインに使用、変更不可）
	DisplayName  string // 表示名（自由に変更可能）
	PasswordHash string
}
