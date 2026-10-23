package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/platform/logger"
)

type HTTPResponseHandler struct {
	log *logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(l *logger.Logger, w http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: l,
		w:   w,
	}
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(statusCode)

	if statusCode != http.StatusNoContent {
		if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
			h.log.Error("write HTTP response", logger.Err(err))
		}
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...any)
	)

	switch {
	case errors.Is(err, errs.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Info
	case errors.Is(err, errs.ErrAlreadyExists):
		statusCode = http.StatusConflict
		logFunc = h.log.Error
	case errors.Is(err, errs.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Info
	case errors.Is(err, errs.ErrInvalidRequest):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Error
	case errors.Is(err, errs.ErrNotImplemented):
		statusCode = http.StatusNotImplemented
		logFunc = h.log.Error
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, logger.Err(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, logger.Err(err), logger.String("stack", string(debug.Stack())))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	h.JSONResponse(NewErrorResponse(err, msg), statusCode)
}
