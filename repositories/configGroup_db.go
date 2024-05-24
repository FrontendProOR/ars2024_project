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

	// For each config in the group, put it in its own "file" within the "configs" folder
	for _, config := range configGroup.Configs {
		_, err := repo.db.Put(fmt.Sprintf("config-groups/%s/%s/configs", configGroup.Name, configGroup.Version), config.Name, config.Version, config)
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

func (repo *ConfigGroupDBRepository) AddConfigToGroup(groupName string, version string, configName string, configVersion string) error {
	var config model.Config
	err := repo.db.Get(fmt.Sprintf("configs/%s/%s", configName, configVersion), &config)
	if err != nil {
		return err
	}
	configWithLabels := model.ConfigWithLabels{
		Config: config,
		Labels: []model.Label{},
	}
	return repo.AddConfigWithLabelToGroup(groupName, version, configWithLabels)
}

func (repo *ConfigGroupDBRepository) RemoveConfigFromGroup(groupName string, version string, configName string, configVersion string) error {
	return repo.db.Delete(fmt.Sprintf("config-groups/%s/%s/configs/%s/%s", groupName, version, configName, configVersion))
}

func (repo *ConfigGroupDBRepository) AddConfigWithLabelToGroup(groupName string, version string, config model.ConfigWithLabels) error {
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}
	configGroup.Configs = append(configGroup.Configs, &config)
	return repo.Add(configGroup)
}

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
	configGroup, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}
	var newConfigs []*model.ConfigWithLabels
	for _, config := range configGroup.Configs {
		if !containsLabels(config.Labels, labels) {
			newConfigs = append(newConfigs, config)
		}
	}
	configGroup.Configs = newConfigs
	return repo.Add(configGroup)
}

// Implement other methods (AddConfigToGroup, RemoveConfigFromGroup, etc.) similarly