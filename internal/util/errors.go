package util

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AppError struct {
	Code    int
	Message string
	cause   error
	stack   string
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

func NewInternalError(msg string, cause ...error) *AppError {
	var c error
	if len(cause) > 0 {
		c = cause[0]
	}
	return &AppError{Code: http.StatusInternalServerError, Message: msg, cause: c, stack: string(debug.Stack())}
}

func HandleError(c *gin.Context, err error, log *zap.SugaredLogger) {
	if appErr, ok := err.(*AppError); ok {
		if log != nil {
			if appErr.Code >= http.StatusInternalServerError {
				log.Errorw("server error", "status", appErr.Code, "error", appErr.Message, "cause", appErr.cause, "request_id", GetRequestIDFromContext(c), "stacktrace", appErr.stack)
			} else {
				log.Warnw("client error", "status", appErr.Code, "error", appErr.Message, "request_id", GetRequestIDFromContext(c))
			}
		}
		SendErrorResponse(c, appErr.Code, appErr.Message)
	} else {
		if log != nil {
			log.Errorw("unexpected error", "error", err, "request_id", GetRequestIDFromContext(c), "stacktrace", string(debug.Stack()))
		}
		SendErrorResponse(c, http.StatusInternalServerError, "unexpected error")
	}
}
