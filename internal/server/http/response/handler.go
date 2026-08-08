package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"go.uber.org/zap"
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
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, errs.ErrNoTestsProvided):
		statusCode = http.StatusNotFound
		logFunc = h.log.Error
	case errors.Is(err, errs.ErrEmptyCode):
		fallthrough
	case errors.Is(err, errs.ErrInvalidLanguage):
		fallthrough
	case errors.Is(err, errs.ErrTooLargeCode):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Info
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	errorResponse := ErrorResponse{
		Error:   err,
		Message: msg,
	}

	h.JSONResponse(errorResponse, statusCode)
}
