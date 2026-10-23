package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/response"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		const requestIDHeader = "x-request-id"
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

func Logger(l *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		const requestIDHeader = "x-request-id"
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			requestLogger := l.With(
				logger.String("request_id", requestID),
				logger.String("method", r.Method),
				logger.String("path", r.URL.String()),
			)

			ctx := logger.InContext(r.Context(), requestLogger)

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
			log := logger.FromContext(ctx)
			responseHandler := response.NewHTTPResponseHandler(log, w)

			defer func() {
				if r := recover(); r != nil {
					responseHandler.PanicResponse(r, "during handle http request got unexpected panic")
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
			log := logger.FromContext(ctx)

			writer := response.NewResponseWriter(w)

			start := time.Now().UTC()

			next.ServeHTTP(writer, r)

			log.Debug(
				"http request",
				logger.Int("status_code", writer.GetStatusCode()),
				logger.String("latency", time.Since(start).String()),
			)
		})
	}
}

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:8080": {},
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
