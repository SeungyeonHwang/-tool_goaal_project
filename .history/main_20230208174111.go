package main

import "net/http"

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

func getUserInfoHandler(w http.ResponseWriter r *http.Request) {

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
}
