package data

import (
	"encoding/json"
	"io"
)

// decodes json to product struct
// takes io reader as a parameter (Response Body) and returns an error
func FromJSON(i interface{}, r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(i)
}

// encodes products struct to json
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w) //creates new encoder with io writer input
	return e.Encode(i)
}
