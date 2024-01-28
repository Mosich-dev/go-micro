package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const maxBytesToRead = 1048576 // = 1 Megabyte

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	data    any    `json:"data,omitempty"`
}

func (app Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytesToRead)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// check if we receive only a single JSON value if not, return error
	err = dec.Decode(struct{}{})
	if err != io.EOF {
		return errors.New("body most have only a single JSON value")
	}

	return nil
}

func (app Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		fmt.Println("Headers:", headers)
		fmt.Println("Headers[0]:", headers[0])
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		fmt.Println("status:", status)
		fmt.Println("status[0]:", status[0])
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
