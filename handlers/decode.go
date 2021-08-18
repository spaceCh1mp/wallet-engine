package handlers

import (
	"encoding/json"
	"io"
)

// Decode takes a struct that satisfies the io.ReadCloser interface and decodes the data into value.
// value should be a pointer to a struct
func Decode(r io.ReadCloser, value interface{}) error {
	return json.NewDecoder(r).Decode(&value)
}
