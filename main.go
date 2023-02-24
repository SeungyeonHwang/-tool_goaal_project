package main

import (
	"log"
	"net/http"

	"github.com/SeungyeonHwang/tool-goaal/app"
)

func main() {
	// port := os.Getenv("PORT")
	// mux := pat.New()

	// handler := todo.MakeHandler("./db/main.db")
	app := app.MakeHandler("./db/main.db")
	defer app.Close()

	log.Println("Start Goaal App...")
	// err := http.ListenAndServe(":"+port, m)
	err := http.ListenAndServe("127.0.0.1:3000", app)
	if err != nil {
		panic(err)
	}
}
