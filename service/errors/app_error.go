package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	StatusCode int                    `json:"statusCode"`        // HTTP status
	Code       string                 `json:"code"`              // machine-readable code
	Message    string                 `json:"message"`           // human-readable msg
	Details    map[string]interface{} `json:"details,omitempty"` // optional metadata
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewInternalError(msg string, details map[string]interface{}) *AppError {
	return &AppError{
		StatusCode: 500,
		Code:       "internal_error",
		Message:    msg,
		Details:    details,
	}
}

func NewBadRequest(msg string, details map[string]interface{}) *AppError {
	return &AppError{
		StatusCode: 400,
		Code:       "bad_request",
		Message:    msg,
		Details:    details,
	}
}

func NewProviderError(provider, msg string, details map[string]interface{}) *AppError {
	return &AppError{
		StatusCode: 502,
		Code:       provider + "_error",
		Message:    msg,
		Details:    details,
	}
}

func NewFeatureError(feature, msg string, details map[string]interface{}) *AppError {
	return &AppError{
		StatusCode: 500,
		Code:       feature + "_feature_error",
		Message:    msg,
		Details:    details,
	}
}

func WriteError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.StatusCode)
		json.NewEncoder(w).Encode(appErr)
		return
	}

	// fallback unexpected error (should be rare)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(NewInternalError(err.Error(), nil))
}
