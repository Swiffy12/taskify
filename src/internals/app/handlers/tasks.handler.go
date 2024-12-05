package handlers

import (
	"fmt"
	"net/http"

	"github.com/Swiffy12/taskify/src/internals/app/services"
)

type TaskHandler struct {
	service services.TasksService
}

func NewTaskHandler(service services.TasksService) *TaskHandler {
	taskHandler := new(TaskHandler)
	taskHandler.service = service
	return taskHandler
}

func (th *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get all tasks")
}
