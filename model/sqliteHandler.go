package model

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteHandler struct {
	db *sql.DB
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func (s *sqliteHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.Id, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *sqliteHandler) AddTodo(sessionId string, name string, userInfo map[string]interface{}) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, picture, ompleted, createdAt) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	picture := userInfo["picture"]
	email := userInfo["email"]

	log.Println(picture)
	log.Println(email)
	rs, err := stmt.Exec(sessionId, name, false)
	if err != nil {
		panic(err)
	}
	id, _ := rs.LastInsertId()
	var todo Todo
	todo.Id = int(id)
	todo.Name = name
	todo.Picture = picture.(string)
	todo.Email = email.(string)
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo
}

func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
	if err != nil {
		panic(err)
	}
	rs, err := stmt.Exec(complete, id)
	if err != nil {
		panic(err)
	}
	cnt, _ := rs.RowsAffected()
	return cnt > 0
}

func (s *sqliteHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
	if err != nil {
		panic(err)
	}
	rs, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	cnt, _ := rs.RowsAffected()
	return cnt > 0
}

func newSqliteHandler(dbDir string) DBHandler {
	database, err := sql.Open("sqlite3", dbDir)
	if err != nil {
		panic(err)
	}
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			sessionId STRING,
			name      TEXT,
			picture STRING,
			email STRING,
			completed BOOLEAN,
			createdAt DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		)`)
	statement.Exec()
	return &sqliteHandler{db: database}
}
