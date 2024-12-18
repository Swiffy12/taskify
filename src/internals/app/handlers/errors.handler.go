package handlers

import (
	"errors"
	"net/http"
)

func WrapErrorNotFound(w http.ResponseWriter, r *http.Request) {
	WrapErrorWithStatus(w, errors.New("метод не реализован"), http.StatusNotFound)
}

func WrapErrorBadRequest(w http.ResponseWriter, err error) {
	WrapErrorWithStatus(w, err, http.StatusBadRequest)
}

func WrapErrorUnauthorized(w http.ResponseWriter, err error) {
	WrapErrorWithStatus(w, err, http.StatusUnauthorized)
}

func WrapErrorInternalServerError(w http.ResponseWriter) {
	WrapErrorWithStatus(w, errors.New("ошибка внутреннего сервера"), http.StatusInternalServerError)
}
