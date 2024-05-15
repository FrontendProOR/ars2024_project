package repositories

import (
	"errors"
	"project/model"
)

type ConfigGroupInMemRepository struct {
	configGroups map[string]map[string]model.ConfigGroup
	configRepo   model.ConfigRepository
}

func NewConfigGroupInMemRepository(configRepo model.ConfigRepository) *ConfigGroupInMemRepository {
	return &ConfigGroupInMemRepository{
		configGroups: make(map[string]map[string]model.ConfigGroup),
		configRepo:   configRepo,
	}
}

func (repo *ConfigGroupInMemRepository) Add(configGroup model.ConfigGroup) error {
	if _, exists := repo.configGroups[configGroup.Name]; !exists {
		repo.configGroups[configGroup.Name] = make(map[string]model.ConfigGroup)
	}
	if _, exists := repo.configGroups[configGroup.Name][configGroup.Version]; exists {
		return errors.New("config group version already exists")
	}
	repo.configGroups[configGroup.Name][configGroup.Version] = configGroup
	return nil
}

func (repo *ConfigGroupInMemRepository) Get(name string, version string) (model.ConfigGroup, error) {
	if versions, ok := repo.configGroups[name]; ok {
		if group, exists := versions[version]; exists {
			return group, nil
		}
		return model.ConfigGroup{}, errors.New("config group version not found")
	}
	return model.ConfigGroup{}, errors.New("config group not found")
}

func (repo *ConfigGroupInMemRepository) Delete(name string, version string) error {
	if versions, exists := repo.configGroups[name]; exists {
		if _, ok := versions[version]; ok {
			delete(versions, version)
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
func (repo *ConfigGroupInMemRepository) AddConfigToGroup(groupName string, version string, configName string, configVersion string) error {
	group, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}
	// Provera da li konfiguracija već postoji pomoću imena i verzije
	for _, config := range group.Configs {
		if config.Name == configName && config.Version == configVersion {
			return errors.New("config already exists in group")
		}
	}
	config, err := repo.configRepo.Get(configName, configVersion)
	if err != nil {
		return err
	}
	// Dodavanje nove konfiguracije u listu
	group.Configs = append(group.Configs, &model.ConfigWithLabels{
		Config: model.Config{
			Name:    configName,
			Version: config.Version,
			Params:  config.Params,
		},
		Labels: []model.Label{},
	})
	// Ažuriranje grupe
	return repo.Add(group)
}

func (repo *ConfigGroupInMemRepository) RemoveConfigFromGroup(groupName string, groupVersion string, configName string, configVersion string) error {
	group, err := repo.Get(groupName, groupVersion)
	if err != nil {
		return err
	}
	for i, config := range group.Configs {
		if config.Name == configName && config.Version == configVersion {
			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
			return repo.Add(group)
		}
	}
	return errors.New("config not found")
}

func (repo *ConfigGroupInMemRepository) AddConfigWithLabelToGroup(groupName string, version string, config model.ConfigWithLabels) error {
	group, err := repo.Get(groupName, version)
	if err != nil {
		return err
	}
	for _, c := range group.Configs {
		if c.Name == config.Name && c.Version == config.Version {
			return errors.New("config already exists in group")
		}
	}
	group.Configs = append(group.Configs, &config)
	return repo.Add(group)
}

func (repo *ConfigGroupInMemRepository) RemoveConfigsWithLabelsFromGroup(groupName string, groupVersion string, labels []model.Label) error {
	group, err := repo.Get(groupName, groupVersion)
	if err != nil {
		return err
	}
	configFound := false
	for i := 0; i < len(group.Configs); {
		config := group.Configs[i]
		if containsAllLabels(config.Labels, labels) {
			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
			configFound = true
		} else {
			i++
		}
	}
	if !configFound {
		return errors.New("config not found with the provided labels")
	}
	return repo.Update(group)
}

func (repo *ConfigGroupInMemRepository) SearchConfigsWithLabelsInGroup(groupName string, version string, searchLabels []model.Label) ([]*model.ConfigWithLabels, error) {
	group, err := repo.Get(groupName, version)
	if err != nil {
		return nil, err
	}
	var result []*model.ConfigWithLabels
	for _, config := range group.Configs {
		if containsLabels(config.Labels, searchLabels) {
			result = append(result, config)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("no configs with specified labels found")
	}
	return result, nil
}

func containsLabels(configLabels []model.Label, searchLabels []model.Label) bool {
	for _, searchLabel := range searchLabels {
		found := false
		for _, configLabel := range configLabels {
			if configLabel.Key == searchLabel.Key && configLabel.Value == searchLabel.Value {
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

func containsAllLabels(configLabels []model.Label, searchLabels []model.Label) bool {
	if len(configLabels) != len(searchLabels) {
		return false
	}

	labelSet := make(map[string]bool)
	for _, label := range configLabels {
		labelSet[label.Key+":"+label.Value] = true
	}

	for _, label := range searchLabels {
		if !labelSet[label.Key+":"+label.Value] {
			return false
		}
	}

	return true
}

func (repo *ConfigGroupInMemRepository) Update(configGroup model.ConfigGroup) error {
	if _, exists := repo.configGroups[configGroup.Name]; !exists {
		return errors.New("config group does not exist")
	}
	repo.configGroups[configGroup.Name][configGroup.Version] = configGroup
	return nil
}
