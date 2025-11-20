package server

import (
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

func RecoverMiddleware(logger *zap.Logger) Middleware {
    return func(next http.Handler) http.Handler {

        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            defer func() {
                if rec := recover(); rec != nil {

                    logger.Error("panic recovered",
                        zap.Any("error", rec),
                        zap.String("path", r.URL.Path),
                        zap.ByteString("stack", debug.Stack()),
                    )
                }
            }()

            next.ServeHTTP(w, r)
        })
    }
}

