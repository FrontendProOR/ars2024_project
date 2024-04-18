package repositories

import (
	"errors"
	"project/model"
)

type ConfigGroupInMemRepository struct {
	configGroups map[string]map[int]model.ConfigGroup
}

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return &ConfigGroupInMemRepository{
		configGroups: make(map[string]map[int]model.ConfigGroup),
	}
}

func (repo *ConfigGroupInMemRepository) Add(configGroup model.ConfigGroup) error {
	if _, exists := repo.configGroups[configGroup.Name]; !exists {
		repo.configGroups[configGroup.Name] = make(map[int]model.ConfigGroup)
	}
	repo.configGroups[configGroup.Name][configGroup.Version] = configGroup
	return nil
}

func (repo *ConfigGroupInMemRepository) Get(name string, version int) (model.ConfigGroup, error) {
	if versions, ok := repo.configGroups[name]; ok {
		if group, exists := versions[version]; exists {
			return group, nil
		}
		return model.ConfigGroup{}, errors.New("config group version not found")
	}
	return model.ConfigGroup{}, errors.New("config group not found")
}

func (repo *ConfigGroupInMemRepository) Delete(name string, version int) error {
	if versions, exists := repo.configGroups[name]; exists {
		if _, ok := versions[version]; ok {
			delete(versions, version)
			// Ako nema više verzija, ukloni i ime
			if len(versions) == 0 {
				delete(repo.configGroups, name)
			}
			return nil
		}
		return errors.New("config group version not found")
	}
	return errors.New("config group not found")
}

// Ove metode (AddConfigToGroup i RemoveConfigFromGroup) su primeri i mogu zahtevati dodatnu implementaciju
func (repo *ConfigGroupInMemRepository) AddConfigToGroup(groupName string, version int, config model.ConfigWithLabels) error {
	// Implementacija dodavanja konfiguracije u grupu
	group, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}
	// Dodavanje konfiguracije u grupu
	group.Configs[config.Name] = &config
	// Ažuriranje grupe
	return repo.Add(group)
}

func (repo *ConfigGroupInMemRepository) RemoveConfigFromGroup(groupName string, version int, configName string) error {
	// Implementacija uklanjanja konfiguracije iz grupe
	group, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}
	// Uklanjanje konfiguracije iz grupe
	delete(group.Configs, configName)
	// Ažuriranje grupe
	return repo.Add(group)
}
