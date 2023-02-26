package project

import (
	"net/http"

	"github.com/SeungyeonHwang/tool-goaal/login"
	"github.com/SeungyeonHwang/tool-goaal/model"
	"github.com/unrolled/render"
)

var rd *render.Render = render.New()

type Handler struct {
	db model.DBHandler
}

func NewHandler(db model.DBHandler) *Handler {
	return &Handler{db}
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

	switch r.URL.Path {
	case "/projects/sorted-by-name":
		list = h.db.GetProjectsSortedByName(userId, sort)
	case "/projects/sorted-by-code":
		list = h.db.GetProjectsSortedByCode(userId, sort)
	case "/projects/sorted-by-priority":
		list = h.db.GetProjectsSortedByPriority(userId, sort)
	case "/projects/sorted-by-color":
		list = h.db.GetProjectsSortedByColor(userId, sort)
	default:
		list = h.db.GetProjects(userId, sort)
	}
	rd.JSON(w, http.StatusOK, list)
}
