package main

import (
	"net/http"
	"project/handlers"
	"project/repositories"
	"project/services"

	"github.com/gorilla/mux"
)

func main() {
	// Inicijalizacija repozitorijuma, servisa i handlera za Config
	configRepo := repositories.NewConfigInMemRepository()
	configService := services.NewConfigService(configRepo)
	configHandler := handlers.NewConfigHandler(configService)

	// Inicijalizacija repozitorijuma, servisa i handlera za ConfigGroup
	configGroupRepo := repositories.NewConfigGroupInMemRepository() // Pretpostavka da postoji ovaj repozitorijum
	configGroupService := services.NewConfigGroupService(configGroupRepo)
	configGroupHandler := handlers.NewConfigGroupHandler(configGroupService)

	// Kreiranje novog router-a
	router := mux.NewRouter()

	// Registracija ruta za ConfigHandler
	router.HandleFunc("/configs", configHandler.Add).Methods("POST")                       // Ruta za dodavanje konfiguracije
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")       // Ruta za pregled konfiguracije po imenu i verziji
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE") // Ruta za brisanje konfiguracije po imenu i verziji

	// Registracija ruta za ConfigGroupHandler
	router.HandleFunc("/config-groups", configGroupHandler.AddGroup).Methods("POST")                       // Ruta za dodavanje konfiguracione grupe
	router.HandleFunc("/config-groups/{name}/{version}", configGroupHandler.GetGroup).Methods("GET")       // Ruta za pregled konfiguracione grupe po imenu i verziji
	router.HandleFunc("/config-groups/{name}/{version}", configGroupHandler.DeleteGroup).Methods("DELETE") // Ruta za brisanje konfiguracione grupe po imenu i verziji
	router.HandleFunc("/config-groups/{name}/{version}/add-config", configGroupHandler.AddConfigToGroup).Methods("POST")
	router.HandleFunc("/config-groups/{name}/{version}/configs/{configName}/{configVersion}", configGroupHandler.RemoveConfigFromGroup).Methods("DELETE")

	// Pokretanje HTTP servera
	http.ListenAndServe("0.0.0.0:8000", router)
}
