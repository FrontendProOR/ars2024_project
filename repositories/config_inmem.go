package repositories

import (
	"errors"
	"project/model"
)

type ConfigInMemRepository struct {
	configs map[string]map[int]model.Config
}

func NewConfigInMemRepository() model.ConfigRepository {
	return &ConfigInMemRepository{
		configs: make(map[string]map[int]model.Config),
	}
}

func (r *ConfigInMemRepository) Add(config model.Config) error {
	if _, exists := r.configs[config.Name]; !exists {
		r.configs[config.Name] = make(map[int]model.Config)
	}
	// Provera da li postoji verzija konfiguracije
	if _, exists := r.configs[config.Name][config.Version]; exists {
		return errors.New("config version already exists")
	}
	r.configs[config.Name][config.Version] = config
	return nil
}

func (repo *ConfigInMemRepository) Get(name string, version int) (model.Config, error) {
	if versions, ok := repo.configs[name]; ok {
		if config, exists := versions[version]; exists {
			return config, nil
		}
		return model.Config{}, errors.New("config version not found")
	}
	return model.Config{}, errors.New("config not found")
}

func (repo *ConfigInMemRepository) Delete(name string, version int) error {
	if versions, ok := repo.configs[name]; ok {
		if _, exists := versions[version]; exists {
			delete(versions, version)
			// Ako nema više verzija, ukloni i ime
			if len(versions) == 0 {
				delete(repo.configs, name)
			}
			return nil
		}
		return errors.New("config version not found")
	}
	return errors.New("config not found")
}
