// The code defines a ConfigDBRepository struct that implements model.ConfigRepository interface
// methods for adding, getting, and deleting configuration data in a database.
package repositories

import (
	"fmt"
	"project/data"
	"project/model"
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
func (r *ConfigDBRepository) Add(config model.Config) error {
	_, err := r.db.Put("configs", config.Name, config.Version, config)
	return err
}

// Get retrieves a configuration from the database based on the name and version.
func (r *ConfigDBRepository) Get(name string, version string) (model.Config, error) {
	var config model.Config
	err := r.db.Get(fmt.Sprintf("configs/%s/%s", name, version), &config)
	if err != nil {
		return model.Config{}, err
	}
	return config, nil
}

// Delete deletes a configuration from the database based on the name and version.
func (r *ConfigDBRepository) Delete(name string, version string) error {
	return r.db.Delete(fmt.Sprintf("configs/%s/%s", name, version))
}
