// The `NewRouter` function creates a new router using Gorilla Mux and registers various routes for
// handling configuration and configuration group related HTTP requests.
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

	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("/docs/"))))

	// Serve Swagger UI
	router.PathPrefix("/swagger-ui/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	// Registration of routes for ConfigHandler
	router.Handle("/configs", middleware.RateLimiter(http.HandlerFunc(configHandler.Add))).Methods("POST")
	router.Handle("/configs/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configHandler.Get))).Methods("GET")
	router.Handle("/configs/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configHandler.Delete))).Methods("DELETE")

	// Registration of routes for ConfigGroupHandler
	router.Handle("/config-groups", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.GetGroup))).Methods("GET")
	router.Handle("/config-groups/{name}/{version}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveGroup))).Methods("DELETE")
	router.Handle("/config-groups/{name}/{version}/{configName}/{configVersion}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddConfigToGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}/configs/search", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.SearchConfigsWithLabelsInGroup))).Methods("GET")
	router.Handle("/config-groups/{name}/{version}/configs", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.AddConfigWithLabelToGroup))).Methods("POST")
	router.Handle("/config-groups/{name}/{version}/configs/delete", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveConfigsWithLabelsFromGroup))).Methods("DELETE")
	router.Handle("/config-groups/{name}/{version}/{configName}/{configVersion}", middleware.RateLimiter(http.HandlerFunc(configGroupHandler.RemoveConfigFromGroup))).Methods("DELETE")

	// Registration of route for serving the frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/app.html")
	})

	return router
}
