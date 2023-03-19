package login

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SeungyeonHwang/tool-goaal/model"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type LoginHandler struct {
	db model.DBHandler
}

type GoogleUserId struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

var googleOauthConfig = oauth2.Config{
	RedirectURL:  os.Getenv("DOMAIN_NAME") + "/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	// RedirectURL:  "http://localhost:3000/auth/google/callback",
	// ClientID:     "436991097398-h5bejll8dmsup6pi6r0gt0nk4sdjgai6.apps.googleusercontent.com",
	// ClientSecret: "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",
	Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: google.Endpoint,
}

// ユーザーがログインしているかどうかを確認するために使用されます。
func CheckLogin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// if reqeust URL is /signin.html, then next() => for login
	if strings.Contains(r.URL.Path, "/login") ||
		strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	// if user already signin
	sessionId := GetSessionId(r)
	if sessionId != "" {
		next(w, r)
		return
	}
	// if not user sign in => rediect login.html
	http.Redirect(w, r, "/login/login.html", http.StatusTemporaryRedirect)
}

// HTTPリクエストからセッションIDを取得するために使用されます。
var GetSessionId = func(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	// セッションに保存された"id"という名前の値を取得し、string型にキャストして返します。
	val := session.Values["id"]
	if val == nil {
		return ""
	}

	return val.(string)
}

// Google OAuth2.0を使用して、ユーザー情報を取得するために使用されます。
func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	// Google APIを使用して、ユーザー情報を取得します。
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body)
}

// OAuth2.0の認証フローで使用されるランダムな状態を生成し、Cookieに保存します。
// 生成された状態を返します。
func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)

	// ランダムなバイトスライスをURLエンコーディングして、stateという名前のCookieに保存します。
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

// Google OAuth2.0の認証フローを開始するために、
// ランダムな状態を生成し、Cookieに保存し、OAuth2.0認証ページにリダイレクトします。
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleのOAuth2.0を利用してログイン後のコールバックを処理するハンドラ
func GoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthstate.Value {
		errMsg := fmt.Sprintf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, r.FormValue("state"))
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Access Tokenを取得するためにGoogleのAPIを呼び出す
	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスデータをパースしてUserIdを取得する
	var userInfo GoogleUserId
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// UserIdをセッションに保存する
	session, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id"] = userInfo.Id
	session.Values["email"] = userInfo.Email
	session.Values["picture"] = userInfo.Picture

	err = session.Save(r, w)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	l := &LoginHandler{
		db: model.NewDBHandler("./db/main.db"),
	}
	log.Println(userInfo.Id)
	l.db.AddUser(userInfo.Id, userInfo.Email, userInfo.Picture)
	l.db.Close()

	//redirect to main
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
