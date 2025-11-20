package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.Status = code
	r.ResponseWriter.WriteHeader(code)
}

func RequestLogger(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, Status: 200}

			reqID := r.Header.Get("X-Request-ID")

			logger.Info("request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remoteAddr", r.RemoteAddr),
				zap.String("requestId", reqID),
			)

			next.ServeHTTP(rec, r)

			latency := time.Since(start)

			logger.Info("request finished",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rec.Status),
				zap.String("remoteAddr", r.RemoteAddr),
				zap.String("requestId", reqID),
				zap.Duration("latency", latency),
			)
		})
	}
}
