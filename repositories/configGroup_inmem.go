package repositories

import (
	"errors"
	"project/model"
)

type ConfigGroupInMemRepository struct {
	configGroups map[string]model.ConfigGroup
}

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return &ConfigGroupInMemRepository{
		configGroups: make(map[string]model.ConfigGroup),
	}
}

func (repo *ConfigGroupInMemRepository) Add(configGroup model.ConfigGroup) error {
	repo.configGroups[configGroup.Name] = configGroup
	return nil
}

func (repo *ConfigGroupInMemRepository) Get(name string, version int) (model.ConfigGroup, error) {
	configGroup, exists := repo.configGroups[name]
	if !exists {
		return model.ConfigGroup{}, errors.New("config group not found")
	}
	return configGroup, nil
}

func (repo *ConfigGroupInMemRepository) Delete(name string, version int) error {
	if _, exists := repo.configGroups[name]; !exists {
		return errors.New("config group not found")
	}
	delete(repo.configGroups, name)
	return nil
}

// Ove metode (AddConfigToGroup i RemoveConfigFromGroup) su primeri i mogu zahtevati dodatnu implementaciju
func (repo *ConfigGroupInMemRepository) AddConfigToGroup(groupName string, version int, config model.ConfigWithLabels) error {
	// Implementacija dodavanja konfiguracije u grupu
	return nil
}

func (repo *ConfigGroupInMemRepository) RemoveConfigFromGroup(groupName string, version int, configName string) error {
	// Implementacija uklanjanja konfiguracije iz grupe
	return nil
}
