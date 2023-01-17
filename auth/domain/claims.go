package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type AccessTokenClaims struct {
	TokenType string `json:"token_type"`
	Id        string `json:"id"`
	Username  string `json:"user_name"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	TokenType string `json:"token_type"`
	Id        string `json:"id"`
	Username  string `json:"user_name"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

// リフレッシュトークンクレームからアクセストークンクレームを作成
func (r RefreshTokenClaims) GenerateAccessTokenClaims() AccessTokenClaims {
	return AccessTokenClaims{
		TokenType: "access_token",
		Id:        r.Id,
		Username:  r.Username,
		Role:      r.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(ACCESS_TOKEN_DURATION)),
		},
	}
}
