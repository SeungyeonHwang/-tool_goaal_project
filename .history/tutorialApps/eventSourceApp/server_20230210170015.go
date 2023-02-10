package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/antage/eventsource"
)

func main() {
	es := eventsource.New(nil, nil)
	defer es.Close()
	http.Handle("/events", es)
	go func() {
		id := 1
		for {
			es.SendEventMessage("tick", "", strconv.Itoa(id))
			id++
			time.Sleep(1 * time.Second)
		}
	}()
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":9000", nil))
}
