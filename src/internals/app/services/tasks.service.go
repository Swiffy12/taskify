package services

import (
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

func (service TasksService) CreateOneTask(userId string, taskBody models.CreateTaskRequestDTO) (*models.Task, error) {
	return service.storage.CreateTask(userId, taskBody)
}

func (service TasksService) GetTasksWithFilter(queryParams models.GetTasksRequestDTO) ([]models.Task, error) {
	return service.storage.GetTasksWithFilter(queryParams)
}

func (service TasksService) GetOneTask(id uint64) (*models.Task, error) {
	return service.storage.GetOneTask(id)
}

func (service TasksService) DeleteOneTask(id uint64) error {
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

func (service TasksService) UpdateOneTask(id uint64, taskBody models.UpdateTaskRequestDTO) (*models.Task, error) {
	foundTask, err := service.storage.GetOneTask(id)
	if err != nil {
		return nil, err
	}

	return service.storage.UpdateOneTask(foundTask.Id, taskBody)
}
