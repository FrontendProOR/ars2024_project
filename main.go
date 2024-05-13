package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/handlers"
	"project/handlers/middleware"
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

	// 10 requests inistally
	// refills at rate 10 req/min (1 req per 6 seconds)
	limiter := rate.NewLimiter(0.167, 10)

	// Use the RateLimitMiddleware from the middleware package

	// Registracija ruta za ConfigHandler
	router.Handle("/configs", middleware.RateLimit(limiter, configHandler.Add)).Methods("POST")                       // Ruta za dodavanje konfiguracije
	router.Handle("/configs/{name}/{version}", middleware.RateLimit(limiter, configHandler.Get)).Methods("GET")       // Ruta za pregled konfiguracije po imenu i verziji
	router.Handle("/configs/{name}/{version}", middleware.RateLimit(limiter, configHandler.Delete)).Methods("DELETE") // Ruta za brisanje konfiguracije po imenu i verziji

	// Registracija ruta za ConfigGroupHandler
	router.Handle("/config-groups/{name}/{version}", middleware.RateLimit(limiter, configGroupHandler.GetGroup)).Methods("GET")                                              // Ruta za pregled konfiguracione grupe po imenu i verziji
	router.Handle("/config-groups", middleware.RateLimit(limiter, configGroupHandler.AddGroup)).Methods("POST")                                                              // Ruta za dodavanje konfiguracione grupe
	router.Handle("/config-groups/{name}/{version}", middleware.RateLimit(limiter, configGroupHandler.DeleteGroup)).Methods("DELETE")                                        // Ruta za brisanje konfiguracione grupe po imenu i verziji
	router.Handle("/config-groups/{name}/{version}/{configName}/{configVersion}", middleware.RateLimit(limiter, configGroupHandler.AddConfigToGroup)).Methods("POST")        // Ruta za dodavanje postojeće konfiguracije u grupu
	router.Handle("/config-groups/{name}/{version}/{configName}/{configVersion}", middleware.RateLimit(limiter, configGroupHandler.RemoveConfigFromGroup)).Methods("DELETE") // Ruta za uklanjanje konfiguracije iz grupe

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
