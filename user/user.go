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

// Webアプリケーションのセキュリティ機能を実装するために使用されます。
// このメソッドは、データベースハンドラーとHTTPレスポンスライターとHTTPリクエストを引数として受け取り、ブール値を返します。
func CheckAccessPermission(db model.DBHandler, w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == "/todo/todo.html" {
		projectId, err := strconv.Atoi(r.URL.Query().Get("project-id"))
		if err != nil {
			http.Error(w, "Invalid project ID", http.StatusBadRequest)
			return false
		}
		sessionId := login.GetSessionId(r)
		userId := db.GetUserIdBySessionId(sessionId)

		// プロジェクトの参加者リストを取得
		participants := db.GetProjectParticipants(projectId)

		// プロジェクトが存在するかどうかを確認
		if len(participants) == 0 {
			http.Error(w, "Project not found", http.StatusNotFound)
			return false
		}

		// ユーザーがプロジェクトの参加者であるかどうかをチェック
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
