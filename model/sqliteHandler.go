package model

import (
	"database/sql"
	"math"
	"time"

	"github.com/leekchan/timeutil"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteHandler struct {
	db *sql.DB
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func (s *sqliteHandler) getTodosList(query string, sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query(query, sessionId)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.Id, &todo.Name, &todo.Picture, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *sqliteHandler) GetTodos(sessionId string) []*Todo {
	query := `
        SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
        FROM todos
        JOIN user ON todos.sessionId = user.sessionId
        WHERE todos.sessionId = ?`
	return s.getTodosList(query, sessionId)
}

// TODO
func (s *sqliteHandler) GetTodosSortedByUser(sessionId string) []*Todo {
	query := `
		SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
		FROM todos
		JOIN user ON todos.sessionId = user.sessionId
		WHERE todos.sessionId = ?`
	return s.getTodosList(query, sessionId)
}

func (s *sqliteHandler) GetTodosSortedByCompleted(sessionId string) []*Todo {
	query := `
		SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
		FROM todos
		JOIN user ON todos.sessionId = user.sessionId
		WHERE todos.sessionId = ? AND todos.completed = 1`
	return s.getTodosList(query, sessionId)
}

func (s *sqliteHandler) AddTodo(sessionId string, name string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}

	n := time.Now()
	formattedTime := timeutil.Strftime(&n, "%Y-%m-%d %H:%M")
	if err != nil {
		panic(err)
	}

	rs, err := stmt.Exec(sessionId, name, false, formattedTime)
	if err != nil {
		panic(err)
	}
	id, _ := rs.LastInsertId()
	var todo Todo
	todo.Id = int(id)
	todo.Name = name

	row := s.db.QueryRow("SELECT picture FROM user WHERE sessionId = ?", sessionId)
	err = row.Scan(&todo.Picture)
	if err != nil {
		panic(err)
	}

	todo.Completed = false
	todo.CreatedAt = formattedTime
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

func (s *sqliteHandler) GetProgress(sessionId string) int {
	rows, err := s.db.Query(`
		SELECT 
		COUNT(*) AS total_count, 
		COUNT(CASE WHEN completed = 1 THEN 1 ELSE NULL END) AS completed_count 
		FROM todos 
		WHERE sessionId = ?`, sessionId)
	if err != nil {
		panic(err)
	}

	var totalCount, completedCount int
	if rows.Next() {
		if err := rows.Scan(&totalCount, &completedCount); err != nil {
			panic(err)
		}
	}
	defer rows.Close()
	if totalCount == 0 {
		return 0
	}
	return int(math.Floor(float64(completedCount) / float64(totalCount) * 100))
}

func (s *sqliteHandler) AddUser(sessionId string, email string, picture string) {
	stmt, err := s.db.Prepare(
		`INSERT INTO user (sessionId, email, picture, createdAt)
			SELECT ?, ?, ?, datetime('now')
			WHERE NOT EXISTS (
				SELECT sessionId FROM user WHERE sessionId = ?
		);
		`)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(sessionId, email, picture, sessionId)
	if err != nil {
		panic(err)
	}
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
			completed BOOLEAN,
			createdAt STRING
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		)`)
	statement.Exec()

	statement, _ = database.Prepare(
		`CREATE TABLE IF NOT EXISTS user (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			sessionId STRING,
			email     TEXT,
			picture   TEXT,
			createdAt STRING
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		)`)
	statement.Exec()
	return &sqliteHandler{db: database}
}
