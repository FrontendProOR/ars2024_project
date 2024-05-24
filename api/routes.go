// The `NewRouter` function creates a new router using Gorilla Mux and registers various routes for
// handling configuration and configuration group related HTTP requests.
package api

import (
	"net/http"
	"project/api/middleware"
	"project/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(configHandler *handlers.ConfigHandler, configGroupHandler *handlers.ConfigGroupHandler) *mux.Router {
	router := mux.NewRouter()

	// Registration of routes for ConfigHandler
	router.Handle("/configs", middleware.RateLimiter(http.HandlerFunc(configHandler.Add))).Methods("POST")
	router.Handle("/configs/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configHandler.Get))).Methods("GET")
	router.Handle("/configs/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configHandler.Delete))).Methods("DELETE")

	// Registration of routes for ConfigGroupHandler
	router.Handle("/config-groups", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.GetGroup))).Methods("GET")
	router.Handle("/config-groups/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveGroup))).Methods("DELETE")
	router.Handle("/config-groups/{name}/{version}/{configName}/{configVersion}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddConfigToGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}/config/search", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.SearchConfigsWithLabelsInGroup))).Methods("GET")
	router.Handle("/config-groups/{name}/{version}/config", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddConfigWithLabelToGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}/config/delete", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveConfigsWithLabelsFromGroup))).Methods("DELETE")
	router.Handle("/config-groups/{name}/{version}/{configName}/{configVersion}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveConfigFromGroup))).Methods("DELETE")

	// Registration of route for serving the frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/app.html")
	})

	return router
}
