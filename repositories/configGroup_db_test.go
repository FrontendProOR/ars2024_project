package repositories

import (
	"project/data"
	"project/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGroupDBRepository_Add_Get_Delete(t *testing.T) {
	// Create a new database instance
	db, err := data.NewDatabase()
	assert.NoError(t, err)

	// Create a new ConfigGroupDBRepository instance
	repo := NewConfigGroupDBRepository(db)

	// Create a new config group
	configGroup := model.ConfigGroup{
		Name:    "test-group",
		Version: "1.0",
		Configs: []*model.ConfigWithLabels{
			{
				Config: model.Config{
					Name:    "config1",
					Version: "1.0",
					Params:  map[string]string{"key1": "value1"},
				},
				Labels: []model.Label{
					{Key: "label1", Value: "value1"},
					{Key: "label2", Value: "value2"},
				},
			},
			{
				Config: model.Config{
					Name:    "config2",
					Version: "1.0",
					Params:  map[string]string{"key2": "value2"},
				},
				Labels: []model.Label{
					{Key: "label3", Value: "value3"},
					{Key: "label4", Value: "value4"},
				},
			},
		},
	}

	// Add the config group to the repository
	err = repo.Add(configGroup)
	assert.NoError(t, err)

	// Retrieve the config group from the repository
	retrievedConfigGroup, err := repo.Get(configGroup.Name, configGroup.Version)
	assert.NoError(t, err)

	// Check if the retrieved config group is the same as the original config group
	assert.Equal(t, configGroup, retrievedConfigGroup)

	// Delete the config group from the repository
	err = repo.Delete(configGroup.Name, configGroup.Version)
	assert.NoError(t, err)
}
