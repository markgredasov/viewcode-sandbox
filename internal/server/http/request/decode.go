package request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err) //TODO: add core invalid request error
	}
	return nil
}
