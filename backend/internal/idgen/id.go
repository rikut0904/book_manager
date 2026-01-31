package idgen

import (
	"crypto/rand"
	"encoding/base64"
)

// New は新しいランダムIDを生成します。
// 生成に失敗した場合はfallbackPrefixに"_demo"を付けた値を返します。
func New(fallbackPrefix string) string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return fallbackPrefix + "_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}

// NewBook は書籍用のIDを生成します。
func NewBook() string {
	return New("book")
}

// NewUserBook はユーザー書籍用のIDを生成します。
func NewUserBook() string {
	return New("userbook")
}

// NewFavorite はお気に入り用のIDを生成します。
func NewFavorite() string {
	return New("fav")
}

// NewSeries はシリーズ用のIDを生成します。
func NewSeries() string {
	return New("series")
}

// NewRecommendation はおすすめ用のIDを生成します。
func NewRecommendation() string {
	return New("rec")
}

// NewNextToBuy は次に買う本用のIDを生成します。
func NewNextToBuy() string {
	return New("ntb")
}
