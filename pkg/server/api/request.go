package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetJsonRequestBody(r *http.Request) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := DecodeJsonRequestBody(r, &m)
	return m, err
}

func DecodeJsonRequestBody(r *http.Request, v interface{}) error {

	defer func() {
		_, _ = io.Copy(io.Discard, r.Body)
	}()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil && err != io.EOF {
		return fmt.Errorf(`invalid request body: %v`, err)
	}

	return nil
}
