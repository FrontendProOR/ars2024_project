package main

import (
	"net/http"
	"project/handlers"

	//"project/model"
	"project/repositories"
	"project/services"

	"github.com/gorilla/mux"
)

func main() {
	// Inicijalizacija repozitorijuma, servisa i handlera
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
	handler := handlers.NewConfigHandler(service)

	// Kreiranje novog router-a
	router := mux.NewRouter()

	// Registracija ruta za ConfigHandler
	router.HandleFunc("/configs", handler.Add).Methods("POST")                       // Ruta za dodavanje konfiguracije
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")       // Ruta za pregled konfiguracije po imenu i verziji
	router.HandleFunc("/configs/{name}/{version}", handler.Delete).Methods("DELETE") // Ruta za brisanje konfiguracije po imenu i verziji

	// Pokretanje HTTP servera
	http.ListenAndServe("0.0.0.0:8000", router)
}
