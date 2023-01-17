package service

import (
	"errors"
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

// リフレッシュできなかった時にnilを返したいので、dto.LoginResponseはポインタ型にする
func Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, error) {

	// 1. リフレッシュトークンがリクエストに存在するか確認
	if !request.RefreshTokenExists() {
		return nil, errors.New("リフレッシュトークンが存在しません。")
	}

	// 2. アクセストークンからトークンを復号する
	_, err := jwt.Parse(request.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})
	// {eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzX3Rva2VuIiwiaWQiOiIxMjM0IiwidXNlcl9uYW1lIjoiSm9obiBEb2UiLCJyb2xlIjoidXNlciIsImV4cCI6MTY3MzkzNjg3My4yNzM0MX0.IfLFBIAItf9Nc2Bm0hF9fx0A4WfdhnIB0bOu32DBsjs 0xc0000121c8 map[alg:HS256 typ:JWT] map[exp:1.67393687327341e+09 id:1234 role:user token_type:access_token user_name:John Doe] IfLFBIAItf9Nc2Bm0hF9fx0A4WfdhnIB0bOu32DBsjs false}

	if err != nil {
		var expErr *jwt.TokenExpiredError

		// 3. アクセストークンの有効期限が切れていた場合
		if errors.As(err, &expErr) {

			// 4. リフレッシュトークンがDBに存在するか確認
			if err = domain.RefreshTokenExists(request.RefreshToken); err != nil {
				return nil, err
			}

			// 5. リフレッシュトークンからトークンを生成
			token, err := jwt.ParseWithClaims(request.RefreshToken, &domain.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(mySigningKey), nil })
			if err != nil {
				return nil, err
			}
			// token: {eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImlkIjoiMTIzNCIsInVzZXJfbmFtZSI6IkpvaG4gRG9lIiwicm9sZSI6InVzZXIiLCJleHAiOjE2NzM5NDA0NDMuMjczNDE0fQ.96-OjqJRDyqGTtIND_rkiuW5hko_Tm5pprAAArzIXnI 0xc0000121c8 map[alg:HS256 typ:JWT] 0xc000000fa0 96-OjqJRDyqGTtIND_rkiuW5hko_Tm5pprAAArzIXnI true}

			// 6. トークンをリフレッシュトークンクレームで型アサーションする
			rtc, ok := token.Claims.(*domain.RefreshTokenClaims)
			if !ok {
				return nil, errors.New("リフレッシュトークン型への型アサーションに失敗しました。")
			}
			// rtc: {refresh_token 1234 John Doe user {[] 2023-01-17 16:27:23.273413 +0900 JST  <nil>  <nil> }}

			// 7. リフレッシュトークンクレームからアクセストークンクレームを作成
			atc := rtc.GenerateAccessTokenClaims()
			// atc: {access_token 1234 John Doe user {[] 2023-01-17 15:31:26.988402 +0900 JST  <nil>  <nil> }}

			// 8. アクセストークン（暗号化前）の作成
			tfa := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)
			// tfa: { 0xc0000121c8 map[alg:HS256 typ:JWT] {access_token 1234 John Doe user {[] 2023-01-17 15:31:26.988402 +0900 JST  <nil>  <nil> }}  false}

			// 9. アクセストークン（暗号化後）の作成
			var accessToken string
			if accessToken, err = tfa.SignedString([]byte(mySigningKey)); err != nil {
				return nil, fmt.Errorf("アクセストークンの暗号化に失敗しました。: %v", err)
			}
			// accessToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzX3Rva2VuIiwiaWQiOiIxMjM0IiwidXNlcl9uYW1lIjoiSm9obiBEb2UiLCJyb2xlIjoidXNlciIsImV4cCI6MTY3MzkzNzA4Ni45ODg0MDJ9.FpPb14rImzlmuiOFm6fAEHjc7cXdv6x7Tbx-5soxgoQ

			// 10. アクセストークンを返す（リフレッシュトークンは不要）
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}

		// アクセストークンの有効期限は切れていないが、有効なものでない場合
		return nil, errors.New("無効なアクセストークンです。")
	}

	// アクセストークンが有効な場合
	return nil, errors.New("アクセストークンの期限が切れていません。")
}
