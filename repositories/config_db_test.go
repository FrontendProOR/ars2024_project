// The TestConfigDBRepository_Get function tests the functionality of retrieving a configuration from a
// database using a ConfigDBRepository instance.
package repositories

import (
	"testing"

	"project/data"
	"project/model"

	"github.com/stretchr/testify/assert"
)

func TestConfigDBRepository_Get(t *testing.T) {
	// Create a new database instance
	db, err := data.NewDatabase()
	assert.NoError(t, err)

	// Create a new ConfigDBRepository instance
	repo := NewConfigDBRepository(db)

	// Add a configuration to the database
	config := model.Config{
		Name:    "test",
		Version: "1.0",
		// Add other fields as needed
		Params: map[string]string{
			"param1": "value1",
			"param2": "value2",
		},
	}
	err = repo.Add(config)
	assert.NoError(t, err)

	// Get the configuration from the database
	retrievedConfig, err := repo.Get(config.Name, config.Version)
	assert.NoError(t, err)

	// Check if the retrieved configuration is the same as the original configuration
	assert.Equal(t, config, retrievedConfig)
}
