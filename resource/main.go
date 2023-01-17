package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	// エンドポイント
	r.HandleFunc("/resource/{id:[0-9]+}", getResource).Methods(http.MethodGet).Name("getResource")
	r.HandleFunc("/admin", adminPage).Methods(http.MethodGet).Name("adminPage")

	// 認証処理をミドルウェアとして利用
	r.Use(authHandler())

	fmt.Println("Start listening on port 8000...")
	http.ListenAndServe(":8000", r)
}

// リソース取得処理
func getResource(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "アクセストークンでリソース取得に成功！\n")
}

// 管理者ページ
func adminPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "管理者ページにログインしました。\n")
}

// 認証処理
func authHandler() func(http.Handler) http.Handler {

	// シグネチャを合わせるためにクロージャにする
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			currentRoute := mux.CurrentRoute(r)
			// currentRoute: getResource

			currentRouteVars := mux.Vars(r)
			// currentRouteVars: map[id:1234]

			// Authorizationに設定したアクセストークンを取得
			authHeader := r.Header.Get("Authorization")
			// authHeader: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzX3Rva2VuIiwiaWQiOiIxMjM0IiwidXNlcl9uYW1lIjoiSm9obiBEb2UiLCJyb2xlIjoidXNlciIsImV4cCI6MTY3MzkzOTEzNi45MzE4MTJ9.SCfR7gYpbw1kF64r3EKteiMIzXkxpGlxIDeCrEhhu7o

			if authHeader != "" {
				// アクセストークン取得
				accessToken := getTokenFromHeader(authHeader)

				//
				isAuthorized, msg := isAuthorized(accessToken, currentRoute.GetName(), currentRouteVars)

				if isAuthorized {
					next.ServeHTTP(w, r)
				} else {
					writeResponse(w, http.StatusNotAcceptable, msg)
				}
			} else {
				writeResponse(w, http.StatusUnauthorized, "アクセストークンがありません。")
			}
		})
	}
}

// ヘッダーからのアクセストークン取得処理
func getTokenFromHeader(header string) string {

	// Bearer askljfsdk... という形になっているので、アクセストークン部分のみ取り出す
	splitToken := strings.Split(header, "Bearer")
	if len(splitToken) == 2 {
		return strings.TrimSpace(splitToken[1])
	}
	return ""
}

func isAuthorized(accessToken string, routeName string, vars map[string]string) (bool, string) {

	// 認証サーバーのURLを構築
	u := buildVerifyUrl(accessToken, routeName, vars)
	// u: http://localhost:8001/auth/verify?accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzX3Rva2VuIiwiaWQiOiIxMjM0IiwidXNlcl9uYW1lIjoiSm9obiBEb2UiLCJyb2xlIjoidXNlciIsImV4cCI6MTY3MzkzOTEzNi45MzE4MTJ9.SCfR7gYpbw1kF64r3EKteiMIzXkxpGlxIDeCrEhhu7o&id=1234&routeName=getResource

	// 認証サーバーにリクエスト送信
	response, err := http.Get(u)

	if err != nil {
		fmt.Println("認証サーバーとの通信に失敗しました。" + err.Error())
		return false, ""
	} else {

		m := map[string]interface{}{}
		if err = json.NewDecoder(response.Body).Decode(&m); err != nil {
			fmt.Println("デコードに失敗しました: " + err.Error())
			return false, ""
		}

		isAuthorized := m["isAuthorized"].(bool)
		msg := m["message"].(string)
		return isAuthorized, msg
	}
}

func buildVerifyUrl(accessToken string, routeName string, vars map[string]string) string {

	// 認証サーバーのURL構築
	u := url.URL{Host: "localhost:8001", Path: "/auth/verify", Scheme: "http"}
	// u: {http   localhost:8001 /auth/verify  false false   }

	q := u.Query()
	fmt.Println("q: ", q)
	// q:  map[]

	// クエリパラメータをURlに設定
	q.Add("accessToken", accessToken)
	q.Add("routeName", routeName)
	for k, v := range vars {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// レスポンスの共通処理
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
