package util

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func NewBadRequestError(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func NewInternalError(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func HandleError(c *gin.Context, err error, handler string) {
	if appErr, ok := err.(*AppError); ok {
		if handler != "" {
			if appErr.Code >= http.StatusInternalServerError {
				slog.Error("server error", "handler", handler, "status", appErr.Code, "error", appErr.Message)
			} else {
				slog.Warn("client error", "handler", handler, "status", appErr.Code, "error", appErr.Message)
			}
		}
		SendErrorResponse(c, appErr.Code, appErr.Message)
	} else {
		if handler != "" {
			slog.Error("unexpected error", "handler", handler, "error", err)
		}
		SendErrorResponse(c, http.StatusInternalServerError, "unexpected error")
	}
}
