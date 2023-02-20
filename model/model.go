package model

import (
	"time"
)

type Todo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(sessionId string, name string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool

	AddUser(sessionId string, email string, picture string)

	Close()
}

func NewDBHandler(dbDir string) DBHandler {
	return newSqliteHandler(dbDir)
}
