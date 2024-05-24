// The main function initializes database, repositories, services, handlers, and a router for a Go
// project and starts the server.
package main

import (
	"log"
	"project/api"
	"project/data"
	"project/handlers"
	"project/repositories"
	"project/services"
)

func main() {
	// Initialisation of database
	db, err := data.NewDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Initialisation of repositories, services, and handlers for Config
	configRepo := repositories.NewConfigDBRepository(db)
	configService := services.NewConfigService(configRepo)
	configHandler := handlers.NewConfigHandler(configService)
	// Initialisation of repositories, services, and handlers for ConfigGroup
	configGroupRepo := repositories.NewConfigGroupDBRepository(db)
	configGroupService := services.NewConfigGroupService(configGroupRepo)
	configGroupHandler := handlers.NewConfigGroupHandler(configGroupService)
	// Creating a new router
	router := api.NewRouter(configHandler, configGroupHandler)

	// Running the server
	api.RunServer(router)
}
