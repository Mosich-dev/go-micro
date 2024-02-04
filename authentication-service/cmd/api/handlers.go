package main

import (
	"errors"
	"fmt"
	"github.com/Mosich-dev/JSONProc"
	"github.com/Mosich-dev/go-micro/authentication-service/data"
	"net/http"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := JSONProc.ReadJSON(w, r, &reqPayload)
	if err != nil {
		JSONProc.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	user, err := data.GetByEmail(reqPayload.Email)
	if err != nil {
		JSONProc.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	valid, err := user.PasswordMatches(reqPayload.Password)
	if err != nil || !valid {
		JSONProc.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	resPayload := JSONProc.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user: %s", user.Email),
		Data:    user,
	}

	err = JSONProc.WriteJSON(w, http.StatusAccepted, resPayload)
	if err != nil {
		JSONProc.ErrorJSON(w, err, http.StatusNotFound)
	}
}
