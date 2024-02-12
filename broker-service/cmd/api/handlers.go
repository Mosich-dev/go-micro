package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Mosich-dev/JSONProc"
	"github.com/Mosich-dev/go-micro/broker-service/cmd/types"
	"io"
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

func HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload types.RequestPayload

	err := JSONProc.ReadJSON(w, r, &reqPayload)
	if err != nil {
		JSONProc.ErrorJSON(w, err)
		return
	}

	switch reqPayload.Action {
	case "auth":
		authenticate(w, reqPayload.Auth)
	default:
		JSONProc.ErrorJSON(w, errors.New("unknown action"))
	}
}

func authenticate(w http.ResponseWriter, a types.AuthPayload) {
	data, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(data))
	if err != nil {
		JSONProc.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		JSONProc.ErrorJSON(w, err)
		return
	}

	if res.StatusCode == http.StatusUnauthorized {
		JSONProc.ErrorJSON(w, errors.New("invalid credentials")) // maybe return code 401 (Unauthorized)
		return
	} else if res.StatusCode != http.StatusAccepted {
		JSONProc.ErrorJSON(w, errors.New("error: auth service"))
		return
	}
	var jsonFromAuthService JSONProc.JsonResponse
	err = json.NewDecoder(res.Body).Decode(&jsonFromAuthService)
	if err != nil {
		JSONProc.ErrorJSON(w, err)
		return
	}
	if jsonFromAuthService.Error {
		JSONProc.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	payload := JSONProc.JsonResponse{
		Error:   false,
		Message: "Authorized",
		Data:    jsonFromAuthService.Data,
	}
	err = JSONProc.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		JSONProc.ErrorJSON(w, err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
			JSONProc.ErrorJSON(w, err)
		}
	}(res.Body)
}
