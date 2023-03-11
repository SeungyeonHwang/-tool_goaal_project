package todo

import (
	"net/http"
	"strconv"

	"github.com/SeungyeonHwang/tool-goaal/login"
	"github.com/SeungyeonHwang/tool-goaal/model"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var rd *render.Render = render.New()

type Handler struct {
	db model.DBHandler
}

func NewHandler(db model.DBHandler) *Handler {
	return &Handler{db}
}

type Success struct {
	Success bool `json:"success"`
}

func (t *Handler) GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
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

func (t *Handler) AddTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := login.GetSessionId(r)
	name := r.FormValue("name")
	todo := t.db.AddTodo(sessionId, name)
	rd.JSON(w, http.StatusCreated, todo)
}

func (t *Handler) RemoveTodoListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := t.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

// func (t *Handler) RemoveCompletedTodoListHandler(w http.ResponseWriter, r *http.Request) {
// 	ok := t.db.RemoveCompletedTodo()
// 	if ok {
// 		rd.JSON(w, http.StatusOK, Success{true})
// 	} else {
// 		rd.JSON(w, http.StatusOK, Success{false})
// 	}
// }

func (t *Handler) CompleteTodoListHandler(w http.ResponseWriter, r *http.Request) {
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

func (t *Handler) GetTodoProgressHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := login.GetSessionId(r)
	progress := t.db.GetProgress(sessionId)
	rd.JSON(w, http.StatusOK, progress)
}
