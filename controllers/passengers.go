package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/abdullahi/go-drive/models"
	u "github.com/abdullahi/go-drive/utils"
)

var CreatePassenger = func(w http.ResponseWriter, r *http.Request) {

	passenger := &models.Passenger{}
	err := json.NewDecoder(r.Body).Decode(passenger)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := passenger.Create()
	u.Respond(w, resp)
}

var AuthPassenger = func(w http.ResponseWriter, r *http.Request) {

	passenger := &models.Passenger{}
	err := json.NewDecoder(r.Body).Decode(passenger) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.PassengerLogin(passenger.Email, passenger.Password)
	u.Respond(w, resp)
}
