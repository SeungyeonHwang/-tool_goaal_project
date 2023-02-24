package app

import (
	"net/http"

	"github.com/SeungyeonHwang/tool-goaal/login"
	"github.com/SeungyeonHwang/tool-goaal/model"
	"github.com/SeungyeonHwang/tool-goaal/project"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type AppHandler struct {
	http.Handler
	db      model.DBHandler
	project *project.Handler
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func MakeHandler(dbDir string) *AppHandler {
	r := mux.NewRouter()
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(login.CheckLogin),
		negroni.NewStatic(http.Dir("public")),
	)

	n.UseHandler(r)

	a := &AppHandler{
		Handler: n,
		db:      model.NewDBHandler(dbDir),
		project: project.NewHandler(model.NewDBHandler(dbDir)),
	}

	//LOGIN
	r.HandleFunc("/auth/google/login", login.GoogleLoginHandler)
	r.HandleFunc("/auth/google/callback", login.GoogleAuthCallback)

	//HOME
	r.HandleFunc("/", a.project.IndexHandler)

	//TODO
	// r.HandleFunc("/todos", t.getTodoListHandler).Methods("GET")
	// r.HandleFunc("/todos/sorted-by-user", t.getTodoListHandler).Methods("GET")
	// r.HandleFunc("/todos/sorted-by-completed", t.getTodoListHandler).Methods("GET")
	// r.HandleFunc("/todos/sorted", t.getTodoListHandler).Methods("GET")

	// r.HandleFunc("/complete-todo/{id:[0-9]+}", t.completeTodoListHandler).Methods("GET")
	// r.HandleFunc("/todos/progress", t.getTodoProgressHandler).Methods("GET")

	// r.HandleFunc("/todos", t.addTodoListHandler).Methods("POST")

	// r.HandleFunc("/todos/{id:[0-9]+}", t.removeTodoListHandler).Methods("DELETE")
	// r.HandleFunc("/todos-completed-clear", t.removeCompletedTodoListHandler).Methods("DELETE")

	return a
}
