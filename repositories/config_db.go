// The code defines a ConfigDBRepository struct that implements model.ConfigRepository interface
// methods for adding, getting, and deleting configuration data in a database.
package repositories

import (
	"errors"
	"fmt"
	"project/data"
	"project/model"
	"strings"
)

type ConfigDBRepository struct {
	db *data.Database
}

func NewConfigDBRepository(db *data.Database) model.ConfigRepository {
	return &ConfigDBRepository{
		db: db,
	}
}

// Add adds a new configuration to the database.
func (repo *ConfigDBRepository) Add(config model.Config) error {
	// Validation
	if strings.TrimSpace(config.Name) == "" {
		return errors.New("config name cannot be empty")
	}
	if strings.TrimSpace(config.Version) == "" {
		return errors.New("config version cannot be empty")
	}

	// Check if the config already exists
	existingConfig, err := repo.Get(config.Name, config.Version)
	if err == nil && existingConfig.Name != "" && existingConfig.Version != "" {
		return errors.New("config with this name and version already exists")
	}

	// Add the config
	_, err = repo.db.Put("configs", config.Name, config.Version, config)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves a configuration from the database based on the name and version.
func (r *ConfigDBRepository) Get(name string, version string) (model.Config, error) {
	var config model.Config
	err := r.db.Get(fmt.Sprintf("configs/%s/%s", name, version), &config)
	if err != nil {
		return model.Config{}, err
	}

	// Check if the retrieved config is empty
	if config.Name == "" && config.Version == "" && config.Params == nil {
		return model.Config{}, fmt.Errorf("no configuration found with name %s and version %s", name, version)
	}

	return config, nil
}

// Delete deletes a configuration from the database based on the name and version.
func (repo *ConfigDBRepository) Delete(name string, version string) error {
	// Check if the config exists
	_, err := repo.Get(name, version)
	if err != nil {
		// If the config does not exist, return the error
		return err
	}

	// Delete the config
	return repo.db.Delete(fmt.Sprintf("configs/%s/%s", name, version))
}
