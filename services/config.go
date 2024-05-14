package services

import (
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

func (s ConfigService) Add(config model.Config) error {
	return s.repo.Add(config)
}

func (s ConfigService) Get(name string, version string) (model.Config, error) {
	return s.repo.Get(name, version)
}

func (s ConfigService) Delete(name string, version string) error {
	err := s.repo.Delete(name, version)
	if err != nil {
		return err
	}
	return nil
}
