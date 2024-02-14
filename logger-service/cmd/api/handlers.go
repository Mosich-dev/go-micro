package main

import (
	"github.com/Mosich-dev/JSONProc"
	"github.com/Mosich-dev/go-micro/logger-service/cmd/types"
	"github.com/Mosich-dev/go-micro/logger-service/data"
	"log"
	"net/http"
)

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var reqPayload types.LoggerPayload
	err := JSONProc.ReadJSON(w, r, &reqPayload)
	if err != nil {
		log.Println("error reading json: ", err)
		JSONProc.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	event := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}
	err = app.Models.Insert(event)
	if err != nil {
		log.Println("error inserting log into db: ", err)
		JSONProc.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	res := JSONProc.JsonResponse{
		Error:   false,
		Message: "logged",
		Data:    nil,
	}
	err = JSONProc.WriteJSON(w, http.StatusCreated, res)
	if err != nil {
		log.Println("error writing json: ", err)
		JSONProc.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
}
