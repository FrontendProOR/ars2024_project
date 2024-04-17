package repositories

import (
	"errors"
	"fmt"
	"project/model"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigInMemRepository() model.ConfigRepository {
	return &ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

// Add implements model.ConfigRepository.
func (c ConfigInMemRepository) Add(config model.Config) {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	c.configs[key] = config
}

// Get implements model.ConfigRepository.
func (c ConfigInMemRepository) Get(name string, version int) (model.Config, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Config{}, errors.New("config not found")
	}
	return config, nil
}

func (c ConfigInMemRepository) Delete(name string, version int) error {
	key := fmt.Sprintf("%s/%d", name, version)
	if _, ok := c.configs[key]; !ok {
		return errors.New("config not found")
	}
	delete(c.configs, key)
	return nil
}
