package main

import (
	"net/http"

	fileapp "github.com/SeungyeonHwang/tool-goaal/fileApp"
)

func main() {
	// httpApp
	// http.ListenAndServe("127.0.0.1:3000", httpApp.NewHttpHandler())

	//FileApp
	http.ListenAndServe("127.0.0.1:3000", fileapp.NewHttpHandler())
}
