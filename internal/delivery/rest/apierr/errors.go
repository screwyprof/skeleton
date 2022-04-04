package apierr

import (
	"fmt"
	"net/http"
)

const (
	BadRequest             = 400
	Conflict               = 409
	NotFound               = 404
	InvalidLoginOrPassword = 401

	InternalServerError = 500
)

var codeToStatus = map[int]int{
	BadRequest:             http.StatusBadRequest,
	Conflict:               http.StatusConflict,
	NotFound:               http.StatusNotFound,
	InvalidLoginOrPassword: http.StatusBadRequest,

	InternalServerError: http.StatusInternalServerError,
}

var codeToMsg = map[int]string{
	BadRequest:             "Bad Request",
	NotFound:               "Not Found",
	Conflict:               "Conflict",
	InternalServerError:    "Internal Server Error",
	InvalidLoginOrPassword: "Invalid login or password",
}

func Wrap(cause error, code int, extra map[string]interface{}) error {
	return &APIError{
		ErrStatus: status(code),
		ErrCode:   code,
		ErrMsg:    msg(code),
		ErrExtra:  extra,
		ErrCause:  cause,
	}
}

func status(code int) int {
	status := http.StatusInternalServerError
	if s, ok := codeToStatus[code]; ok {
		status = s
	}

	return status
}

func msg(code int) string {
	msg := "Internal Server Error"
	if m, ok := codeToMsg[code]; ok {
		msg = m
	}

	return msg
}

type APIError struct {
	ErrStatus int                    `json:"-"`
	ErrCode   int                    `json:"code"`
	ErrMsg    string                 `json:"msg"`
	ErrExtra  map[string]interface{} `json:"extra,omitempty"`
	ErrCause  error                  `json:"-"`
}

func (a *APIError) Code() int {
	return a.ErrCode
}

func (a *APIError) Status() int {
	return a.ErrStatus
}

func (a *APIError) Message() string {
	return a.ErrMsg
}

func (a *APIError) Extra() map[string]interface{} {
	return a.ErrExtra
}

func (a *APIError) Cause() error {
	return a.ErrCause
}

func (a *APIError) Error() string {
	return fmt.Sprintf("[%d] %s", a.Code(), a.ErrMsg)
}
