package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"StatusCode,omitempty"`
	Error      string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("failed %s is a required field", err.Field()))
			break
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("failed %s is a not a valid URL", err.Field()))
			break
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("failed %s is a not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
