package handlers

import (
	"encoding/json"
	"net/http"
	"project/model"
	"project/services"

	"github.com/gorilla/mux"
)

type ConfigHandler struct {
	service services.ConfigService
}

func NewConfigHandler(service services.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		service: service,
	}
}

// Add a new configuration
// @Summary Add configuration
// @Description Add a new configuration
// @Tags config
// @Accept  json
// @Produce  json
// @Param config body model.Config true "Configuration object"
// @Success 201 {string} string "Config successfully added"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /configs [post]
func (c ConfigHandler) Add(w http.ResponseWriter, r *http.Request) {
	var config model.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.service.Add(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Config successfully added"))
}

// Get configuration by name and version
// @Summary Get configuration
// @Description Get a configuration by name and version
// @Tags config
// @Accept  json
// @Produce  json
// @Param name path string true "Configuration name"
// @Param version path string true "Configuration version"
// @Success 200 {object} model.Config
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /configs/{name}/{version} [get]
func (c ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	config, err := c.service.Get(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// Delete configuration by name and version
// @Summary Delete configuration
// @Description Delete a configuration by name and version
// @Tags config
// @Accept  json
// @Produce  json
// @Param name path string true "Configuration name"
// @Param version path string true "Configuration version"
// @Success 200 {string} string "Config successfully deleted"
// @Failure 404 {string} string "Not Found"
// @Router /configs/{name}/{version} [delete]
func (c ConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	if err := c.service.Delete(name, version); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Config successfully deleted"))
}
