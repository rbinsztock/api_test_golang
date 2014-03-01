package main

import (
	"encoding/json"
)

// An Encoder implements an encoding format of values to be sent as response to
// requests on the API endpoints.
type Encoder interface {
	Encode(v ...interface{}) (string, error)
}

// Because `panic`s are caught by martini's Recovery handler, it can be used
// to return server-side errors (500). Some helpful text message should probably
// be sent, although not the technical error (which is printed in the log).
func Must(data string, err error) string {
	if err != nil {
		panic(err)
	}
	return data
}

type jsonEncoder struct{}

// jsonEncoder is an Encoder that produces JSON-formatted responses.
func (_ jsonEncoder) Encode(v ...interface{}) (string, error) {
	var data interface{} = v
	if v == nil {
		// So that empty results produces `[]` and not `null`
		data = []interface{}{}
	} else if len(v) == 1 {
		data = v[0]
	}
	b, err := json.Marshal(data)
	return string(b), err
}
