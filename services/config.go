package services

import (
	"fmt"
	"project/model"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

func (s ConfigService) Hello() {
	fmt.Println("hello from config service")
}

func (s ConfigService) AddConfig(config model.Config) error {
	return s.repo.AddConfig(config)
}

func (s ConfigService) GetConfig(id string) (model.Config, error) {
	return s.repo.GetConfig(id)
}

func (s ConfigService) DeleteConfig(id string) error {
	return s.repo.DeleteConfig(id)
}

func (s ConfigService) AddConfigGroup(group model.ConfigGroup) error {
	return s.repo.AddConfigGroup(group)
}

func (s ConfigService) GetConfigGroup(id string) (model.ConfigGroup, error) {
	return s.repo.GetConfigGroup(id)
}

func (s ConfigService) DeleteConfigGroup(id string) error {
	return s.repo.DeleteConfigGroup(id)
}
