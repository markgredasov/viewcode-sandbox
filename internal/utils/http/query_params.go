package utils_http

import (
	"net/http"
)

func GetStringQueryParam(r *http.Request, key string) *string {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	return &param
}
