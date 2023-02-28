package user

import (
	"net/http"
	"strconv"

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
