package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		wrappedWriter := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		h.ServeHTTP(wrappedWriter, r)

		timeStop := time.Since(timeStart)

		slog.Info("Запрос",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", timeStop),
			slog.Int("status", wrappedWriter.statusCode),
		)
	})
}
