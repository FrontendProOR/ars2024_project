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

// Add a new config
// @Summary Add a new config
// @Description Add a new config
// @Tags config
// @Accept json
// @Produce json
// @Param body body model.Config true "Config object to add"
// @Success 201 {string} string "Config successfully added"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
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

// Get a config by name and version
// @Summary Get a config by name and version
// @Description Get a config by name and version
// @Tags config
// @Accept json
// @Produce json
// @Param name path string true "Config name"
// @Param version path string true "Config version"
// @Success 200 {object} model.Config
// @Failure 404 {string} string "Config not found"
// @Failure 500 {string} string "Internal server error"
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

// Delete a config by name and version
// @Summary Delete a config by name and version
// @Description Delete a config by name and version
// @Tags config
// @Param name path string true "Config name"
// @Param version path string true "Config version"
// @Success 200 {string} string "Config successfully deleted"
// @Failure 404 {string} string "Config not found"
// @Failure 500 {string} string "Internal server error"
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
