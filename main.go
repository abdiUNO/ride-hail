package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/abdullahi/go-drive/controllers"
	"github.com/abdullahi/go-drive/services"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	services.SendSMS()

	router.Use(services.JwtAuthentication)

	router.HandleFunc("/drivers/new", controllers.CreateDriver).Methods("POST")
	router.HandleFunc("/drivers/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/drivers/{id}/status", controllers.ChangeStatus).Methods("POST")

	router.HandleFunc("/passengers/new", controllers.CreatePassenger).Methods("POST")
	router.HandleFunc("/passengers/login", controllers.AuthPassenger).Methods("POST")

	router.HandleFunc("/rides", controllers.CreateRide).Methods("POST")
	router.HandleFunc("/rides/current", controllers.GetRide).Methods("GET")
	router.HandleFunc("/rides/current", controllers.UpdateRide).Queries("status", "{status}").Methods("PUT")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		fmt.Print(err)
	}
}
