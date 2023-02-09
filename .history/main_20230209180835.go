package main

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	sendMessage(name, msg)
}

func sendMessage(name, msg string) {
	//send message to every clients
}

func main() {
	mux := pat.New()
	mux.Post("/messages", postMessageHandler)

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe("127.0.0.1:3000", n)
}
