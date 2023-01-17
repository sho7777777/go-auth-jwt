package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"main/dto"
	"main/service"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// エンドポイント
	r.HandleFunc("/auth/login", http.HandlerFunc(Login)).Methods(http.MethodPost, "OPTIONS")
	r.HandleFunc("/auth/refresh", http.HandlerFunc(Refresh)).Methods(http.MethodPost, "OPTIONS")
	r.HandleFunc("/auth/verify", http.HandlerFunc(Verify)).Methods(http.MethodGet, "OPTIONS")

	// サーバー起動
	fmt.Println("Starting server on port 8001")
	if err := http.ListenAndServe(":8001", r); err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	// CORS対応
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Access-Control-Allow-Origin")
	if r.Method == "OPTIONS" {
		return
	}

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
	// CORS対応
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Access-Control-Allow-Origin")
	if r.Method == "OPTIONS" {
		return
	}

	// アクセストークンとリフレッシュトークンを格納するDTO
	var refreshRequest dto.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		err := fmt.Errorf("エラーが発生しました: %v", err)
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		// アクセストークンの取得
		token, err := service.Refresh(refreshRequest)

		// リフレッシュトークンの有効期限が切れている場合
		if _, ok := err.(*jwt.TokenExpiredError); ok {
			writeResponse(w, http.StatusNotAcceptable, authResponse(false, "リフレッシュトークンの有効期限が切れました。"))
		} else if err != nil {
			writeResponse(w, http.StatusNotAcceptable, err.Error())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

func Verify(w http.ResponseWriter, r *http.Request) {

	urlParams := make(map[string]string)

	// URLからパラメータを取得し、mapに格納
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}
	// urlParams: map[accessToken:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzX3Rva2VuIiwiaWQiOiIxMjM0IiwidXNlcl9uYW1lIjoiSm9obiBEb2UiLCJyb2xlIjoidXNlciIsImV4cCI6MTY3MzkzOTYzOS4yOTE5MzF9.6AjWbUPSD0O8w9eL7TnZ9yyFCYw_WNrm6tRKSoJHTZQ id:1234 routeName:getResource]

	var err error
	// アクセストークンがある場合
	if urlParams["accessToken"] != "" {
		err = service.Verify(urlParams)

		// アクセストークンの有効期限が切れている場合
		if _, ok := err.(*jwt.TokenExpiredError); ok {
			writeResponse(w, http.StatusNotAcceptable, authResponse(false, "アクセストークンの有効期限が切れました。"))

		} else if err != nil {
			writeResponse(w, http.StatusNotAcceptable, authResponse(false, err.Error()))
		} else {
			writeResponse(w, http.StatusOK, authResponse(true, ""))
		}

		// アクセストークンがない場合
	} else {
		writeResponse(w, http.StatusForbidden, authResponse(false, "アクセストークンがありません。"))
	}
}

// アクセス可否とエラーメッセージを返す
func authResponse(isAuthorized bool, message string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": isAuthorized,
		"message":      message,
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
