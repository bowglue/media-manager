package middleware

import (
	"fmt"
	"net/http"
	"time"

	"mms/common/logger"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Log(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			log.Info(fmt.Sprintf("%d %s %s %v", wrapped.statusCode, r.Method, r.URL.Path, time.Since(start)))
		})
	}
}
