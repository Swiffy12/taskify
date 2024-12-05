package services

import (
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
