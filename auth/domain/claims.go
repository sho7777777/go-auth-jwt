package domain

import (
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
