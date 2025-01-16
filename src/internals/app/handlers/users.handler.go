package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Swiffy12/taskify/src/internals/app/services"
	"github.com/gorilla/mux"
)

type UsersHandler struct {
	service *services.UsersService
}

func NewUsersHandler(service *services.UsersService) *UsersHandler {
	usersHandler := new(UsersHandler)
	usersHandler.service = service
	return usersHandler
}

func (usersHandler *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	fullname := vars.Get("fullname")
	rank := vars.Get("rank")

	listUsers := usersHandler.service.GetAllUsers(fullname, rank)
	WrapOK(w, listUsers)
}

func (usersHandler *UsersHandler) GetOneUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		WrapErrorBadRequest(w, errors.New("пропущено id пользователя"))
		return
	}

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WrapErrorBadRequest(w, errors.New("недопустимый формат ввода id"))
		return
	}

	user, err := usersHandler.service.GetOneUser(id)
	if err != nil {
		WrapErrorNotFound(w, err)
		return
	}
	WrapOK(w, user)
}
