package httpapp

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
)

func decodeJSON[T any](r *http.Request, target *T) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return fmt.Errorf("expected application/json")
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var value *T
	if err := decoder.Decode(&value); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("expected an object")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("expected one JSON object")
	}
	*target = *value
	return nil
}
