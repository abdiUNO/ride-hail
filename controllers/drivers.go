package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/abdullahi/go-drive/models"
	u "github.com/abdullahi/go-drive/utils"
	"github.com/gorilla/mux"
)

type StatusRequest struct {
	Online bool `json:"online"`
}

var CreateDriver = func(w http.ResponseWriter, r *http.Request) {
	driver := &models.Driver{}
	err := json.NewDecoder(r.Body).Decode(driver)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := driver.Create()
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	driver := &models.Driver{}
	err := json.NewDecoder(r.Body).Decode(driver) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(driver.Email, driver.Password)
	u.Respond(w, resp)
}

var ChangeStatus = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	body := &StatusRequest{}

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := models.ChangeStatus(id, body.Online)
	u.Respond(w, resp)
}
