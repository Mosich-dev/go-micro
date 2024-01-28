package main

import (
	"log"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit",
	}
	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Fatalf("Broker Failed:\n%v", err)
	}
}
