package handlers

import (
	"encoding/json"
	"net/http"
	"project/model"
	"project/services"
	"strings"

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
	version := mux.Vars(r)["version"]

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

// Uklanjanje konfiguracione grupe
func (h *ConfigGroupHandler) RemoveGroup(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	if err := h.repo.Delete(name, version); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Config group successfully removed"))
}

// Dodavanje konfiguracije u grupu
func (h *ConfigGroupHandler) AddConfigToGroup(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	configName := mux.Vars(r)["configName"]
	configVersion := mux.Vars(r)["configVersion"]

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
	groupVersion := mux.Vars(r)["version"]

	configName := mux.Vars(r)["configName"]
	configVersion := mux.Vars(r)["configVersion"]

	if err := h.repo.RemoveConfigFromGroup(groupName, groupVersion, configName, configVersion); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Config successfully removed from group"))
}

// Dodavanje konfiguracije sa labelom u grupu
func (h *ConfigGroupHandler) AddConfigWithLabelToGroup(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	var config model.ConfigWithLabels
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.AddConfigWithLabelToGroup(groupName, version, config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Config with label successfully added to group"))
}

// Uklanjanje konfiguracija sa svim labelama iz grupe
func (h *ConfigGroupHandler) RemoveConfigsWithLabelsFromGroup(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["name"]
	groupVersion := mux.Vars(r)["version"]

	labelsParam := r.URL.Query().Get("labels")

	labelPairs := strings.Split(labelsParam, ";")

	labels := make([]model.Label, 0, len(labelPairs))
	for _, pair := range labelPairs {
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			http.Error(w, "Invalid label format. Expected format is key:value", http.StatusBadRequest)
			return
		}
		labels = append(labels, model.Label{Key: parts[0], Value: parts[1]})
	}

	if err := h.repo.RemoveConfigsWithLabelsFromGroup(groupName, groupVersion, labels); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configs with labels successfully removed from group"))
}

// Pretraga konfiguracije sa labelom unutar grupe
func (h *ConfigGroupHandler) SearchConfigsWithLabelsInGroup(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	labelsParam := r.URL.Query().Get("labels")

	labelPairs := strings.Split(labelsParam, ";")

	searchLabels := make([]model.Label, 0, len(labelPairs))
	for _, pair := range labelPairs {
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			http.Error(w, "Invalid label format. Expected format is key:value", http.StatusBadRequest)
			return
		}
		searchLabels = append(searchLabels, model.Label{Key: parts[0], Value: parts[1]})
	}

	configs, err := h.repo.SearchConfigsWithLabelsInGroup(groupName, version, searchLabels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(configs) == 0 {
		http.Error(w, "No configs with all labels found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(configs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
