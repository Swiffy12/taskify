package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

const (
	ResultOK    = "OK"
	ResultError = "Error"
)

func OK(data any) Response {
	return Response{
		Result: ResultOK,
		Data:   data,
	}
}

func Error(msg string) Response {
	return Response{
		Result: ResultError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.StructField()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Result: ResultError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
