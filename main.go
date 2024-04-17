package main

import (
	"fmt"
	"project/model"
	"project/repositories"
	"project/services"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)

	config := model.Config{
		ID:      "1",
		Name:    "Test Config",
		Version: "1.0",
		Params:  map[string]string{"param1": "value1"},
	}

	err := service.AddConfig(config)
	if err != nil {
		fmt.Println("Error adding config:", err)
		return
	}

	retrievedConfig, err := service.GetConfig("1")
	if err != nil {
		fmt.Println("Error getting config:", err)
		return
	}

	fmt.Println("Retrieved config:", retrievedConfig)

	err = service.DeleteConfig("1")
	if err != nil {
		fmt.Println("Error deleting config:", err)
		return
	}

	fmt.Println("Config deleted successfully")
}
