package model

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	GetTodosSortedByUser(sessionId string) []*Todo
	GetTodosSortedByCompleted(sessionId string) []*Todo
	AddTodo(sessionId string, name string) *Todo

	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool

	GetProgress(sessionId string) int

	AddUser(sessionId string, email string, picture string)

	Close()
}

func NewDBHandler(dbDir string) DBHandler {
	return newSqliteHandler(dbDir)
}
