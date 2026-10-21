package httputils

import (
	"net/http"
	"strings"
)

func GetStringQueryParam(r *http.Request, key string) *string {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	return &param
}

func GetStringArrayQueryParam(r *http.Request, key string) *[]string {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	strs := strings.Split(param, ",")
	return &strs
}
