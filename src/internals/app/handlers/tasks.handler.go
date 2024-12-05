package handlers

import (
	"fmt"
	"net/http"

	"github.com/Swiffy12/taskify/src/internals/app/services"
)

type TasksHandler struct {
	service *services.TasksService
}

func NewTasksHandler(service *services.TasksService) *TasksHandler {
	tasksHandler := new(TasksHandler)
	tasksHandler.service = service
	return tasksHandler
}

func (th *TasksHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get all tasks")
}

func (th *TasksHandler) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)

}
