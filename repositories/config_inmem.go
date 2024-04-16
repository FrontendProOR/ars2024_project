package repositories

import (
	"errors"
	"project/model"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
	groups  map[string]model.ConfigGroup
}

func (repo *ConfigInMemRepository) AddConfig(config model.Config) error {
	repo.configs[config.ID] = config
	return nil
}

func (repo *ConfigInMemRepository) GetConfig(id string) (model.Config, error) {
	config, ok := repo.configs[id]
	if !ok {
		return model.Config{}, errors.New("config not found")
	}
	return config, nil
}

func (repo *ConfigInMemRepository) DeleteConfig(id string) error {
	_, ok := repo.configs[id]
	if !ok {
		return errors.New("config not found")
	}
	delete(repo.configs, id)
	return nil
}

func (repo *ConfigInMemRepository) AddConfigGroup(group model.ConfigGroup) error {
	repo.groups[group.ID] = group
	return nil
}

func (repo *ConfigInMemRepository) GetConfigGroup(id string) (model.ConfigGroup, error) {
	group, ok := repo.groups[id]
	if !ok {
		return model.ConfigGroup{}, errors.New("group not found")
	}
	return group, nil
}

func (repo *ConfigInMemRepository) DeleteConfigGroup(id string) error {
	_, ok := repo.groups[id]
	if !ok {
		return errors.New("group not found")
	}
	delete(repo.groups, id)
	return nil
}

func (repo *ConfigInMemRepository) UpdateConfigGroupAddConfig(groupId string, config model.Config) error {
	group, ok := repo.groups[groupId]
	if !ok {
		return errors.New("group not found")
	}
	group.Configs = append(group.Configs, config)
	repo.groups[groupId] = group
	return nil
}

func (repo *ConfigInMemRepository) UpdateConfigGroupRemoveConfig(groupId string, configId string) error {
	group, ok := repo.groups[groupId]
	if !ok {
		return errors.New("group not found")
	}
	for i, config := range group.Configs {
		if config.ID == configId {
			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
			repo.groups[groupId] = group
			return nil
		}
	}
	return errors.New("config not found in group")
}

func NewConfigInMemRepository() model.ConfigRepository {
	return &ConfigInMemRepository{
		configs: make(map[string]model.Config),
		groups:  make(map[string]model.ConfigGroup),
	}
}
