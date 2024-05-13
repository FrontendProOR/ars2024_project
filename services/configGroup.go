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

func (s ConfigGroupService) Add(group model.ConfigGroup) error {
	// Provera validnosti ConfigGroup objekta
	// Na primer, provera da li grupa već postoji, da li su svi potrebni podaci prisutni itd.

	// Dodavanje grupe u repozitorijum
	return s.repo.Add(group)
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

func (s ConfigGroupService) AddConfigToGroup(groupName string, version int, configName string, configVersion int) error {
	// Implementacija bi uključivala dobijanje postojeće grupe,
	// dodavanje konfiguracije u grupu, i ažuriranje grupe u repozitorijumu.
	// Ovo zahteva dodatnu logiku i možda izmene u model.ConfigGroupRepository interfejsu.
	return s.repo.AddConfigToGroup(groupName, version, configName, configVersion)
}

func (s ConfigGroupService) RemoveConfigFromGroup(groupName string, groupVersion int, configName string, configVersion int) error {
	return s.repo.RemoveConfigFromGroup(groupName, groupVersion, configName, configVersion)
}

func (s ConfigGroupService) GetConfigsByLabels(groupName string, versionInt int, labels map[string]string) (model.Config, error) {
	configs, err := s.repo.GetConfigsByLabels(groupName, versionInt, labels)
	if err != nil {
		return model.Config{}, err // it returns empty config in case of error
	}
	return configs, nil
}

func (s ConfigGroupService) RemoveConfigsByLabels(groupName string, versionInt int, labels map[string]string) error {
	err := s.repo.RemoveConfigsByLabels(groupName, versionInt, labels)
	if err != nil {
		return err
	}
	return nil
}
