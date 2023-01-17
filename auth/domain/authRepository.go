package domain

import (
	"errors"
	"fmt"
)

// ユーザー取得処理
func FindBy(id, password string) (*User, error) {
	var u User

	// データ取得処理（本来はDB等から取得）
	if id == "1234" && password == "pass1234" {
		u = User{
			Id:       "1234",
			Username: "John Doe",
			Role:     "user",
		}
		return &u, nil
	}

	// 取得に失敗した場合、エラーを返す
	return nil, errors.New("Idかパスワードが正しくありません。")
}

// リフレッシュトークンをDBに格納
func SaveRefreshToken(rt string) error {

	// DBにリフレッシュトークンを保存する処理 ここでは保存できたとし、nilを返す
	fmt.Println("DBにリフレッシュトークンを保存しました。")
	return nil
}
