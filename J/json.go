package J

import (
	"bytes"
	"encoding/json"
	"io"
)

type Json map[string]any

type Array []Json

func (j *Json) Reader() (io.Reader, error) {
	b, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (a *Array) Reader() (io.Reader, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
