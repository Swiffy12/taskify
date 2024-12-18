package handlers

import (
	"errors"
	"fmt"
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

func (usersHandler *UsersHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	listUsers := usersHandler.service.GetAllUsers()
	WrapOK(w, listUsers)
}

func (usersHandler *UsersHandler) GetOneUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s", vars["id"])
	if vars["id"] == "" {
		WrapErrorBadRequest(w, errors.New("missing id"))
		return
	}

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WrapErrorBadRequest(w, err)
		return
	}

	user := usersHandler.service.GetOneUser(id)
	WrapOK(w, user)
}
