package user

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

func (h *Handler) GetUserInfoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])
	userInfo := h.db.GetUserInfoById(userId)
	rd.JSON(w, http.StatusOK, userInfo)
}

func CheckAccessPermission(db model.DBHandler, w http.ResponseWriter, r *http.Request) bool {
	// check if the request is for /todo/todo.html?project-id={project_id}
	if r.URL.Path == "/todo/todo.html" {
		projectId, err := strconv.Atoi(r.URL.Query().Get("project-id"))
		if err != nil {
			http.Error(w, "Invalid project ID", http.StatusBadRequest)
			return false
		}
		sessionId := login.GetSessionId(r)
		userId := db.GetUserIdBySessionId(sessionId)

		// get participants list of the project
		participants := db.GetProjectParticipants(projectId)

		// check if the project exists
		if len(participants) == 0 {
			http.Error(w, "Project not found", http.StatusNotFound)
			return false
		}

		// check if the user is a participant of the project
		for _, p := range participants {
			if p.Id == userId {
				return true
			}
		}
		http.Error(w, "Access denied", http.StatusForbidden)
		return false
	}
	return true
}
