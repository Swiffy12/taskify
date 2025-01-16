package handlers

import (
	"errors"
	"net/http"
)

func WrapErrorMethodNotFound(w http.ResponseWriter, r *http.Request) {
	WrapErrorWithStatus(w, errors.New("метод не реализован"), http.StatusNotFound)
}

func WrapErrorInternalServerError(w http.ResponseWriter) {
	WrapErrorWithStatus(w, errors.New("ошибка внутреннего сервера"), http.StatusInternalServerError)
}

func WrapErrorBadRequest(w http.ResponseWriter, err error) {
	WrapErrorWithStatus(w, err, http.StatusBadRequest)
}

func WrapErrorUnauthorized(w http.ResponseWriter, err error) {
	WrapErrorWithStatus(w, err, http.StatusUnauthorized)
}

func WrapErrorNotFound(w http.ResponseWriter, err error) {
	WrapErrorWithStatus(w, err, http.StatusNotFound)
}
