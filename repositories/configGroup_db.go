// The above code defines a repository for managing configuration groups and their associated
// configurations in a Go project.
package repositories

import (
	"errors"
	"fmt"
	"project/data"
	"project/model"
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
	// Check if the config exists in the group
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}
	found := false
	for _, config := range configGroup.Configs {
		if config.Name == configName && config.Version == configVersion {
			found = true
			break
		}
	}
	if !found {
		return errors.New("config not found in the group")
	}

	// If the config exists, delete it
	return repo.db.Delete(fmt.Sprintf("config-groups/%s/%s/configs/%s/%s", groupName, version, configName, configVersion))
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

	// Add the config to the group
	keyType := "config-groups"
	name := fmt.Sprintf("%s/%s/configs/%s", groupName, version, config.Name)
	_, err = repo.db.Put(keyType, name, config.Version, config)
	if err != nil {
		return err
	}

	return nil
}

// This `SearchConfigsWithLabelsInGroup` method in the `ConfigGroupDBRepository` struct is responsible
// for searching and retrieving configurations within a specific configuration group that match a given
// set of labels. Here's a breakdown of what the method does:
func (repo *ConfigGroupDBRepository) SearchConfigsWithLabelsInGroup(groupName string, version string, labels []model.Label) ([]*model.ConfigWithLabels, error) {
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return nil, err
	}
	var result []*model.ConfigWithLabels
	for _, config := range configGroup.Configs {
		if containsLabels(config.Labels, labels) {
			result = append(result, config)
		}
	}
	return result, nil
}

func containsLabels(configLabels []model.Label, labels []model.Label) bool {
	for _, label := range labels {
		found := false
		for _, configLabel := range configLabels {
			if configLabel.Key == label.Key && configLabel.Value == label.Value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (repo *ConfigGroupDBRepository) RemoveConfigsWithLabelsFromGroup(groupName string, version string, labels []model.Label) error {
	// Get the config group
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}

	// Search for configs with the given labels and number of labels
	var configsToRemove []*model.ConfigWithLabels
	for _, config := range configGroup.Configs {
		if containsAllLabels(config.Labels, labels) {
			configsToRemove = append(configsToRemove, config)
		}
	}

	// If there is config with exactly the same labels (in terms of number and values) as the given labels
	if len(configsToRemove) == 0 {
		return errors.New("no configs with all labels found")
	}

	// Remove the configs
	for _, config := range configsToRemove {
		err := repo.RemoveConfigFromGroup(groupName, version, config.Name, config.Version)
		if err != nil {
			return err
		}
	}

	return nil
}

func containsAllLabels(configLabels []model.Label, labels []model.Label) bool {
	// Check if all labels are present in the config in terms of number and values
	if len(configLabels) != len(labels) {
		return false
	}

	for _, label := range labels {
		found := false
		for _, configLabel := range configLabels {
			if configLabel.Key == label.Key && configLabel.Value == label.Value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
