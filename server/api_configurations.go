package server

import (
	"net/http"

	"dev.azure.com/daimler-mic/content-aggregator/handler"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
)

func ConfigureRoutes(
	s *Server,
	cfg *props.Config,
	contentHandler handler.HTTPHandler,
	logger *zap.Logger,
) {

	s.mux.Handle("/content/aggregate",
		Chain(
			contentHandler,
			RequestLogger(logger),
			BodyLogger(logger, cfg.Logging.BodyLogging),
			RecoverMiddleware(logger),
		),
	)

	s.mux.Handle("/healthz",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		}),
	)
}
