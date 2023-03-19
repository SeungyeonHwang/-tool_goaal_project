package main

import (
	"log"
	"net/http"

	"github.com/SeungyeonHwang/tool-goaal/app"
)

func main() {
	// port := os.Getenv("PORT")

	app := app.MakeHandler("./db/main.db")
	defer app.Close()

	log.Println("Start Goaal App...")

	// dev
	err := http.ListenAndServe("127.0.0.1:3000", app)

	// prod
	// err := http.ListenAndServe(":"+port, app)
	if err != nil {
		panic(err)
	}
}
