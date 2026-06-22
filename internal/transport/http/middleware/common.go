package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	l "github.com/Zerferrous/Task-Manager-API/internal/logger"
	"github.com/Zerferrous/Task-Manager-API/internal/transport/http/handler/response"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(l *l.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			log := l.With(
				slog.String("request_id", requestID),
				slog.String("url", r.URL.String()),
				slog.String("method", r.Method),
			)

			ctx := context.WithValue(r.Context(), "logger", log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := l.FromContext(ctx)
			responseHandler := response.NewHandler(logger, w)

			defer func() {
				if p := recover(); p != nil {
					logger.Error(
						"panic recovered",
						slog.Any("panic", p),
						slog.String("stack", string(debug.Stack())),
					)
					responseHandler.Panic(p, "caught panic during handling HTTP request")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := l.FromContext(ctx)
			rw := response.NewWriter(w)

			start := time.Now()
			logger.Debug("Handling HTTP request", slog.Time("time", start.UTC()))

			next.ServeHTTP(rw, r.WithContext(ctx))

			logger.Debug("Finished HTTP request",
				slog.Int("status_code", rw.GetStatusCodeOrPanic()),
				slog.Time("time", time.Now().UTC()),
				slog.Duration("latency", time.Since(start)),
			)
		})
	}
}
