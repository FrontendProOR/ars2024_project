package api

import (
	"net/http"
	"project/api/middleware"
	"project/handlers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(configHandler *handlers.ConfigHandler, configGroupHandler *handlers.ConfigGroupHandler) *mux.Router {
	router := mux.NewRouter()

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Registracija ruta za ConfigHandler
	router.Handle("/configs", middleware.RateLimiter(http.HandlerFunc(configHandler.Add))).Methods("POST")
	router.Handle("/configs/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configHandler.Get))).Methods("GET")
	router.Handle("/configs/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configHandler.Delete))).Methods("DELETE")

	// Registracija ruta za ConfigGroupHandler
	router.Handle("/config-groups", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.GetGroup))).Methods("GET")
	router.Handle("/config-groups/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveGroup))).Methods("DELETE")
	router.Handle("/config-groups/{name}/{version}/{configName}/{configVersion}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddConfigToGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}/config/search", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.SearchConfigsWithLabelsInGroup))).Methods("GET")
	router.Handle("/config-groups/{name}/{version}/config", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddConfigWithLabelToGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}/config/delete", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveConfigsWithLabelsFromGroup))).Methods("DELETE")
	router.Handle("/config-groups/{name}/{version}/{configName}/{configVersion}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveConfigFromGroup))).Methods("DELETE")

	return router
}
