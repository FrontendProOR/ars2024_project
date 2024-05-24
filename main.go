package main

import (
	"project/api"
	"project/handlers"
	"project/repositories"
	"project/services"

	_ "project/docs" // uključite generisanu dokumentaciju

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Config API
// @version 1.0
// @description API for managing configurations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

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

	// Dodavanje Swagger endpointa
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Pokretanje servera
	api.RunServer(router)
}
