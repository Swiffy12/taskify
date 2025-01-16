package services

import (
	"errors"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/Swiffy12/taskify/src/internals/app/storages"
)

type TasksService struct {
	storage *storages.TasksStorage
}

func NewTasksService(storage *storages.TasksStorage) *TasksService {
	tasksService := new(TasksService)
	tasksService.storage = storage
	return tasksService
}

func (service TasksService) CreateOneTask(userId string, taskBody models.Task) (models.Task, error) {

	createdUser, err := service.storage.CreateTask(userId, taskBody)
	if err != nil {
		return createdUser, errors.New("не удалось создать задачу")
	}

	return createdUser, err
}

func (service TasksService) GetTasksWithFilter(id string, title string, creator string, assigned string) []models.Task {
	return service.storage.GetTasksWithFilter(id, title, creator, assigned)
}

func (service TasksService) GetOneTask(id int64) (models.Task, error) {
	return service.storage.GetOneTask(id)
}

func (service TasksService) DeleteOneTask(id int64) error {
	task, err := service.storage.GetOneTask(id)
	if err != nil {
		return err
	}

	err = service.storage.DeleteOneTask(task.Id)
	if err != nil {
		return err
	}
	return nil
}

func (service TasksService) UpdateOneTask(id int64, taskBody models.UpdateTaskRequestDTO) (models.Task, error) {
	foundTask, err := service.storage.GetOneTask(id)
	if err != nil {
		return foundTask, err
	}
	return service.storage.UpdateOneTask(id, taskBody), nil
}
