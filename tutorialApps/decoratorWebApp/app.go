package decoratorwebapp

import (
	"fmt"
	"log"
	"net/http"
	"time"

	decohandler "github.com/SeungyeonHwang/tool-goaal/decoratorWebApp/decoHandler"
)

func logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER1] Started")
	h.ServeHTTP(w, r)
	log.Println("[LOGGER1] Completed time:", time.Since(start).Microseconds())
}

func logger2(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER2]  Started")
	h.ServeHTTP(w, r)
	log.Println("[LOGGER2] Completed time:", time.Since(start).Microseconds())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	return mux
}

func NewHttpHandler() http.Handler {
	h := NewHandler()
	h = decohandler.NewDecoHandler(h, logger)
	h = decohandler.NewDecoHandler(h, logger2)
	return h
}
