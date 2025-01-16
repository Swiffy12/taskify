package api

import (
	"github.com/Swiffy12/taskify/src/internals/app/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(
	tasksHandler *handlers.TasksHandler,
	usersHandler *handlers.UsersHandler,
	authHandler *handlers.AuthHandler,
) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	router.HandleFunc("/tasks", tasksHandler.CreateTask).Methods("POST")
	router.HandleFunc("/tasks", tasksHandler.GetTasks).Methods("GET")
	router.HandleFunc("/tasks/{id:.+}", tasksHandler.GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id:.+}", tasksHandler.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id:.+}", tasksHandler.UpdateTask).Methods("PATCH")

	router.HandleFunc("/users", usersHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id:.+}", usersHandler.GetOneUser).Methods("GET")

	router.NotFoundHandler = router.NewRoute().HandlerFunc(handlers.WrapErrorMethodNotFound).GetHandler()
	return router
}
