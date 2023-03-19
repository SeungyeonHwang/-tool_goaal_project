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

// 指定されたプロジェクトのTodoアイテムリストを取得するために使用されます。
func (h *Handler) GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
	projectId := r.FormValue("projectId")
	sort := r.FormValue("sort")
	filter := r.FormValue("filter")
	var list = make([]*model.Todo, 0)

	if sort != "" && filter != "" {
		switch filter {
		case "user":
			list = h.db.GetTodosSortedByUser(projectId, sort)
		case "completed":
			list = h.db.GetTodosSortedByCompleted(projectId, sort)
		default:
			list = h.db.GetTodos(projectId, sort)
		}
	} else {
		switch r.URL.Path {
		case "/todos/sorted-by-user":
			list = h.db.GetTodosSortedByUser(projectId, "")
		case "/todos/sorted-by-completed":
			list = h.db.GetTodosSortedByCompleted(projectId, "")
		default:
			list = h.db.GetTodos(projectId, "")
		}
	}
	rd.JSON(w, http.StatusOK, list)
}

// 新しいTodoアイテムを追加するために使用されます。
func (h *Handler) AddTodoListHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	sessionId := login.GetSessionId(r)
	userId := h.db.GetUserIdBySessionId(sessionId)
	projectId, _ := strconv.Atoi(r.FormValue("projectId"))
	todo := h.db.AddTodo(name, userId, projectId)
	rd.JSON(w, http.StatusCreated, todo)
}

// 指定されたTodoアイテムを削除するために使用されます。
func (h *Handler) RemoveTodoListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := h.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

// 指定されたプロジェクトの完了済みTodoアイテムを削除するために使用されます。
func (h *Handler) RemoveCompletedTodoListHandler(w http.ResponseWriter, r *http.Request) {
	projectId, _ := strconv.Atoi(r.FormValue("projectId"))
	ok := h.db.RemoveCompletedTodo(projectId)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

// 指定されたTodoアイテムを完了済みにマークするために使用されます。
func (h *Handler) CompleteTodoListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	ok := h.db.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

// 指定されたプロジェクトの進捗状況を取得するために使用されます。
func (h *Handler) GetTodoProgressHandler(w http.ResponseWriter, r *http.Request) {
	projectId, _ := strconv.Atoi(r.FormValue("projectId"))
	progress := h.db.GetProgress(projectId)
	rd.JSON(w, http.StatusOK, progress)
}
