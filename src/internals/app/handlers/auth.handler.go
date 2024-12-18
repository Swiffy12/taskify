package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/Swiffy12/taskify/src/internals/app/services"
	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	authHandler := new(AuthHandler)
	authHandler.service = service
	return authHandler
}

func (authHandler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var authModel models.Auth

	err := json.NewDecoder(r.Body).Decode(&authModel)
	if err != nil {
		WrapErrorBadRequest(w, errors.New("ошибка в параметрах запроса"))
		return
	}

	validation, err := govalidator.ValidateStruct(authModel)
	if !validation {
		WrapErrorBadRequest(w, errors.New("недопустимая форма ввода"))
		return
	}
	if err != nil {
		logrus.Errorln(err)
		WrapErrorInternalServerError(w)
		return
	}

	token, err := authHandler.service.Register(authModel)
	if err != nil {
		WrapErrorUnauthorized(w, err)
		return
	}
	WrapOK(w, token)
}

func (authHandler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var authModel models.Auth

	err := json.NewDecoder(r.Body).Decode(&authModel)
	if err != nil {
		WrapErrorBadRequest(w, errors.New("ошибка в параметрах запроса"))
		return
	}

	validationEmail := govalidator.IsEmail(authModel.Email)
	if !validationEmail || authModel.Password == "" {
		WrapErrorBadRequest(w, errors.New("недопустимая форма ввода"))
		return
	}

	token, err := authHandler.service.Login(authModel.Email, authModel.Password)
	if err != nil {
		WrapErrorUnauthorized(w, err)
		return
	}
	WrapOK(w, token)
}
