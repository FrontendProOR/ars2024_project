package handlers

import (
	"encoding/json"
	"io"
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

	configName := mux.Vars(r)["configName"]
	configVersion, err := strconv.Atoi(mux.Vars(r)["configVersion"])
	if err != nil {
		http.Error(w, "Invalid config version format", http.StatusBadRequest)
		return
	}

	if err := h.repo.AddConfigToGroup(groupName, version, configName, configVersion); err != nil {
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

func decodeBodyLabels(body io.Reader) (map[string]string, error) {
	var labels map[string]string
	err := json.NewDecoder(body).Decode(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func renderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ConfigGroupHandler) GetConfigsByLabels(w http.ResponseWriter, r *http.Request) {
	// Extract name and version from request
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	// Convert version to int
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the group based on name and version
	group, err := h.repo.Get(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Decode labels from request body
	labels, err := decodeBodyLabels(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get configs by labels
	configs, err := h.repo.GetConfigsByLabels(group.Name, versionInt, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the filtered configs as JSON response
	renderJSON(w, configs)
}

func (h *ConfigGroupHandler) RemoveConfigsByLabels(w http.ResponseWriter, r *http.Request) {
	// Extract name and version from request
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	// Convert version to int
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the group based on name and version
	group, err := h.repo.Get(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Decode labels from request body
	labels, err := decodeBodyLabels(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete configs by labels
	err = h.repo.RemoveConfigsByLabels(group.Name, versionInt, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success message
	renderJSON(w, "deleted")
}
