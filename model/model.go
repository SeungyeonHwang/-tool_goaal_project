package model

type User struct {
	Id      int    `json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type Project struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Color       string `json:"color"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Priority    string `json:"priority"`
	UserId      int    `json:"user_id"`
}

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

/*
AddUser: ユーザーをデータベースに追加します。
GetUserInfoById: 指定されたIDのユーザー情報を取得します。
GetUserIdBySessionId: 指定されたセッションIDのユーザーIDを取得します。
AddProject: プロジェクトをデータベースに追加します。
GetProjects: 指定されたユーザーIDに紐付くプロジェクトのリストを取得します。
GetProjectParticipants: 指定されたプロジェクトに参加しているユーザーのリストを取得します。
GetProjectAvailableUsers: 指定されたプロジェクトに招待可能なユーザーのリストを取得します。
GetProjectById: 指定されたプロジェクトIDのプロジェクト情報を取得します。
CheckProjectEditAuth: 指定されたプロジェクトIDとセッションIDを使用して、ユーザーがプロジェクトの編集権限を持っているかどうかを確認します。
UpdateProject: 指定されたプロジェクトIDのプロジェクト情報を更新します。
RemoveProject: 指定されたプロジェクトをデータベースから削除します。
GetTodos: 指定されたプロジェクトに紐付くタスクのリストを取得します。
GetTodosSortedByUser: 指定されたプロジェクトに紐付くタスクをユーザー別にソートしたリストを取得します。
GetTodosSortedByCompleted: 指定されたプロジェクトに紐付くタスクを完了状態でソートしたリストを取得します。
AddTodo: 指定されたプロジェクトに新しいタスクを追加します。
CompleteTodo: 指定されたタスクの完了状態を更新します。
GetProgress: 指定されたプロジェクトの進捗状況を取得します。
RemoveTodo: 指定されたタスクをデータベースから削除します。
RemoveCompletedTodo: 指定されたプロジェクトから完了したタスクを削除します。
Close: データベース接続をクローズします。
*/

type DBHandler interface {
	//User
	AddUser(sessionId string, email string, picture string)
	GetUserInfoById(id int) *User
	GetUserIdBySessionId(sessionId string) int

	//PROJECT
	AddProject(name string, code string, description string, color string, priority string, userId int) *Project
	GetProjects(userId int, sort string) []*Project
	GetProjectParticipants(projectId int) []*User
	GetProjectAvailableUsers(projectId int) []*User
	GetProjectById(id int) *Project
	CheckProjectEditAuth(id int, sessionId string) bool
	UpdateProject(id int, name string, code string, description string, color string, priority string, userId int, participantIds []string, availableUserIds []string) *Project
	RemoveProject(id int) bool

	//TODO
	GetTodos(projectId string, sort string) []*Todo
	GetTodosSortedByUser(projectId string, sort string) []*Todo
	GetTodosSortedByCompleted(projectId string, sort string) []*Todo
	AddTodo(name string, userId int, projectId int) *Todo
	CompleteTodo(id int, complete bool) bool
	GetProgress(projectId int) int
	RemoveTodo(id int) bool
	RemoveCompletedTodo(projectId int) bool

	Close()
}

func NewDBHandler(dbDir string) DBHandler {
	return newSqliteHandler(dbDir)
}
