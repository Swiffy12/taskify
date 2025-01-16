package storages

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type TasksStorage struct {
	databasePool *pgxpool.Pool
}

func NewTasksStorage(pool *pgxpool.Pool) *TasksStorage {
	tasksStorage := new(TasksStorage)
	tasksStorage.databasePool = pool
	return tasksStorage
}

func (storage TasksStorage) CreateTask(userId string, taskBody models.Task) (models.Task, error) {
	query := `INSERT INTO tasks (title, description, priority, status, creator_id, assigned_id) values ($1, $2, $3, $4, $5, $6)
	RETURNING *`

	var createdUser models.Task
	err := pgxscan.Get(context.Background(), storage.databasePool, &createdUser, query, taskBody.Title, taskBody.Description, taskBody.Priority, 1, userId, taskBody.AssignedId)
	if err != nil {
		logrus.Errorln(err)
		return createdUser, err
	}

	return createdUser, nil
}

func (storage TasksStorage) GetTasksWithFilter(id string, title string, creator string, assigned string) []models.Task {
	query := `SELECT * FROM tasks WHERE 1=1`

	var receivedTasks []models.Task
	placeholderNumber := 1
	args := make([]any, 0)

	if id != "" {
		query += fmt.Sprintf(" AND id = $%d", placeholderNumber)
		args = append(args, id)
		placeholderNumber++
	}

	if title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", placeholderNumber)
		args = append(args, fmt.Sprintf("%%%s%%", title))
		placeholderNumber++
	}

	if creator != "" {
		query += fmt.Sprintf(" AND creator_id = $%d", placeholderNumber)
		args = append(args, creator)
		placeholderNumber++
	}

	if assigned != "" {
		query += fmt.Sprintf(" AND assigned_id = $%d", placeholderNumber)
		args = append(args, assigned)
		placeholderNumber++
	}

	err := pgxscan.Select(context.Background(), storage.databasePool, &receivedTasks, query, args...)
	if err != nil {
		logrus.Errorln(err)
	}

	return receivedTasks
}

func (storage TasksStorage) GetOneTask(id int64) (models.Task, error) {
	query := "SELECT * FROM tasks WHERE id = $1"
	var receivedTask models.Task

	err := pgxscan.Get(context.Background(), storage.databasePool, &receivedTask, query, id)
	if err = errors.Unwrap(errors.Unwrap(err)); err != nil { // Дублирование кода
		if err == pgx.ErrNoRows {
			return receivedTask, errors.New("не удалось найти данную задачу")
		}
		logrus.Errorln(err)
	}
	return receivedTask, nil
}

func (storage TasksStorage) DeleteOneTask(id int64) error {
	query := "DELETE FROM tasks WHERE id = $1"

	_, err := storage.databasePool.Exec(context.Background(), query, id)
	if err != nil {
		logrus.Errorln(err)
		return errors.New("не удалось удалить данную задачу")
	}
	return nil
}

func (storage TasksStorage) UpdateOneTask(id int64, taskBody models.UpdateTaskRequestDTO) models.Task {
	setValues := make([]string, 0)
	placeholderNumber := 1
	args := make([]any, 0)

	if taskBody.Id != 0 {
		setValues = append(setValues, fmt.Sprintf("id = $%d", placeholderNumber))
		args = append(args, taskBody.Id)
		placeholderNumber++
	}

	if taskBody.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title = $%d", placeholderNumber))
		args = append(args, taskBody.Title)
		placeholderNumber++
	}

	if taskBody.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description = $%d", placeholderNumber))
		args = append(args, taskBody.Description)
		placeholderNumber++
	}

	if taskBody.Priority != 0 {
		setValues = append(setValues, fmt.Sprintf("priority = $%d", placeholderNumber))
		args = append(args, taskBody.Priority)
		placeholderNumber++
	}

	if taskBody.Status != 0 {
		setValues = append(setValues, fmt.Sprintf("status = $%d", placeholderNumber))
		args = append(args, taskBody.Status)
		placeholderNumber++
	}

	if taskBody.CreatorId != 0 {
		setValues = append(setValues, fmt.Sprintf("creator_id = $%d", placeholderNumber))
		args = append(args, taskBody.CreatorId)
		placeholderNumber++
	}

	if taskBody.AssignedId != 0 {
		setValues = append(setValues, fmt.Sprintf("assigned_id = $%d", placeholderNumber))
		args = append(args, taskBody.AssignedId)
		placeholderNumber++
	}

	setQuery := strings.Join(setValues, ", ")
	var updatedTask models.Task
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d RETURNING *", setQuery, placeholderNumber)
	args = append(args, id)

	err := pgxscan.Get(context.Background(), storage.databasePool, &updatedTask, query, args...)
	if err != nil {
		logrus.Errorln(err)
	}
	return updatedTask
}
