package handlers

import (
	"encoding/json"
	"net/http"
	"project/model"
	"project/services"
	"strconv"

	"github.com/gorilla/mux"
)

type ConfigGroupHandler struct {
	repo services.ConfigGroupService
}

func NewConfigGroupHandler(repo services.ConfigGroupService) *ConfigGroupHandler {
	return &ConfigGroupHandler{
		repo: repo,
	}
}

// Dodavanje konfiguracione grupe
func (h *ConfigGroupHandler) AddGroup(w http.ResponseWriter, r *http.Request) {
	var group model.ConfigGroup
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Add(group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Config group successfully added"))
}

// Pregled konfiguracione grupe
func (h *ConfigGroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version, err := strconv.Atoi(mux.Vars(r)["version"])
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}

	group, err := h.repo.Get(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// Brisanje konfiguracione grupe
func (h *ConfigGroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version, err := strconv.Atoi(mux.Vars(r)["version"])
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(name, version); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Config group successfully deleted"))
}

// Dodavanje konfiguracije u grupu
func (h *ConfigGroupHandler) AddConfigToGroup(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["name"]
	version, err := strconv.Atoi(mux.Vars(r)["version"])
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}

	var config model.ConfigWithLabels
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.AddConfigToGroup(groupName, version, config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Config successfully added to group"))
}

// Uklanjanje konfiguracije iz grupe
func (h *ConfigGroupHandler) RemoveConfigFromGroup(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["name"]
	groupVersion, err := strconv.Atoi(mux.Vars(r)["version"])
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}
	configName := mux.Vars(r)["configName"]
	configVersion, err := strconv.Atoi(mux.Vars(r)["configVersion"]) // Pretpostavka da je configVersion dostupan
	if err != nil {
		http.Error(w, "Invalid config version format", http.StatusBadRequest)
		return
	}

	if err := h.repo.RemoveConfigFromGroup(groupName, groupVersion, configName, configVersion); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Config successfully removed from group"))
}
