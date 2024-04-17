package repositories

import (
	"errors"
	"project/model"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

func NewConfigInMemRepository() model.ConfigRepository {
	return &ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

func (r *ConfigInMemRepository) Add(config model.Config) error {
	// Provera da li konfiguracija već postoji i vraćanje greške ako je imutabilnost pravilo
	if _, exists := r.configs[config.Name]; exists {
		return errors.New("config is immutable and cannot be modified")
	}
	r.configs[config.Name] = config
	return nil
}

func (repo ConfigInMemRepository) Get(name string, version int) (model.Config, error) {
	config, ok := repo.configs[name]
	if !ok {
		return model.Config{}, errors.New("config not found")
	}
	return config, nil
}

func (repo ConfigInMemRepository) Delete(name string, version int) error {
	_, ok := repo.configs[name]
	if !ok {
		return errors.New("config not found")
	}
	delete(repo.configs, name)
	return nil
}
