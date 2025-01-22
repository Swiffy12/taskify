package storages

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

var createdTaskStatus uint8 = 1

type TasksStorage struct {
	databasePool *pgxpool.Pool
}

func NewTasksStorage(pool *pgxpool.Pool) *TasksStorage {
	tasksStorage := new(TasksStorage)
	tasksStorage.databasePool = pool
	return tasksStorage
}

func (storage TasksStorage) CreateTask(userId string, taskBody models.CreateTaskRequestDTO) (*models.Task, error) {
	query := `INSERT INTO tasks (title, description, priority_id, status_id, creator_id, assigned_id) values ($1, $2, $3, $4, $5, $6)
	RETURNING *`

	var createdUser models.Task
	err := pgxscan.Get(context.Background(), storage.databasePool, &createdUser, query, taskBody.Title, taskBody.Description, taskBody.PriorityId, createdTaskStatus, userId, taskBody.AssignedId)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	return &createdUser, nil
}

func (storage TasksStorage) GetTasksWithFilter(queryParams models.GetTasksRequestDTO) ([]models.Task, error) {
	query := "SELECT * FROM tasks WHERE 1=1"

	var receivedTasks []models.Task
	placeholderNumber := 1
	args := make([]any, 0)

	if queryParams.Id != 0 {
		query += fmt.Sprintf(" AND id = $%d", placeholderNumber)
		args = append(args, queryParams.Id)
		placeholderNumber++
	}

	if queryParams.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", placeholderNumber)
		args = append(args, fmt.Sprintf("%%%s%%", queryParams.Title))
		placeholderNumber++
	}

	if queryParams.CreatorId != 0 {
		query += fmt.Sprintf(" AND creator_id = $%d", placeholderNumber)
		args = append(args, queryParams.CreatorId)
		placeholderNumber++
	}

	if queryParams.AssignedId != 0 {
		query += fmt.Sprintf(" AND assigned_id = $%d", placeholderNumber)
		args = append(args, queryParams.AssignedId)
		placeholderNumber++
	}

	err := pgxscan.Select(context.Background(), storage.databasePool, &receivedTasks, query, args...)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	return receivedTasks, nil
}

func (storage TasksStorage) GetOneTask(id uint64) (*models.Task, error) {
	query := "SELECT * FROM tasks WHERE id = $1"
	var receivedTask models.Task

	err := pgxscan.Get(context.Background(), storage.databasePool, &receivedTask, query, id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			logrus.Errorln(err)
		}
		return nil, err
	}
	return &receivedTask, nil
}

func (storage TasksStorage) DeleteOneTask(id uint64) error {
	query := "DELETE FROM tasks WHERE id = $1"

	_, err := storage.databasePool.Exec(context.Background(), query, id)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}

func (storage TasksStorage) UpdateOneTask(id uint64, taskBody models.UpdateTaskRequestDTO) (*models.Task, error) {
	setValues := make([]string, 0)
	placeholderNumber := 1
	args := make([]any, 0)

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

	if taskBody.PriorityId != 0 {
		setValues = append(setValues, fmt.Sprintf("priority_id = $%d", placeholderNumber))
		args = append(args, taskBody.PriorityId)
		placeholderNumber++
	}

	if taskBody.StatusId != 0 {
		setValues = append(setValues, fmt.Sprintf("status_id = $%d", placeholderNumber))
		args = append(args, taskBody.StatusId)
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
	updateTime := time.Now()
	query := fmt.Sprintf("UPDATE tasks SET %s, updated_at = $%d WHERE id = $%d RETURNING *", setQuery, placeholderNumber, placeholderNumber+1)
	args = append(args, updateTime, id)
	var updatedTask models.Task

	err := pgxscan.Get(context.Background(), storage.databasePool, &updatedTask, query, args...)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			logrus.Errorln(err)
		}
		return nil, err
	}
	return &updatedTask, nil
}
