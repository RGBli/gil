package binding

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONBinding struct{}

func (JSONBinding) Name() string {
	return "json"
}

func (JSONBinding) Bind(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	return decodeJSON(req.Body, obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
