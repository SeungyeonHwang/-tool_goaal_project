package project

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

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "project/project.html", http.StatusTemporaryRedirect)
}

func (h *Handler) AddProjectListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := login.GetSessionId(r)
	name := r.FormValue("name")
	code := r.FormValue("code")
	description := r.FormValue("description")
	color := r.FormValue("color")
	priority := r.FormValue("priority")
	userId := h.db.GetUserIdBySessionId(sessionId)
	project := h.db.AddProject(name, code, description, color, priority, userId)
	rd.JSON(w, http.StatusCreated, project)
}

func (h *Handler) GetProjectListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := login.GetSessionId(r)
	userId := h.db.GetUserIdBySessionId(sessionId)
	sort := r.FormValue("sort")
	var list = make([]*model.Project, 0)
	list = h.db.GetProjects(userId, sort)
	rd.JSON(w, http.StatusOK, list)
}

func (h *Handler) GetProjectParticipantListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId, _ := strconv.Atoi(vars["id"])
	var list = make([]*model.User, 0)
	list = h.db.GetProjectParticipants(projectId)
	rd.JSON(w, http.StatusOK, list)
}

func (h *Handler) GetProjectAvailableUsersListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId, _ := strconv.Atoi(vars["id"])
	var list = make([]*model.User, 0)
	list = h.db.GetProjectAvailableUsers(projectId)
	rd.JSON(w, http.StatusOK, list)
}

func (h *Handler) GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	project := h.db.GetProjectById(id)
	rd.JSON(w, http.StatusOK, project)
}

func (h *Handler) CheckProjectEditAuthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	sessionId := login.GetSessionId(r)
	canEdit := h.db.CheckProjectEditAuth(id, sessionId)
	rd.JSON(w, http.StatusOK, canEdit)
}

func (h *Handler) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := r.ParseForm()
	if err != nil {
		rd.JSON(w, http.StatusBadRequest, "잘못된 요청입니다.")
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	description := r.FormValue("description")
	color := r.FormValue("color")
	priority := r.FormValue("priority")

	userIdStr := r.FormValue("managerId")
	userId, _ := strconv.Atoi(userIdStr)

	participantIds := r.Form["participantIds[]"]
	availableUserIds := r.Form["availableUserIds[]"]

	project := h.db.UpdateProject(id, name, code, description, color, priority, userId, participantIds, availableUserIds)

	rd.JSON(w, http.StatusOK, project)
}

func (h *Handler) RemoveProjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := h.db.RemoveProject(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (h *Handler) GoToTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoURL := "/todo/todo.html"
	http.Redirect(w, r, todoURL, http.StatusSeeOther)
}
