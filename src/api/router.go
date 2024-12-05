package api

import (
	"github.com/Swiffy12/taskify/src/internals/app/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(tasksHandler *handlers.TasksHandler, usersHandler *handlers.UsersHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", tasksHandler.GetAllTasks).Methods("GET")
	router.HandleFunc("/tasks", tasksHandler.Create).Methods("POST")

	router.HandleFunc("/users", usersHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", usersHandler.Create).Methods("POST")

	router.NotFoundHandler = router.NewRoute().HandlerFunc(handlers.NotFound).GetHandler()
	return router
}
