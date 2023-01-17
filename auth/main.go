package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"main/dto"
	"main/service"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// エンドポイント
	r.HandleFunc("/auth/login", http.HandlerFunc(Login)).Methods(http.MethodPost)
	r.HandleFunc("/auth/refresh", http.HandlerFunc(Refresh)).Methods(http.MethodPost)

	// サーバー起動
	fmt.Println("Starting server on port 8001")
	if err := http.ListenAndServe(":8001", r); err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	// リクエストで受け取ったIDとパスワード格納用DTO
	var loginRequest dto.LoginRequest

	// アクセストークン、リフレッシュトークン格納用DTO
	// 取得に失敗した際にサービスからnilを返したいのでポインタにする
	var loginResponse *dto.LoginResponse

	// リクエストをDTOにデコード
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		err := fmt.Errorf("エラーが発生しました: %v", err)
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		// アクセストークン、リフレッシュトークンの取得
		loginResponse, err = service.Login(loginRequest)
		if err != nil {
			writeResponse(w, http.StatusBadRequest, err.Error())
		} else {
			// アクセストークン、リフレッシュトークンをレスポンスとして返す
			writeResponse(w, http.StatusOK, *loginResponse)
		}
	}
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	// アクセストークンとリフレッシュトークンを格納するDTO
	var refreshRequest dto.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		err := fmt.Errorf("エラーが発生しました: %v", err)
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		// アクセストークンの取得
		token, err := service.Refresh(refreshRequest)
		if err != nil {
			writeResponse(w, http.StatusNotAcceptable, err.Error())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

// レスポンスの際の共通処理
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	// json形式でやり取り
	w.Header().Add("Content-Type", "application/json")
	// HTTPレスポンスコード
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
