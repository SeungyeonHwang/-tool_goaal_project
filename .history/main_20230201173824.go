package main

import (
	"net/http"

	"github.com/SeungyeonHwang/tool-goaal/myapp"
)

func main() {

	http.ListenAndServe("127.0.0.1:3000", myapp.NewHttpHandler())
}
