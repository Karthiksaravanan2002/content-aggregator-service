package errors

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type AppError struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	err 			 error 	`json:"-"` 
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.err
}

func New(status int, code, message string, err error) *AppError {
	return &AppError{
		StatusCode: status,
		Code:       code,
		Message:    message,
		err:        err,
	}
}

func BadRequest(message string, err error) *AppError {
	return New(400, "BAD_REQUEST", message, err)
}

func InternalError(message string, err error) *AppError {
	return New(500, "INTERNAL_ERROR", message, err)
}

func ProviderFailure(provider string, err error) *AppError {
	return New(502, "PROVIDER_FAILURE", fmt.Sprintf("%s provider failed", provider), err)
}

func Unsupported(entity string) *AppError {
	return New(400, "UNSUPPORTED", fmt.Sprintf("%s is not supported", entity), nil)
}

func WriteError(w http.ResponseWriter, logger *zap.Logger, err error) {
    // Not an AppError? → Internal error wrapper
    appErr, ok := err.(*AppError)
    if !ok {
        logger.Error("unhandled error type", zap.Error(err))

        internal := InternalError("internal server error", err)

        writeJSON(w, logger, internal.StatusCode, internal)
        return
    }

    // Log the error (structured)
    logger.Error("request failed",
        zap.String("code", appErr.Code),
        zap.String("message", appErr.Message),
        zap.Error(appErr.err), // internal error if present
    )

    writeJSON(w, logger, appErr.StatusCode, appErr)
}
