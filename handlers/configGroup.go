// The above code defines a set of HTTP handlers for managing configuration groups and configurations
// with labels in a Go application.
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

// AddGroup adds a new configuration group
// @Summary Add a new configuration group
// @Description Add a new configuration group
// @Tags config-groups
// @Accept  json
// @Produce  json
// @Param group body model.ConfigGroup true "Config Group"
// @Success 201 {string} string "Config group successfully added"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups [post]
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

// GetGroup retrieves a configuration group
// @Summary Get a configuration group
// @Description Get a configuration group
// @Tags config-groups
// @Accept  json
// @Produce  json
// @Param name path string true "Config Group Name"
// @Param version path string true "Config Group Version"
// @Success 200 {object} model.ConfigGroup
// @Failure 404 {string} string "Config group not found"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups/{name}/{version} [get]
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

// RemoveGroup removes a configuration group
// @Summary Remove a configuration group
// @Description Remove a configuration group
// @Tags config-groups
// @Accept  json
// @Produce  json
// @Param name path string true "Config Group Name"
// @Param version path string true "Config Group Version"
// @Success 200 {string} string "Config group successfully removed"
// @Failure 404 {string} string "Config group not found"
// @Router /config-groups/{name}/{version} [delete]
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

// AddConfigToGroup adds a configuration to a group
// @Summary Add a configuration to a group
// @Description Add a configuration to a group
// @Tags config-groups
// @Accept  json
// @Produce  json
// @Param name path string true "Config Group Name"
// @Param version path string true "Config Group Version"
// @Param configName path string true "Config Name"
// @Param configVersion path string true "Config Version"
// @Success 201 {string} string "Config successfully added to group"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups/{name}/{version}/{configName}/{configVersion} [post]
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

// RemoveConfigFromGroup removes a configuration from a group
// @Summary Remove a configuration from a group
// @Description Remove a configuration from a group
// @Tags config-groups
// @Accept  json
// @Produce  json
// @Param name path string true "Config Group Name"
// @Param version path string true "Config Group Version"
// @Param configName path string true "Config Name"
// @Param configVersion path string true "Config Version"
// @Success 200 {string} string "Config successfully removed from group"
// @Failure 404 {string} string "Config group not found"
// @Router /config-groups/{name}/{version}/configs/{configName}/{configVersion} [delete]
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

// AddConfigWithLabelToGroup adds a configuration with labels to a group
// @Summary Add a configuration with labels to a group
// @Description Add a configuration with labels to a group
// @Tags config-groups
// @Accept  json
// @Produce  json
// @Param name path string true "Config Group Name"
// @Param version path string true "Config Group Version"
// @Param config body model.ConfigWithLabels true "Config with Labels"
// @Success 201 {string} string "Config with label successfully added to group"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups/{name}/{version}/configs [post]
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

// RemoveConfigsWithLabelsFromGroup removes configurations with labels from a group
// @Summary Remove configurations with labels from a group
// @Description Remove configurations with labels from a group
// @Tags config-groups
// @Accept  json
// @Produce  json
// @Param name path string true "Config Group Name"
// @Param version path string true "Config Group Version"
// @Param labels path string true "Labels (key:value;key:value)"
// @Param configName path string true "Config Name"
// @Param configVersion path string true "Config Version"
// @Success 200 {string} string "Configs with labels successfully removed from group"
// @Failure 400 {string} string "Invalid label format"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups/{name}/{version}/configs/{labels}/{configName}/{configVersion} [delete]
func (h *ConfigGroupHandler) RemoveConfigsWithLabelsFromGroup(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	vars := mux.Vars(r)
	labelsParam := vars["labels"]

	labelPairs := strings.Split(labelsParam, ";")

	configName := mux.Vars(r)["configName"]
	configVersion := mux.Vars(r)["configVersion"]

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

	if err := h.repo.RemoveConfigsWithLabelsFromGroup(groupName, version, labels, configName, configVersion); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configs with labels successfully removed from group"))
}

// SearchConfigsWithLabelsInGroup searches for configurations with labels in a group
// @Summary Search for configurations with labels in a group
// @Description Search for configurations with labels in a group
// @Tags config-groups
// @Accept  json
// @Produce  json
// @Param name path string true "Config Group Name"
// @Param version path string true "Config Group Version"
// @Param labels path string true "Labels (key:value;key:value)"
// @Param configName path string true "Config Name"
// @Param configVersion path string true "Config Version"
// @Success 200 {array} model.ConfigWithLabels
// @Failure 400 {string} string "Invalid label format"
// @Failure 404 {string} string "No configs with all labels found"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups/{name}/{version}/configs/{labels}/{configName}/{configVersion} [get]
func (h *ConfigGroupHandler) SearchConfigsWithLabelsInGroup(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	vars := mux.Vars(r)
	labelsParam := vars["labels"]

	labelPairs := strings.Split(labelsParam, ";")

	configName := mux.Vars(r)["configName"]
	configVersion := mux.Vars(r)["configVersion"]

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

	configs, err := h.repo.SearchConfigsWithLabelsInGroup(groupName, version, searchLabels, configName, configVersion)
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
