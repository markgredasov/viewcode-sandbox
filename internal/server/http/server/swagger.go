package server_http

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *HTTPServer) RegisterSwagger() {
	h.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)

	h.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(docs.SwaggerInfo.ReadDoc())) // nolint:errcheck,gosec // todo
		},
	)
}
