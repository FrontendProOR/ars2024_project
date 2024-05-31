// The above code defines a repository for managing configuration groups and their associated
// configurations in a Go project.
package repositories

import (
	"errors"
	"fmt"
	"project/data"
	"project/model"
	"sort"
	"strings"
)

type ConfigGroupDBRepository struct {
	db *data.Database
}

func NewConfigGroupDBRepository(db *data.Database) *ConfigGroupDBRepository {
	return &ConfigGroupDBRepository{
		db: db,
	}
}

// This `Add` method in the `ConfigGroupDBRepository` struct is responsible for adding a new
// configuration group to the repository. Here's a breakdown of what the method does:
func (repo *ConfigGroupDBRepository) Add(configGroup model.ConfigGroup) error {
	// Validation
	if strings.TrimSpace(configGroup.Name) == "" {
		return errors.New("configGroup name cannot be empty")
	}
	if strings.TrimSpace(configGroup.Version) == "" {
		return errors.New("configGroup version cannot be empty")
	}

	// Check if the group already exists
	existingKeys, err := repo.db.List("config-groups/")
	if err != nil {
		return err
	}
	for key := range existingKeys {
		// Split the key into parts
		parts := strings.Split(key, "/")
		// Check if the key matches the name and version of the group being added
		if len(parts) >= 3 && parts[1] == configGroup.Name && parts[2] == configGroup.Version {
			return errors.New("configGroup with this name and version already exists")
		}
	}

	// Check if the group has configs
	hasConfigs := len(configGroup.Configs) > 0

	// Add the group key without value only if it has no configs
	if !hasConfigs {
		_, err = repo.db.Put("config-groups", configGroup.Name, configGroup.Version, nil)
		if err != nil {
			return err
		}
	}

	// Add configs to the group
	for _, config := range configGroup.Configs {
		keyType := "config-groups"
		name := fmt.Sprintf("%s/%s/configs/%s", configGroup.Name, configGroup.Version, config.Name)
		_, err = repo.db.Put(keyType, name, config.Version, config)
		if err != nil {
			return err
		}
	}

	return nil
}

// The `Get` method in the `ConfigGroupDBRepository` struct is responsible for retrieving a specific
// configuration group by its name and version from the repository. Here's a breakdown of what the
// method does:
func (repo *ConfigGroupDBRepository) Get(name string, version string) (model.ConfigGroup, error) {
	var configGroup model.ConfigGroup
	configGroup.Name = name
	configGroup.Version = version

	// Check if the config group exists
	existingKeys, err := repo.db.List("config-groups/")
	if err != nil {
		return model.ConfigGroup{}, err
	}
	found := false
	for key := range existingKeys {
		// Split the key into parts
		parts := strings.Split(key, "/")
		// Check if the key matches the name and version of the group being added
		if len(parts) >= 3 && parts[1] == name && parts[2] == version {
			found = true
			break
		}
	}
	if !found {
		return model.ConfigGroup{}, errors.New("configGroup not found")
	}

	configs, err := repo.db.List(fmt.Sprintf("config-groups/%s/%s/configs", name, version))
	if err != nil {
		return model.ConfigGroup{}, err
	}
	for key := range configs {
		var config model.ConfigWithLabels
		err := repo.db.Get(key, &config)
		if err != nil {
			return model.ConfigGroup{}, err
		}
		configGroup.Configs = append(configGroup.Configs, &config)
	}
	return configGroup, nil
}

// This `Delete` method in the `ConfigGroupDBRepository` struct is responsible for deleting a specific
// configuration group by its name and version from the repository. Here's a breakdown of what the
// method does:
func (repo *ConfigGroupDBRepository) Delete(name string, version string) error {
	// Check if the group exists
	_, err := repo.Get(name, version)
	if err != nil {
		// If the group does not exist, return the error
		return err
	}

	configs, err := repo.db.List(fmt.Sprintf("config-groups/%s/%s/configs", name, version))
	if err != nil {
		return err
	}
	for key := range configs {
		err := repo.db.Delete(key)
		if err != nil {
			return err
		}
	}
	return repo.db.Delete(fmt.Sprintf("config-groups/%s/%s", name, version))
}

// The `AddConfigToGroup` method in the `ConfigGroupDBRepository` struct is responsible for adding a
// new configuration to a specific configuration group within the repository. Here's a breakdown of
// what the method does:
func (repo *ConfigGroupDBRepository) AddConfigToGroup(groupName string, version string, configName string, configVersion string) error {
	// Get the config
	var config model.Config
	err := repo.db.Get(fmt.Sprintf("configs/%s/%s", configName, configVersion), &config)
	if err != nil {
		return err
	}

	// Get the config group
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}

	// Check if the config already exists in the group
	for _, existingConfig := range configGroup.Configs {
		if existingConfig.Name == configName && existingConfig.Version == configVersion {
			return errors.New("config already exists in the group")
		}
	}

	// Check if the config group is empty
	isEmpty := len(configGroup.Configs) == 0

	// Add the config to the group
	keyType := "config-groups"
	name := fmt.Sprintf("%s/%s/configs/%s", groupName, version, configName)
	_, err = repo.db.Put(keyType, name, configVersion, config)
	if err != nil {
		return err
	}

	// If the config group was empty, delete the old key
	if isEmpty {
		oldKey := fmt.Sprintf("config-groups/%s/%s", groupName, version)
		err = repo.db.Delete(oldKey)
		if err != nil {
			return err
		}
	}

	return nil
}

// The `Update` method in the `ConfigGroupDBRepository` struct is responsible for updating an existing
// configuration group in the repository. Here's a breakdown of what the method does:
func (repo *ConfigGroupDBRepository) Update(configGroup model.ConfigGroup) error {
	// Validation
	if strings.TrimSpace(configGroup.Name) == "" {
		return errors.New("configGroup name cannot be empty")
	}
	if strings.TrimSpace(configGroup.Version) == "" {
		return errors.New("configGroup version cannot be empty")
	}

	// Update the group
	_, err := repo.db.Put("config-groups", configGroup.Name, configGroup.Version, configGroup)
	if err != nil {
		return err
	}

	return nil
}

// The `RemoveConfigFromGroup` method in the `ConfigGroupDBRepository` struct is responsible for
// removing a specific configuration from a configuration group within the repository. Here's a
// breakdown of what the method does:
func (repo *ConfigGroupDBRepository) RemoveConfigFromGroup(groupName string, version string, configName string, configVersion string) error {
	// Get the config group
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}

	// Check if the config exists in the group
	found := false
	for i, existingConfig := range configGroup.Configs {
		if existingConfig.Name == configName && existingConfig.Version == configVersion {
			// Remove the config from the group
			configGroup.Configs = append(configGroup.Configs[:i], configGroup.Configs[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("config not found in the group")
	}

	// Delete the config from the database
	key := fmt.Sprintf("config-groups/%s/%s/configs/%s/%s", groupName, version, configName, configVersion)
	err = repo.db.Delete(key)
	if err != nil {
		return err
	}

	// Check if there are no more configs in the group
	if len(configGroup.Configs) == 0 {
		// Update the group key in the database
		keyType := "config-groups"
		_, err = repo.db.Put(keyType, groupName, version, configGroup)
		if err != nil {
			return err
		}
	}

	return nil
}

// This `AddConfigWithLabelToGroup` method in the `ConfigGroupDBRepository` struct is responsible for
// adding a new configuration with labels to a specific configuration group within the repository.
// Here's a breakdown of what the method does:
func (repo *ConfigGroupDBRepository) AddConfigWithLabelToGroup(groupName string, version string, config model.ConfigWithLabels) error {
	// Get the config group
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}

	// Check if the config already exists in the group
	for _, existingConfig := range configGroup.Configs {
		if existingConfig.Name == config.Name && existingConfig.Version == config.Version {
			return errors.New("config already exists in the group")
		}
	}

	// Sort the labels to ensure consistency
	sort.Slice(config.Labels, func(i, j int) bool {
		return config.Labels[i].Key < config.Labels[j].Key
	})

	// Check if the config group is empty
	isEmpty := len(configGroup.Configs) == 0

	// Add the config to the group, ensure that labels are part of key in Consul /config-groups/{name}/{version}/{labels}/{configName}/{configVersion}
	keyType := "config-groups"
	labels := ""
	// labels are in format key1:value1;key2:value2
	for _, label := range config.Labels {
		labels += fmt.Sprintf("%s:%s;", label.Key, label.Value)
	}
	name := fmt.Sprintf("%s/%s/configs/%s/%s", groupName, version, labels, config.Name)
	// If the config group is empty, delete the old key
	if isEmpty {
		oldKey := fmt.Sprintf("config-groups/%s/%s", groupName, version)
		err = repo.db.Delete(oldKey)
		if err != nil {
			return err
		}
	}
	_, err = repo.db.Put(keyType, name, config.Version, config)
	if err != nil {
		return err
	}

	return nil

}

// This `SearchConfigsWithLabelsInGroup` method in the `ConfigGroupDBRepository` struct is responsible
// for searching and retrieving configurations within a specific configuration group that match a given
// set of labels. Here's a breakdown of what the method does:
func (repo *ConfigGroupDBRepository) SearchConfigsWithLabelsInGroup(groupName string, version string, labels []model.Label, configName string, configVersion string) ([]*model.ConfigWithLabels, error) {
	// Check if Config Group exists
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return nil, err
	}

	// Convert labels to a map
	labelsMap := make(map[string]string)
	for _, label := range labels {
		if label.Key == "" || label.Value == "" {
			return nil, errors.New("invalid label format. Expected format is key:value")
		}
		labelsMap[label.Key] = label.Value
	}

	// Then search if there are any configs with the given labels
	var matchingConfigs []*model.ConfigWithLabels
	for _, config := range configGroup.Configs {
		if containsAllLabels(config.Labels, labelsMap) {
			if config.Config.Name != configName {
				return nil, errors.New("config name does not match")
			}
			if config.Config.Version != configVersion {
				return nil, errors.New("config version does not match")
			}
			matchingConfigs = append(matchingConfigs, config)
		}
	}

	// Then compare configName and configVersion with found configs
	if len(matchingConfigs) == 0 {
		return nil, errors.New("no configs with all labels found")
	}

	return matchingConfigs, nil
}

// The `RemoveConfigsWithLabelsFromGroup` method in the `ConfigGroupDBRepository` struct is responsible
// for removing configurations from a specific configuration group that match a given set of labels.
// Here's a breakdown of what the method does:
func (repo *ConfigGroupDBRepository) RemoveConfigsWithLabelsFromGroup(groupName string, version string, labels []model.Label, configName string, configVersion string) error {
	// Check if the config name, version and labels are valid
	if configName == "" {
		return errors.New("config name cannot be empty")
	}
	if configVersion == "" {
		return errors.New("config version cannot be empty")
	}
	if len(labels) == 0 {
		return errors.New("at least one label must be provided")
	}

	// Get the config group
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}

	// Convert labels to a map
	labelsMap := make(map[string]string)
	for _, label := range labels {
		if label.Key == "" || label.Value == "" {
			return errors.New("invalid label format. Expected format is key:value")
		}
		labelsMap[label.Key] = label.Value
	}

	// Convert labels to a string
	labelsStr := ""
	for _, label := range labels {
		labelsStr += fmt.Sprintf("%s:%s;", label.Key, label.Value)
	}

	// If configName or configVersion are incorrect, return an error
	if configName == "" || configVersion == "" {
		return errors.New("config name and version must be provided")
	}

	// If configName or configVersion is not found in the group, return an error
	found := false
	for _, config := range configGroup.Configs {
		if config.Config.Name == configName && config.Config.Version == configVersion {
			found = true
			break
		}
	}
	if !found {
		return errors.New("config not found in the group")
	}

	// Find configs with the given labels and matching config name and version
	var configsToRemove []*model.ConfigWithLabels
	for _, config := range configGroup.Configs {
		if containsAllLabels(config.Labels, labelsMap) && config.Config.Name == configName && config.Config.Version == configVersion {
			configsToRemove = append(configsToRemove, config)
		}
	}

	if len(configsToRemove) == 0 {
		return errors.New("no configs with all labels found to remove")
	}

	// Remove the matching configs from the group
	for _, configToRemove := range configsToRemove {
		configKey := fmt.Sprintf("/config-groups/%s/%s/configs/%s/%s/%s", groupName, version, labelsStr, configToRemove.Config.Name, configToRemove.Config.Version)
		err = repo.db.Delete(configKey)
		if err != nil {
			return err
		}
	}

	// If all configs are removed, update the group with an empty configs array
	if len(configGroup.Configs) == len(configsToRemove) {
		configGroup.Configs = []*model.ConfigWithLabels{}
		err = repo.Update(configGroup)
		if err != nil {
			return err
		}
	}

	return nil
}

func containsAllLabels(configLabels []model.Label, labelsMap map[string]string) bool {
	for _, label := range configLabels {
		if labelsMap[label.Key] != label.Value {
			return false
		}
	}
	return true
}
