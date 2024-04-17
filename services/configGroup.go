package services

import "project/model"

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupService(repo model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{
		repo: repo,
	}
}

func (s ConfigGroupService) Add(ConfigGroup model.ConfigGroup) error {
	return s.repo.Add(ConfigGroup)
}

func (s ConfigGroupService) Get(name string, version int) (model.ConfigGroup, error) {
	return s.repo.Get(name, version)
}

func (s ConfigGroupService) Delete(name string, version int) error {
	err := s.repo.Delete(name, version)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigGroupService) AddConfigToGroup(groupName string, version int, config model.Config) error {
	// Implementacija bi uključivala dobijanje postojeće grupe,
	// dodavanje konfiguracije u grupu, i ažuriranje grupe u repozitorijumu.
	// Ovo zahteva dodatnu logiku i možda izmene u model.ConfigGroupRepository interfejsu.
	return nil
}

func (s ConfigGroupService) RemoveConfigFromGroup(groupName string, version int, configName string) error {
	// Slično kao i za dodavanje, ali uklanjanje konfiguracije iz grupe.
	return nil
}
