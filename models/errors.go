package models

import (
	"encoding/json"
	"net/http"
)

type CustomError struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
	Data   any    `json:"error"`
}

func (c CustomError) Error() string {
	return c.Msg
}
func (c CustomError) SendError(w http.ResponseWriter) {
	w.WriteHeader(c.Status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

var (
	ErrUserNotFound = CustomError{
		Msg:    "user not found",
		Status: http.StatusNotFound,
	}
	ErrUserAlreadyExists = CustomError{
		Msg:    "user already exists",
		Status: http.StatusConflict,
	}
	ErrUsernameOrPasswordIncorrect = CustomError{
		Msg:    "username or password is incorrect",
		Status: http.StatusUnauthorized,
	}
	ErrInternalServerError = CustomError{
		Msg:    "Internal server error",
		Status: http.StatusInternalServerError,
	}
	ErrInvalidRequest = CustomError{
		Msg:    "Invalid request",
		Status: http.StatusBadRequest,
	}
	ErrMsgToSelf = CustomError{
		Msg:    "you can't send message to your self",
		Status: http.StatusConflict,
	}
	ErrUnauthorized = CustomError{
		Msg:    "Unauthorized",
		Status: http.StatusUnauthorized,
	}
	ErrResourceNotFound = CustomError{
		Msg:    "Resource not found",
		Status: http.StatusNotFound,
	}
)
