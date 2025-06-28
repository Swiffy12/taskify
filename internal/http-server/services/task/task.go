package taskservice

import (
	"fmt"

	"github.com/Swiffy12/taskify/internal/http-server/models"
)

type Storage interface {
	CreateTask(title, description string) (int, error)
	GetTask(id int) (models.Task, error)
	GetAllTasks(title string) ([]models.Task, error)
	DeleteTask(id int) (int, error)
	UpdateTask(id int, title, description, status string) (models.Task, error)
}

type TaskService struct {
	storage Storage
}

func New(storage Storage) *TaskService {
	return &TaskService{storage: storage}
}

func (t *TaskService) CreateTask(title, description string) (int, error) {
	const op = "services.task.CreateTask"

	id, err := t.storage.CreateTask(title, description)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (t *TaskService) GetAllTasks(title string) ([]models.Task, error) {
	const op = "services.task.GetAllTasks"

	tasks, err := t.storage.GetAllTasks(title)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (t *TaskService) GetTask(id int) (models.Task, error) {
	const op = "services.task.GetTask"

	task, err := t.storage.GetTask(id)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (t *TaskService) DeleteTask(id int) (int, error) {
	const op = "services.task.DeleteTask"

	id, err := t.storage.DeleteTask(id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (t *TaskService) UpdateTask(id int, req models.UpdateTaskRequest) (models.Task, error) {
	const op = "services.task.UpdateTask"

	task, err := t.storage.UpdateTask(id, req.Title, req.Description, req.Status)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}
