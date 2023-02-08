package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "Seungyeon", Email: "syhwang.web@gmail.com"}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
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

	mux := http.NewServeMux()

	mux.HandleFunc("/users", getUserInfoHandler)
	http.ListenAndServe("127.0.0.1:3000", mux)
}
