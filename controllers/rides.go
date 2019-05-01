package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/abdullahi/go-drive/models"
	"github.com/abdullahi/go-drive/services"
	u "github.com/abdullahi/go-drive/utils"
	"github.com/gorilla/mux"
)

var GetRide = func(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*models.Token)
	id := token.UserId

	ride := models.GetRide(id)

	response := u.Message(true, "Found current ride")
	response["ride"] = ride

	u.Respond(w, response)
}

var CreateRide = func(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*models.Token)
	id := token.UserId
	ride := &models.Ride{}
	ride.PassengerID = id
	err := json.NewDecoder(r.Body).Decode(ride)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	driver, err := services.FindDriver(-96.019415, 41.254311)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	//ride.DriverID = driver.ID
	ride.Driver = *driver
	resp := ride.Create(token)
	u.Respond(w, resp)
}

var UpdateRide = func(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*models.Token)
	status, found := mux.Vars(r)["status"]

	if found == false {
		u.Respond(w, u.Message(false, "Status param not defined"))
		return
	}

	resp := models.UpdateStatus(*token, status)

	u.Respond(w, resp)
}
