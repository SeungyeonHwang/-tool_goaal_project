package todo

import (
	"net/http"
	"strconv"

	"github.com/SeungyeonHwang/tool-goaal/login"
	"github.com/SeungyeonHwang/tool-goaal/model"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render = render.New()

type TodoHandler struct {
	http.Handler
	db model.DBHandler
}

type Success struct {
	Success bool `json:"success"`
}

func (t *TodoHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "todo/todo.html", http.StatusTemporaryRedirect)
}

func (t *TodoHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := login.GetSessionId(r)
	sort := r.FormValue("sort")
	filter := r.FormValue("filter")
	var list = make([]*model.Todo, 0)

	if sort != "" && filter != "" {
		switch filter {
		case "user":
			list = t.db.GetTodosSortedByUser(sessionId, sort)
		case "completed":
			list = t.db.GetTodosSortedByCompleted(sessionId, sort)
		default:
			list = t.db.GetTodos(sessionId, sort)
		}
	} else {
		switch r.URL.Path {
		case "/todos/sorted-by-user":
			list = t.db.GetTodosSortedByUser(sessionId, "")
		case "/todos/sorted-by-completed":
			list = t.db.GetTodosSortedByCompleted(sessionId, "")
		default:
			list = t.db.GetTodos(sessionId, "")
		}
	}

	rd.JSON(w, http.StatusOK, list)
}

func (t *TodoHandler) addTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := login.GetSessionId(r)
	name := r.FormValue("name")
	todo := t.db.AddTodo(sessionId, name)
	rd.JSON(w, http.StatusCreated, todo)
}

func (t *TodoHandler) removeTodoListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := t.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

// TODO
// func (t *TodoHandler) removeCompletedTodoListHandler(w http.ResponseWriter, r *http.Request) {
// 	ok := t.db.RemoveCompletedTodo()
// 	if ok {
// 		rd.JSON(w, http.StatusOK, Success{true})
// 	} else {
// 		rd.JSON(w, http.StatusOK, Success{false})
// 	}
// }

func (t *TodoHandler) completeTodoListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	ok := t.db.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (t *TodoHandler) getTodoProgressHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := login.GetSessionId(r)
	progress := t.db.GetProgress(sessionId)
	rd.JSON(w, http.StatusOK, progress)
}

func (t *TodoHandler) Close() {
	t.db.Close()
}

func MakeHandler(dbDir string) *TodoHandler {
	r := mux.NewRouter()
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(login.CheckLogin),
		negroni.NewStatic(http.Dir("public")),
	)

	n.UseHandler(r)

	t := &TodoHandler{
		Handler: n,
		db:      model.NewDBHandler(dbDir),
	}

	//HOME
	r.HandleFunc("/", t.indexHandler)

	//LOGIN
	r.HandleFunc("/auth/google/login", login.GoogleLoginHandler)
	r.HandleFunc("/auth/google/callback", login.GoogleAuthCallback)

	//TODO
	r.HandleFunc("/todos", t.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos/sorted-by-user", t.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos/sorted-by-completed", t.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos/sorted", t.getTodoListHandler).Methods("GET")

	r.HandleFunc("/complete-todo/{id:[0-9]+}", t.completeTodoListHandler).Methods("GET")
	r.HandleFunc("/todos/progress", t.getTodoProgressHandler).Methods("GET")

	r.HandleFunc("/todos", t.addTodoListHandler).Methods("POST")

	r.HandleFunc("/todos/{id:[0-9]+}", t.removeTodoListHandler).Methods("DELETE")
	// r.HandleFunc("/todos-completed-clear", t.removeCompletedTodoListHandler).Methods("DELETE")

	return t
}
