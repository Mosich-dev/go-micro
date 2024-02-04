package main

import (
	"github.com/Mosich-dev/JSONProc"
	"log"
	"net/http"
)

func Broker(w http.ResponseWriter, r *http.Request) {
	payload := JSONProc.JsonResponse{
		Error:   false,
		Message: "Hit",
	}
	err := JSONProc.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Fatalf("Broker Failed:\n%v", err)
	}
}
