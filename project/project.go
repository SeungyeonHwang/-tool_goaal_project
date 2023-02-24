package project

import (
	"net/http"

	"github.com/SeungyeonHwang/tool-goaal/model"
)

type Handler struct {
	db model.DBHandler
}

func NewHandler(db model.DBHandler) *Handler {
	return &Handler{db}
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "project/project.html", http.StatusTemporaryRedirect)
}
