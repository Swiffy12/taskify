package postgresql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Swiffy12/taskify/internal/http-server/models"
	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

const (
	statusCreated = "created"
)

func New(host, port, user, password, dbName string) (*Storage, error) {
	const op = "storage.postgresql.New"

	storagePath := newStoragePath(host, port, user, password, dbName)

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS tasks(
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		description TEXT,
		status VARCHAR(30) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtIdx, err := db.Prepare(`
		CREATE INDEX IF NOT EXISTS idx_title ON tasks(title);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmtIdx.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) CreateTask(title, description string) (int, error) {
	const op = "storage.postgresql.CreateTask"

	stmt, err := s.DB.Prepare("INSERT INTO tasks (title, description, status, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	createdAt := time.Now()
	updatedAt := createdAt

	var id int
	err = stmt.QueryRow(title, description, statusCreated, createdAt, updatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetAllTasks(title string) ([]models.Task, error) {
	const op = "storage.postgresql.GetAllTasks"
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks"
	args := []any{}

	if title != "" {
		query += " WHERE title ILIKE $1"
		args = append(args, "%"+title+"%")
	}

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var rows *sql.Rows
	if len(args) > 0 {
		rows, err = stmt.Query(args...)
	} else {
		rows, err = stmt.Query()
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (s *Storage) GetTask(id int) (models.Task, error) {
	const op = "storage.postgresql.GetTask"

	stmt, err := s.DB.Prepare("SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1")
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	var task models.Task
	row := stmt.QueryRow(id)
	err = row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (s *Storage) DeleteTask(id int) (int, error) {
	const op = "storage.postgresql.DeleteTask"

	stmt, err := s.DB.Prepare("DELETE FROM tasks WHERE id = $1 RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRow(id)
	var respId int
	err = row.Scan(&respId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return respId, nil
}

func (s *Storage) UpdateTask(id int, title, description, status string) (models.Task, error) {
	const op = "storage.postgresql.UpdateTask"

	builder := strings.Builder{}
	builder.WriteString("UPDATE tasks SET updated_at = $1")
	updatedAt := time.Now()
	args := []any{updatedAt}
	counter := 2

	if title != "" {
		builder.WriteString(fmt.Sprintf(", title = $%d", counter))
		args = append(args, title)
		counter++
	}

	if description != "" {
		builder.WriteString(fmt.Sprintf(", description = $%d", counter))
		args = append(args, description)
		counter++
	}

	if status != "" {
		builder.WriteString(fmt.Sprintf(", status = $%d", counter))
		args = append(args, status)
		counter++
	}
	builder.WriteString(fmt.Sprintf(" WHERE id = $%d RETURNING *;", counter))

	stmt, err := s.DB.Prepare(builder.String())
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}
	args = append(args, id)

	row := stmt.QueryRow(args...)
	var task models.Task
	err = row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func newStoragePath(host, port, user, password, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)
}
