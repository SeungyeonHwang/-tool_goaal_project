package etcwebapp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:created_at`
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "Seungyeon", Email: "syhwang.web@gmail.com"}

	rd.JSON(w, http.StatusOK, user)
	// w.Header().Add("Content-type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// data, _ := json.Marshal(user)
	// fmt.Fprint(w, string(data))
}

func addUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder((r.Body)).Decode(user)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprint(w, err)
		rd.Text(w, http.StatusBadRequest, err.Error())
		return
	}
	user.CreatedAt = time.Now()
	rd.JSON(w, http.StatusOK, user)
	// w.Header().Add("Content-type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// data, _ := json.Marshal(user)
	// fmt.Fprint(w, string(data))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "Seungyeon", Email: "syhwang.web@gmail.com"}
	rd = render.New(render.Options{
		Directory:  "etcWebApp/templates",
		Extensions: []string{".html", ".tmpl"},
		Layout:     "hello",
	})
	// tmpl, err := template.New("Hello").ParseFiles("etcWebApp/template/hello.tmpl")
	// tmpl.ExecuteTemplate(w, "hello.tmpl", "Seungyeon")
	rd.HTML(w, http.StatusOK, "body", user)
}

func NewHttpHandler() http.Handler {
	rd = render.New()
	mux := pat.New()
	mux.Get("/users", getUserInfoHandler)
	mux.Post("/users", addUserInfoHandler)
	mux.Get("/hello", helloHandler)
	// mux.Handle("/", http.FileServer(http.Dir("etcWebApp/public")))
	// /pulic
	n := negroni.Classic()
	n.UseHandler(mux)

	return n
}
