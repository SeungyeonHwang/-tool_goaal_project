package main

import (
	"log"
	"net/http"

	"github.com/SeungyeonHwang/tool-goaal/todo"
)

func main() {
	// port := os.Getenv("PORT")
	// mux := pat.New()

	m := todo.MakeHandler("./db/main.db")
	defer m.Close()

	log.Println("Start Goaal App...")
	// err := http.ListenAndServe(":"+port, m)
	err := http.ListenAndServe("127.0.0.1:3000", m)
	if err != nil {
		panic(err)
	}
}
