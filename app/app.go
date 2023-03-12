package app

import (
	"net/http"

	"github.com/SeungyeonHwang/tool-goaal/login"
	"github.com/SeungyeonHwang/tool-goaal/model"
	"github.com/SeungyeonHwang/tool-goaal/project"
	"github.com/SeungyeonHwang/tool-goaal/todo"
	"github.com/SeungyeonHwang/tool-goaal/user"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type AppHandler struct {
	http.Handler
	db      model.DBHandler
	project *project.Handler
	todo    *todo.Handler
	user    *user.Handler
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func MakeHandler(dbDir string) *AppHandler {
	r := mux.NewRouter()

	a := &AppHandler{
		db:      model.NewDBHandler(dbDir),
		project: project.NewHandler(model.NewDBHandler(dbDir)),
		todo:    todo.NewHandler(model.NewDBHandler(dbDir)),
		user:    user.NewHandler(model.NewDBHandler(dbDir)),
	}

	//USER
	r.HandleFunc("/auth/google/login", login.GoogleLoginHandler)
	r.HandleFunc("/auth/google/callback", login.GoogleAuthCallback)
	r.HandleFunc("/user/{id:[0-9]+}", a.user.GetUserInfoById).Methods("GET")

	//PROJECT
	r.HandleFunc("/", a.project.IndexHandler)
	r.HandleFunc("/projects", a.project.AddProjectListHandler).Methods("POST")
	r.HandleFunc("/projects", a.project.GetProjectListHandler).Methods("GET")
	r.HandleFunc("/projects/{id:[0-9]+}", a.project.GetProjectHandler).Methods("GET")
	r.HandleFunc("/projects/{id:[0-9]+}/check-edit-auth", a.project.CheckProjectEditAuthHandler).Methods("GET")
	r.HandleFunc("/projects/{id:[0-9]+}/participants", a.project.GetProjectParticipantListHandler).Methods("GET")
	r.HandleFunc("/projects/{id:[0-9]+}/availableUsers", a.project.GetProjectAvailableUsersListHandler).Methods("GET")
	r.HandleFunc("/projects/{id:[0-9]+}/todos", a.project.GoToTodoHandler).Methods("GET")
	r.HandleFunc("/projects/{id:[0-9]+}", a.project.UpdateProjectHandler).Methods("PUT")
	r.HandleFunc("/projects/{id:[0-9]+}", a.project.RemoveProjectHandler).Methods("DELETE")

	//TODO
	r.HandleFunc("/todos", a.todo.GetTodoListHandler).Methods("GET")
	r.HandleFunc("/todos/sorted-by-user", a.todo.GetTodoListHandler).Methods("GET")
	r.HandleFunc("/todos/sorted-by-completed", a.todo.GetTodoListHandler).Methods("GET")
	r.HandleFunc("/todos/sorted", a.todo.GetTodoListHandler).Methods("GET")
	r.HandleFunc("/todos/progress", a.todo.GetTodoProgressHandler).Methods("GET")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", a.todo.CompleteTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.todo.AddTodoListHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.todo.RemoveTodoListHandler).Methods("DELETE")
	r.HandleFunc("/todos/completed", a.todo.RemoveCompletedTodoListHandler).Methods("DELETE")

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(login.CheckLogin),
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			if !user.CheckAccessPermission(a.db, w, r) {
				return
			}
			next(w, r)
		}),
		negroni.NewStatic(http.Dir("public")),
	)

	n.UseHandler(r)
	a.Handler = n

	return a
}
