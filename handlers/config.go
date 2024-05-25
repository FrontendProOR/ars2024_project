// The code defines a ConfigHandler struct with methods for adding, retrieving, and deleting
// configurations using a ConfigService.
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

// Adds a new configuration
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

// Retrieves a configuration
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

// Deletes a configuration
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
