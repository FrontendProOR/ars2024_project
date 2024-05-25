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

// Add a configuration group
// @Summary Add a new configuration group
// @Description Add a new configuration group
// @Tags config-group
// @Accept json
// @Produce json
// @Param body body model.ConfigGroup true "Configuration group object to add"
// @Success 201 {string} string "Configuration group successfully added"
// @Failure 400 {string} string "Bad request"
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

// Get a configuration group by name and version
// @Summary Get a configuration group by name and version
// @Description Get a configuration group by name and version
// @Tags config-group
// @Accept json
// @Produce json
// @Param name path string true "Configuration group name"
// @Param version path string true "Configuration group version"
// @Success 200 {object} model.ConfigGroup
// @Failure 404 {string} string "Configuration group not found"
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

// Remove a configuration group by name and version
// @Summary Remove a configuration group by name and version
// @Description Remove a configuration group by name and version
// @Tags config-group
// @Param name path string true "Configuration group name"
// @Param version path string true "Configuration group version"
// @Success 200 {string} string "Configuration group successfully removed"
// @Failure 404 {string} string "Configuration group not found"
// @Failure 500 {string} string "Internal server error"
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

// Add a configuration to a group
// @Summary Add a configuration to a group
// @Description Add a configuration to a group
// @Tags config-group
// @Accept json
// @Produce json
// @Param name path string true "Configuration group name"
// @Param version path string true "Configuration group version"
// @Param configName path string true "Configuration name"
// @Param configVersion path string true "Configuration version"
// @Success 201 {string} string "Configuration successfully added to group"
// @Failure 400 {string} string "Bad request"
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

// Remove a configuration from a group
// @Summary Remove a configuration from a group
// @Description Remove a configuration from a group
// @Tags config-group
// @Param name path string true "Configuration group name"
// @Param version path string true "Configuration group version"
// @Param configName path string true "Configuration name"
// @Param configVersion path string true "Configuration version"
// @Success 200 {string} string "Configuration successfully removed from group"
// @Failure 404 {string} string "Configuration not found in group"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups/{name}/{version}/{configName}/{configVersion} [delete]
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

// Add a configuration with labels to a group
// @Summary Add a configuration with labels to a group
// @Description Add a configuration with labels to a group
// @Tags config-group
// @Accept json
// @Produce json
// @Param name path string true "Configuration group name"
// @Param version path string true "Configuration group version"
// @Param body body model.ConfigWithLabels true "Configuration with labels object to add to group"
// @Success 201 {string} string "Configuration with label successfully added to group"
// @Failure 400 {string} string "Bad request"
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

// Remove configurations with labels from a group
// @Summary Remove configurations with labels from a group
// @Description Remove configurations with labels from a group
// @Tags config-group
// @Param name path string true "Configuration group name"
// @Param version path string true "Configuration group version"
// @Param labels query []string true "Labels to match in key:value format (e.g., label1:value1;label2:value2)"
// @Success 200 {string} string "Configurations with labels successfully removed from group"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups/{name}/{version}/configs/delete [delete]
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

// Search configurations with labels in a group
// @Summary Search configurations with labels in a group
// @Description Search configurations with labels in a group
// @Tags config-group
// @Accept json
// @Produce json
// @Param name path string true "Configuration group name"
// @Param version path string true "Configuration group version"
// @Param labels query []string true "Labels to match in key:value format (e.g., label1:value1;label2:value2)"
// @Success 200 {object} []model.Config
// @Failure 404 {string} string "No configurations found with specified labels"
// @Failure 500 {string} string "Internal server error"
// @Router /config-groups/{name}/{version}/configs/search [get]
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
