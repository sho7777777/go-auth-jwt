package service

import (
	"fmt"
	"main/domain"
	"main/dto"

	"github.com/dgrijalva/jwt-go/v4"
)

// 暗号化に使用するキー　環境変数等で設定する
const mySigningKey = "samplekey"

func Login(lr dto.LoginRequest) (*dto.LoginResponse, error) {

	// DBから取得したデータを格納する
	var u *domain.User
	var err error

	// ID,ユーザー名取得
	if u, err = domain.FindBy(lr.Id, lr.Password); err != nil {
		return nil, err
	}

	// 1. クレームを取得する
	cfa := u.ClaimsForAccessToken()
	cfr := u.ClaimsForRefreshToken()
	// cfa: {access_token 1234 John Doe user {[] 2023-01-17 14:09:42.833543 +0900 JST  <nil>  <nil> }}
	// cfr: {refresh_token 1234 John Doe user {[] 2023-01-17 15:09:12.833546 +0900 JST  <nil>  <nil> }}

	// 2. トークン（暗号化前）の作成
	tfa := jwt.NewWithClaims(jwt.SigningMethodHS256, cfa)
	tfr := jwt.NewWithClaims(jwt.SigningMethodHS256, cfr)
	// tfa: &{ 0xc0000a21b0 map[alg:HS256 typ:JWT] {access_token 1234 John Doe user {[] 2023-01-17 14:09:42.833543 +0900 JST  <nil>  <nil> }}  false}
	// tfr: &{ 0xc0000a21b0 map[alg:HS256 typ:JWT] {refresh_token 1234 John Doe user {[] 2023-01-17 15:09:12.833546 +0900 JST  <nil>  <nil> }}  false}

	// 3. トークン（暗号化後）の作成
	var accessToken string
	if accessToken, err = tfa.SignedString([]byte(mySigningKey)); err != nil {
		return nil, fmt.Errorf("アクセストークンの暗号化に失敗しました。: %v", err)
	}
	// accessToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzX3Rva2VuIiwiaWQiOiIxMjM0IiwidXNlcl9uYW1lIjoiSm9obiBEb2UiLCJyb2xlIjoidXNlciIsImV4cCI6MTY3MzkzMjE4Mi44MzM1NDI4fQ.1WvLFqiMz565zXMRknKcJS4JBlvEbGet2rNlZznuNKo

	// 3. リフレッシュトークン（暗号化後）の作成
	var refreshToken string
	if refreshToken, err = tfr.SignedString([]byte(mySigningKey)); err != nil {
		return nil, fmt.Errorf("リフレッシュトークンの暗号化に失敗しました。: %v", err)
	}
	// refreshToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImlkIjoiMTIzNCIsInVzZXJfbmFtZSI6IkpvaG4gRG9lIiwicm9sZSI6InVzZXIiLCJleHAiOjE2NzM5MzU3NTIuODMzNTQ2fQ.DET0gEfbbMaJ6BX9SWPIJEYqRCcJZZ0yKBPiRJgNnJ0

	// リフレッシュトークンはDBに保存する
	if err = domain.SaveRefreshToken(refreshToken); err != nil {
		return nil, fmt.Errorf("リフレッシュトークンの保存に失敗しました。: %v", err)
	}

	// アクセストークンを返す
	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
