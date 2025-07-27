package resp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

const (
	ResultOK    = "OK"
	ResultError = "error"
)

func OK(w http.ResponseWriter, r *http.Request, status int, data any) {
	resp := Response{
		Result: ResultOK,
		Data:   data,
	}
	w.WriteHeader(status)
	render.JSON(w, r, resp)
}

func Error(w http.ResponseWriter, r *http.Request, status int, msg string) {
	resp := Response{
		Result: ResultError,
		Error:  msg,
	}

	w.WriteHeader(status)
	render.JSON(w, r, resp)
}

func ValidationError(w http.ResponseWriter, r *http.Request, status int, errs validator.ValidationErrors) {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.StructField()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	resp := Response{
		Result: ResultError,
		Error:  strings.Join(errMsgs, ", "),
	}

	w.WriteHeader(status)
	render.JSON(w, r, resp)
}
