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

// Project
func (s *sqliteHandler) GetUserIdBySessionId(sessionId string) int {
	var userId int
	err := s.db.QueryRow("SELECT id FROM user WHERE sessionId = ?", sessionId).Scan(&userId)
	if err != nil {
		panic(err)
	}
	return userId
}

func (s *sqliteHandler) AddProject(name string, code string, description string, color string, priority string, userId int) *Project {
	stmt, err := s.db.Prepare("INSERT INTO projects (name, code, description, color, createdAt, priority, userId) VALUES (?, ?, ?, ?, ?, ? ,?)")
	if err != nil {
		panic(err)
	}

	n := time.Now()
	formattedTime := timeutil.Strftime(&n, "%Y-%m-%d %H:%M")
	if err != nil {
		panic(err)
	}

	rs, err := stmt.Exec(name, code, description, color, formattedTime, priority, userId)
	if err != nil {
		panic(err)
	}

	id, _ := rs.LastInsertId()
	var project Project
	project.Id = int(id)
	project.Name = name
	project.Code = code
	project.Description = description
	project.Color = color
	project.CreatedAt = formattedTime
	project.Priority = priority
	project.UserId = userId

	s.addUserToProject(project.Id, userId)

	return &project
}

func (s *sqliteHandler) addUserToProject(projectId int, userId int) {
	stmt, err := s.db.Prepare("INSERT INTO project_users (projectId, userId) VALUES (?, ?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(projectId, userId)
	if err != nil {
		panic(err)
	}
}

func (s *sqliteHandler) getProjectsList(query string, userId int) []*Project {
	projects := []*Project{}
	rows, err := s.db.Query(query, userId)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var project Project
		rows.Scan(
			&project.Id,
			&project.Name,
			&project.Code,
			&project.Description,
			&project.Color,
			&project.Priority,
			&project.CreatedAt,
			&project.UserId,
		)
		projects = append(projects, &project)
	}
	return projects
}

func (s *sqliteHandler) GetProjects(userId int, sort string) []*Project {
	query := `
		SELECT projects.id, projects.name, projects.code, projects.description, projects.color, projects.priority, projects.createdAt, projects.userId
		FROM projects
		INNER JOIN project_users
		ON projects.id = project_users.projectId
		WHERE project_users.userId = ?
		ORDER BY projects.createdAt DESC`

	return s.getProjectsList(query, userId)
}

func (s *sqliteHandler) GetProjectsSortedByName(userId int, sort string) []*Project {
	query := `
		SELECT projects.id, projects.name, projects.code, projects.description, projects.color, projects.priority, projects.createdAt, projects.userId
		FROM projects
		INNER JOIN project_users
		ON projects.id = project_users.projectId
		WHERE project_users.userId = ?`

	switch sort {
	case "asc":
		query += " ORDER BY projects.name ASC"
	case "desc":
		query += " ORDER BY projects.name DESC"
	}
	return s.getProjectsList(query, userId)
}

func (s *sqliteHandler) GetProjectsSortedByCode(userId int, sort string) []*Project {
	query := `
		SELECT projects.id, projects.name, projects.code, projects.description, projects.color, projects.priority, projects.createdAt, projects.userId
		FROM projects
		INNER JOIN project_users
		ON projects.id = project_users.projectId
		WHERE project_users.userId = ?`

	switch sort {
	case "asc":
		query += " ORDER BY projects.code ASC"
	case "desc":
		query += " ORDER BY projects.code DESC"
	}
	return s.getProjectsList(query, userId)
}

func (s *sqliteHandler) GetProjectsSortedByPriority(userId int, sort string) []*Project {
	query := `
		SELECT projects.id, projects.name, projects.code, projects.description, projects.color, projects.priority, projects.createdAt, projects.userId
		FROM projects
		INNER JOIN project_users
		ON projects.id = project_users.projectId
		WHERE project_users.userId = ?`

	switch sort {
	case "asc":
		query += " ORDER BY projects.priority ASC"
	case "desc":
		query += " ORDER BY projects.priority DESC"
	}
	return s.getProjectsList(query, userId)
}

func (s *sqliteHandler) GetProjectsSortedByColor(userId int, sort string) []*Project {
	query := `
		SELECT projects.id, projects.name, projects.code, projects.description, projects.color, projects.priority, projects.createdAt, projects.userId
		FROM projects
		INNER JOIN project_users
		ON projects.id = project_users.projectId
		WHERE project_users.userId = ?`

	switch sort {
	case "asc":
		query += " ORDER BY projects.color ASC"
	case "desc":
		query += " ORDER BY projects.color DESC"
	}
	return s.getProjectsList(query, userId)
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

func (s *sqliteHandler) GetTodos(sessionId string, sort string) []*Todo {
	query := `
        SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
        FROM todos
        JOIN user ON todos.sessionId = user.sessionId
        WHERE todos.sessionId = ?`

	switch sort {
	case "asc":
		query += " ORDER BY todos.createdAt ASC"
	case "desc":
		query += " ORDER BY todos.createdAt DESC"
	}
	return s.getTodosList(query, sessionId)
}

// TODO
func (s *sqliteHandler) GetTodosSortedByUser(sessionId string, sort string) []*Todo {
	query := `
		SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
		FROM todos
		JOIN user ON todos.sessionId = user.sessionId
		WHERE todos.sessionId = ?`

	return s.getTodosList(query, sessionId)
}

func (s *sqliteHandler) GetTodosSortedByCompleted(sessionId string, sort string) []*Todo {
	query := `
		SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
		FROM todos
		JOIN user ON todos.sessionId = user.sessionId
		WHERE todos.sessionId = ? AND todos.completed = 0`

	switch sort {
	case "asc":
		query += " ORDER BY todos.createdAt ASC"
	case "desc":
		query += " ORDER BY todos.createdAt DESC"
	}
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

// TODO : Team todo delete
// func (s *sqliteHandler) RemoveCompletedTodo(id int) bool {
// 	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
// 	if err != nil {
// 		panic(err)
// 	}
// 	rs, err := stmt.Exec(id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	cnt, _ := rs.RowsAffected()
// 	return cnt > 0
// }

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
		`CREATE TABLE IF NOT EXISTS projects (
			id          INTEGER  PRIMARY KEY AUTOINCREMENT,
			name        TEXT,
			code        TEXT,
			description TEXT,
			color       TEXT,
			priority    TEXT,
			createdAt   STRING,
			userId     INTEGER,
			FOREIGN KEY (userId) REFERENCES user(id)
		);`)
	statement.Exec()

	statement, _ = database.Prepare(
		`CREATE TABLE IF NOT EXISTS project_users (
			id          INTEGER  PRIMARY KEY AUTOINCREMENT,
			projectId  INTEGER,
			userId     INTEGER,
			FOREIGN KEY (projectId) REFERENCES projects(id),
			FOREIGN KEY (userId) REFERENCES user(id)
		);`)
	statement.Exec()

	// CREATE TABLE IF NOT EXISTS project_todos (
	// 	id          INTEGER  PRIMARY KEY AUTOINCREMENT,
	// 	project_id  INTEGER,
	// 	todo_id     INTEGER,
	// 	FOREIGN KEY (project_id) REFERENCES projects(id),
	// 	FOREIGN KEY (todo_id) REFERENCES todos(id)
	// );

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

	statement, _ = database.Prepare(
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

	return &sqliteHandler{db: database}
}
