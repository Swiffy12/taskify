package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/Swiffy12/taskify/src/internals/app/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
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
	queryParams := models.GetUsersRequestDTO{
		FullName: vars.Get("fullname"),
		Rank:     vars.Get("rank"),
	}

	listUsers, err := usersHandler.service.GetUsersWithFilter(queryParams)
	if err != nil {
		WrapErrorInternalServerError(w)
		return
	}
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
		if errors.Is(err, pgx.ErrNoRows) {
			WrapErrorNotFound(w, errors.New("пользователь не найден"))
			return
		}
		WrapErrorInternalServerError(w)
		return
	}
	WrapOK(w, user)
}
