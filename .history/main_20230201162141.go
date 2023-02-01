package main

import (
	"fmt"
	"net/http"
)

type fooHandler struct{}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Foo")
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Heelo Bar!")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	http.HandleFunc("/bar", barHandler)

	http.Handle("/foo", &fooHandler{})

	http.ListenAndServe(":3000", nil)
}
