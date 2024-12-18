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

	router.HandleFunc("/tasks", tasksHandler.GetAllTasks).Methods("GET")
	router.HandleFunc("/tasks", tasksHandler.Create).Methods("POST")

	router.HandleFunc("/users", usersHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id:.+}", usersHandler.GetOneUser).Methods("GET")

	router.NotFoundHandler = router.NewRoute().HandlerFunc(handlers.WrapErrorNotFound).GetHandler()
	return router
}
