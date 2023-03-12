package model

import (
	"database/sql"
	"log"
	"math"
	"time"

	"github.com/SeungyeonHwang/tool-goaal/util"
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

func (s *sqliteHandler) GetUserInfoById(id int) *User {
	var user User
	err := s.db.QueryRow("SELECT id, email, picture FROM user WHERE id = ?", id).Scan(&user.Id, &user.Email, &user.Picture)
	if err != nil {
		panic(err)
	}
	return &user
}

func (s *sqliteHandler) AddProject(name string, code string, description string, color string, priority string, userId int) *Project {
	stmt, err := s.db.Prepare("INSERT INTO projects (name, code, description, color, createdAt, priority, userId) VALUES (?, ?, ?, ?, ?, ? ,?)")
	if err != nil {
		panic(err)
	}

	n := time.Now()
	formattedTime := timeutil.Strftime(&n, "%Y-%m-%d %H:%M:%S")
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

func (s *sqliteHandler) GetProjectParticipants(id int) []*User {
	users := []*User{}
	query := `
		SELECT DISTINCT user.id, user.email, user.picture
        FROM user
        JOIN project_users ON user.id = project_users.userId
        WHERE project_users.projectId = ?`
	rows, err := s.db.Query(query, id)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		rows.Scan(
			&user.Id,
			&user.Email,
			&user.Picture,
		)
		users = append(users, &user)
	}
	return users
}

func (s *sqliteHandler) GetProjectAvailableUsers(id int) []*User {
	users := []*User{}
	query := `
	SELECT DISTINCT user.id, user.email, user.picture
	FROM user
	LEFT JOIN project_users ON user.id = project_users.userId AND project_users.projectId = ?
	WHERE project_users.userId IS NULL;`
	rows, err := s.db.Query(query, id)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		rows.Scan(
			&user.Id,
			&user.Email,
			&user.Picture,
		)
		users = append(users, &user)
	}
	return users
}

func (s *sqliteHandler) GetProjectById(id int) *Project {
	var project Project

	query := "SELECT id, name, code, description, color, priority, createdAt, userId FROM projects WHERE id = ?"
	row := s.db.QueryRow(query, id)
	err := row.Scan(
		&project.Id,
		&project.Name,
		&project.Code,
		&project.Description,
		&project.Color,
		&project.Priority,
		&project.CreatedAt,
		&project.UserId,
	)

	if err != nil {
		panic(err)
	}

	return &project
}

func (s *sqliteHandler) CheckProjectEditAuth(id int, sessionId string) bool {
	var userSessionId string
	f, err := util.ParseFloat(sessionId)
	if err != nil {
		panic(err)
	}
	sessionIdStr := util.FormatScientific(f)

	err = s.db.QueryRow(`
						SELECT u.sessionId
						FROM projects p
						JOIN user u ON p.userId = u.id
						WHERE p.id=?`, id).Scan(&userSessionId)
	if err != nil {
		panic(err)
	}
	return userSessionId == sessionIdStr
}

func (s *sqliteHandler) UpdateProject(id int, name string, code string, description string, color string, priority string, userId int, participantIds []string, availableUserIds []string) *Project {
	stmt, err := s.db.Prepare("UPDATE projects SET name=?, code=?, description=?, color=?, priority=?, userId=? WHERE id=?")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(name, code, description, color, priority, userId, id)
	if err != nil {
		panic(err)
	}

	var project Project
	row := s.db.QueryRow("SELECT id, name, code, description, color, priority, userId FROM projects WHERE id = ?", id)
	err = row.Scan(&project.Id, &project.Name, &project.Code, &project.Description, &project.Color, &project.Priority, &project.UserId)
	if err != nil {
		panic(err)
	}

	//참가자 추가
	for _, participantId := range participantIds {
		var count int
		row := s.db.QueryRow("SELECT COUNT(*) FROM project_users WHERE projectId=? AND userId=?", id, participantId)
		err := row.Scan(&count)
		if err != nil {
			panic(err)
		}
		if count == 0 {
			_, err = s.db.Exec("INSERT INTO project_users(projectId, userId) VALUES(?, ?)", id, participantId)
			if err != nil {
				panic(err)
			}
		}
	}

	// 참가자 제거
	for _, availableUserId := range availableUserIds {
		// 사용자가 프로젝트에 참여 중인지 확인
		var count int
		row := s.db.QueryRow("SELECT COUNT(*) FROM project_users WHERE projectId=? AND userId=?", project.Id, availableUserId)
		err := row.Scan(&count)
		if err != nil {
			panic(err)
		}
		if count > 0 {
			// 사용자가 프로젝트에 참가 중이므로 제거
			_, err = s.db.Exec("DELETE FROM project_users WHERE projectId=? AND userId=?", id, availableUserId)
			if err != nil {
				panic(err)
			}
		}
	}

	return &project
}

func (s *sqliteHandler) RemoveProject(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM projects WHERE id=?")
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

func (s *sqliteHandler) getTodosList(query string, projectId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query(query, projectId)

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

func (s *sqliteHandler) GetTodos(projectId string, sort string) []*Todo {
	query := `
		SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
		FROM project_todos
		JOIN todos ON project_todos.todoId = todos.id
		JOIN user ON todos.userId = user.id
		WHERE project_todos.projectId = ?`

	switch sort {
	case "asc":
		query += " ORDER BY todos.createdAt ASC"
	case "desc":
		query += " ORDER BY todos.createdAt DESC"
	}
	return s.getTodosList(query, projectId)
}

// TODO:
func (s *sqliteHandler) GetTodosSortedByUser(projectId string, sort string) []*Todo {
	query := `
		SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
		FROM project_todos
		JOIN todos ON project_todos.todoId = todos.id
		JOIN user ON todos.userId = user.id
		WHERE project_todos.projectId = ?`

	return s.getTodosList(query, projectId)
}

func (s *sqliteHandler) GetTodosSortedByCompleted(projectId string, sort string) []*Todo {
	query := `
		SELECT todos.id, todos.name, user.picture, todos.completed, todos.createdAt
		FROM project_todos
		JOIN todos ON project_todos.todoId = todos.id
		JOIN user ON todos.userId = user.id
		WHERE project_todos.projectId = ? AND todos.completed = 0`

	switch sort {
	case "asc":
		query += " ORDER BY todos.createdAt ASC"
	case "desc":
		query += " ORDER BY todos.createdAt DESC"
	}
	return s.getTodosList(query, projectId)
}

func (s *sqliteHandler) AddTodo(name string, userId int, projectId int) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (userId, name, completed, createdAt) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}

	n := time.Now()
	formattedTime := timeutil.Strftime(&n, "%Y-%m-%d %H:%M:%S")
	if err != nil {
		panic(err)
	}

	rs, err := stmt.Exec(userId, name, false, formattedTime)
	if err != nil {
		panic(err)
	}
	id, _ := rs.LastInsertId()
	var todo Todo
	todo.Id = int(id)
	todo.Name = name

	row := s.db.QueryRow("SELECT picture FROM user WHERE id = ?", userId)
	err = row.Scan(&todo.Picture)
	if err != nil {
		panic(err)
	}

	todo.Completed = false
	todo.CreatedAt = formattedTime

	s.addTodoToProject(projectId, todo.Id)
	return &todo
}

func (s *sqliteHandler) addTodoToProject(projectId int, todoId int) {
	stmt, err := s.db.Prepare("INSERT INTO project_todos (projectId, todoId) VALUES (?, ?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(projectId, todoId)
	if err != nil {
		panic(err)
	}
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

func (s *sqliteHandler) GetProgress(projectId int) int {
	log.Println(projectId)
	rows, err := s.db.Query(`
		SELECT 
		COUNT(*) AS total_count, 
		COUNT(CASE WHEN completed = 1 THEN 1 ELSE NULL END) AS completed_count 
		FROM todos 
		WHERE id IN (
			SELECT todoId FROM project_todos WHERE projectId=?
		)`, projectId)
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

	statement, _ = database.Prepare(
		`CREATE TABLE IF NOT EXISTS project_todos (
			id          INTEGER  PRIMARY KEY AUTOINCREMENT,
			projectId   INTEGER,
			todoId      INTEGER,
			FOREIGN KEY (projectId) REFERENCES projects(id),
			FOREIGN KEY (todoId) REFERENCES todos(id)
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

	statement, _ = database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			userId    INTEGER,
			name      TEXT,
			completed BOOLEAN,
			createdAt STRING
		);
		CREATE INDEX IF NOT EXISTS userIdIndexOnTodos ON todos (
			userId ASC
		)`)
	statement.Exec()

	return &sqliteHandler{db: database}
}
