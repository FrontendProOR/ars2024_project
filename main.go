package main

import (
	"project/api"
	"project/handlers"
	"project/repositories"
	"project/services"
)

func main() {
	// Inicijalizacija repozitorijuma, servisa i handlera za Config
	configRepo := repositories.NewConfigInMemRepository()
	configService := services.NewConfigService(configRepo)
	configHandler := handlers.NewConfigHandler(configService)
	// Inicijalizacija repozitorijuma, servisa i handlera za ConfigGroup
	configGroupRepo := repositories.NewConfigGroupInMemRepository(configRepo) // Pretpostavka da postoji ovaj repozitorijum
	configGroupService := services.NewConfigGroupService(configGroupRepo)
	configGroupHandler := handlers.NewConfigGroupHandler(configGroupService)
	// Kreiranje novog router-a
	router := api.NewRouter(configHandler, configGroupHandler)

	// Pokretanje servera
	api.RunServer(router)
}
