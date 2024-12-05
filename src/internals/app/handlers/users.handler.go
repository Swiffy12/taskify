package handlers

import (
	"net/http"

	"github.com/Swiffy12/taskify/src/internals/app/services"
)

type UsersHandler struct {
	service *services.UsersService
}

func NewUsersHandler(service *services.UsersService) *UsersHandler {
	usersHandler := new(UsersHandler)
	usersHandler.service = service
	return usersHandler
}

func (handler *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (handler *UsersHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	listUsers := handler.service.GetAllUsers()

	var m = map[string]any{
		"result": "OK",
		"data":   listUsers,
	}

	WrapOK(w, m)
}
