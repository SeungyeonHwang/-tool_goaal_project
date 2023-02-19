package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SeungyeonHwang/tool-goaal/todo"
	"github.com/antage/eventsource"
)

// Chat
func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	sendMessage(name, msg)
}

type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

var msgCh chan Message

func sendMessage(name, msg string) {
	//send message to every clients
	msgCh <- Message{name, msg}
}

func processMsgCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	sendMessage("", fmt.Sprintf("add user: %s", username))
}

func leftUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("leftUserHandler Called")
	username := r.FormValue("username")
	sendMessage("", fmt.Sprintf("left user: %s", username))
}

func main() {
	// port := os.Getenv("PORT")
	// mux := pat.New()

	//Chat
	// msgCh = make(chan Message)
	// es := eventsource.New(nil, nil)
	// defer es.Close()
	// go processMsgCh(es)

	// mux.Handle("/stream", es)
	// mux.Post("/messages", postMessageHandler)
	// mux.Post("/users", addUserHandler)
	// mux.Delete("/users", leftUserHandler)

	m := todo.MakeHandler("./todo.db")
	defer m.Close()

	log.Println("Start Goaal App...")
	// err := http.ListenAndServe(":"+port, m)
	err := http.ListenAndServe("127.0.0.1:3000", m)
	if err != nil {
		panic(err)
	}
}
