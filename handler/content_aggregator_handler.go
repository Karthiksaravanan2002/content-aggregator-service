package handler

import (
	"context"
	"encoding/json"

	"net/http"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service"
	"dev.azure.com/daimler-mic/content-aggregator/service/helper"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
)

type ctxKey string

const headersKey ctxKey = "headers"

type ContentHandler struct {
	contentService service.ContentService
	logger         *zap.Logger
	props          *props.Config
}

func NewContentHandler(
	contentService service.ContentService,
	logger *zap.Logger,
	props *props.Config,
) *ContentHandler {

	return &ContentHandler{
		contentService: contentService,
		logger:         logger.With(zap.String("handler", "content")),
		props:          props,
	}
}

var _ HTTPHandler = (*ContentHandler)(nil)

func (h *ContentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request JSON
	var aggReq *models.AggregateRequest
	if err := json.NewDecoder(r.Body).Decode(&aggReq); err != nil {
		h.logger.Warn("invalid JSON", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validation → at least 1 provider must be present
	if len(aggReq.Providers) == 0 {
		http.Error(w, "at least one provider must be supplied", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), headersKey, r.Header)
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.props.Server.Timeout)*time.Second)
	defer cancel()

	// Execute content aggregation
	resp := h.contentService.Aggregate(ctx, aggReq)
	w.WriteHeader(helper.SelectRespStatusCode(resp))
	_ = json.NewEncoder(w).Encode(resp)
}
