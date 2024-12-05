package api

import (
	"github.com/Swiffy12/taskify/src/internals/app/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(tasksHandler *handlers.TaskHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", tasksHandler.GetAllTasks).Methods("GET")
	return router
}
