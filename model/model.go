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

type DBHandler interface {
	//User
	AddUser(sessionId string, email string, picture string)
	GetUserInfoById(id int) *User
	GetUserIdBySessionId(sessionId string) int

	//PROJECT
	AddProject(name string, code string, description string, color string, priority string, userId int) *Project
	GetProjects(userId int, sort string) []*Project
	GetProjectParticipants(projectId int) []*User
	GetProjectById(id int) *Project

	//TODO
	GetTodos(sessionId string, sort string) []*Todo
	GetTodosSortedByUser(sessionId string, sort string) []*Todo
	GetTodosSortedByCompleted(sessionId string, sort string) []*Todo
	AddTodo(sessionId string, name string) *Todo
	CompleteTodo(id int, complete bool) bool
	GetProgress(sessionId string) int

	RemoveTodo(id int) bool
	// RemoveCompletedTodo() bool

	Close()
}

func NewDBHandler(dbDir string) DBHandler {
	return newSqliteHandler(dbDir)
}
