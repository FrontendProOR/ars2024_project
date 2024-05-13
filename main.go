package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/handlers"
	"project/middleware"
	"project/repositories"
	"project/services"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
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
	router := mux.NewRouter()

	// Rate limiting 5 requests per second
	rateLimiter := rate.NewLimiter(4, 1)

	// Use the RateLimitMiddleware from the middleware package
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Registracija ruta za ConfigHandler
	router.HandleFunc("/configs", configHandler.Add).Methods("POST")                       // Ruta za dodavanje konfiguracije
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")       // Ruta za pregled konfiguracije po imenu i verziji
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE") // Ruta za brisanje konfiguracije po imenu i verziji

	// Registracija ruta za ConfigGroupHandler
	router.HandleFunc("/config-groups", configGroupHandler.AddGroup).Methods("POST")                                                              // Ruta za dodavanje konfiguracione grupe
	router.HandleFunc("/config-groups/{name}/{version}", configGroupHandler.GetGroup).Methods("GET")                                              // Ruta za pregled konfiguracione grupe po imenu i verziji
	router.HandleFunc("/config-groups/{name}/{version}", configGroupHandler.DeleteGroup).Methods("DELETE")                                        // Ruta za brisanje konfiguracione grupe po imenu i verziji
	router.HandleFunc("/config-groups/{name}/{version}/{configName}/{configVersion}", configGroupHandler.AddConfigToGroup).Methods("POST")        // Ruta za dodavanje postojeće konfiguracije u grupu
	router.HandleFunc("/config-groups/{name}/{version}/{configName}/{configVersion}", configGroupHandler.RemoveConfigFromGroup).Methods("DELETE") // Ruta za uklanjanje konfiguracije iz grupe

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	// Pokretanje HTTP servera u gorutini kako bi se moglo osluškivati za shutdown signal
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Osluškivanje za SIGINT i SIGTERM za Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
