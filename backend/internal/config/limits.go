package config

// ページネーション設定
const (
	DefaultPageSize    = 20
	AdminPageSize      = 50
	MaxPageSize        = 200
	MaxPageNumber      = 1000000
)

// バリデーション設定
const (
	PasswordMinLength   = 8
	UserIDMinLength     = 2
	UserIDMaxLength     = 20
	DisplayNameMaxLength = 50
)
