package errors

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// helper
func writeJSON(w http.ResponseWriter, logger *zap.Logger, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("failed to write JSON response", zap.Error(err))
	}
}
