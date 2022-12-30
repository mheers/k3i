package helpers

import (
	"encoding/json"
	"io"
)

func ReadJSON(body io.ReadCloser, v interface{}) error {
	// reads from io.ReadCloser and unmarshals the json into the given interface
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func ReadBytes(body io.ReadCloser) ([]byte, error) {
	// reads from io.ReadCloser
	return io.ReadAll(body)
}
