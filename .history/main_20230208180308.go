package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/pat"
)

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:created_at`
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "Seungyeon", Email: "syhwang.web@gmail.com"}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func addUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder((r.Body)).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user.CreatedAt = time.Now()
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("Hello").ParseFiles("hello.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	tmpl.ExecuteTemplate(w.Body)
}
func main() {
	// httpApp
	// http.ListenAndServe("127.0.0.1:3000", httpApp.NewHttpHandler())

	//FileApp
	// http.ListenAndServe("127.0.0.1:3000", fileapp.NewHttpHandler())

	//RestfulApp
	// http.ListenAndServe("127.0.0.1:3000", restfulapp.NewHttpHandler())

	//DecoratorWebApp
	// http.ListenAndServe("127.0.0.1:3000", decoratorwebapp.NewHttpHandler())

	// mux := mux.NewRouter()
	// mux.HandleFunc("/users", getUserInfoHandler).Methods("GET")
	// mux.HandleFunc("/users", addUserInfoHandler).Methods("POST")
	mux := pat.New()
	mux.Get("/users", getUserInfoHandler)
	mux.Post("/users", addUserInfoHandler)
	mux.Get("/helloHandler", helloHandler)

	http.ListenAndServe("127.0.0.1:3000", mux)
}
