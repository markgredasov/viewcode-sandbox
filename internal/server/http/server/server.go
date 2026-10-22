package server_http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/middleware"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        *logger.Logger
	middleware []middleware.Middleware
}

func NewHTTPServer(config Config, logger *logger.Logger, middleware ...middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        logger,
		middleware: middleware,
	}
}

func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:              h.config.Address,
		Handler:           mux,
		ReadHeaderTimeout: h.config.ReadHeaderTimeout,
	}

	ch := make(chan error, 1)

	go func() {
		defer func() {
			close(ch)
		}()

		h.log.Warn(
			"http server started",
			logger.String("address", h.config.Address),
		)

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		return fmt.Errorf("listen and serve HTTP server: %w", err)
	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server")
		shutdownCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx),
			h.config.ShutdownTimeout,
		)

		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			err = fmt.Errorf("shutdown HTTP server: %w", err)

			sErr := server.Close()
			if sErr != nil {
				err = fmt.Errorf("close server: %w: %w", sErr, err)
			}

			return err
		}

		h.log.Warn("HTTP server stopped")
	}

	return nil
}
