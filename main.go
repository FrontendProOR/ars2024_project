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

	_ "project/docs"

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

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Running the server
	api.RunServer(router)
}
