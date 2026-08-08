package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
	"go.uber.org/zap"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		const RequestIDHeader = "x-request-id"
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(RequestIDHeader, requestID)
			w.Header().Set(RequestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(l *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		const RequestIDHeader = "x-request-id"
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)

			logger := l.With(
				zap.String("Request-ID", requestID),
				zap.String("URL", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "logger", logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()
			logger := logger.FromContext(ctx)
			responseHandler := response.NewHTTPResponseHandler(logger, w)

			defer func() {
				if r := recover(); r != nil {
					responseHandler.PanicResponse(
						r,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()
			logger := logger.FromContext(ctx)

			writer := response.NewResponseWriter(w)

			start := time.Now().UTC()
			logger.Debug(
				">>> incoming HTTP request",
				zap.Time("time", start),
			)

			next.ServeHTTP(writer, r)

			logger.Debug(
				"<<< done HTTP request",
				zap.Int("status-code", writer.GetStatusCode()),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}
