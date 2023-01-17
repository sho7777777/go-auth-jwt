package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

const ACCESS_TOKEN_DURATION = time.Second * 30 // 後でリフレッシュトークンの検証をするため、短めに設定
const REFRESH_TOKEN_DURATION = time.Minute * 60

type User struct {
	Id       string
	Username string
	Role     string
}

// アクセストークン用クレーム作成処理
func (l User) ClaimsForAccessToken() AccessTokenClaims {
	return AccessTokenClaims{
		TokenType: "access_token",
		Id:        l.Id,
		Username:  l.Username,
		Role:      l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(ACCESS_TOKEN_DURATION)),
		},
	}
}

// リフレッシュトークン用クレーム作成処理
func (l User) ClaimsForRefreshToken() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType: "refresh_token",
		Id:        l.Id,
		Username:  l.Username,
		Role:      l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(REFRESH_TOKEN_DURATION)),
		},
	}
}
