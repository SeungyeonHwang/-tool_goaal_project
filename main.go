package main

import (
	"net/http"

	restfulapp "github.com/SeungyeonHwang/tool-goaal/restfulApp"
)

func main() {
	// httpApp
	// http.ListenAndServe("127.0.0.1:3000", httpApp.NewHttpHandler())

	//FileApp
	// http.ListenAndServe("127.0.0.1:3000", fileapp.NewHttpHandler())

	//RestfulAppÂ˜
	http.ListenAndServe("127.0.0.1:3000", restfulapp.NewHttpHandler())
}
