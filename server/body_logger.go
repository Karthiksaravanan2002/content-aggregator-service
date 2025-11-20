package server

import (
	"bytes"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// BodyLogging mode: none | all | errors
func BodyLogger(logger *zap.Logger, mode string) Middleware {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if mode == "none" {
				next.ServeHTTP(w, r)
				return
			}

			var reqBody []byte
			if r.Body != nil {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewReader(reqBody))
			}

			buf := &bytes.Buffer{}
			rec := &bodyRecorder{
				ResponseWriter: w,
				body:           buf,
				status:         200,
			}

			next.ServeHTTP(rec, r)

			if mode == "all" || (mode == "errors" && rec.status >= 400) {
				logger.Info("body logging",
					zap.String("path", r.URL.Path),
					zap.Int("status", rec.status),
					zap.ByteString("req", reqBody),
					zap.ByteString("res", buf.Bytes()),
				)
			}

			_, _ = w.Write(buf.Bytes())
		})
	}
}

type bodyRecorder struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (r *bodyRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *bodyRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return len(b), nil
}
