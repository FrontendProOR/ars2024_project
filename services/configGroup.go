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
	return s.repo.Add(group)
}

func (s ConfigGroupService) Get(name string, version string) (model.ConfigGroup, error) {
	return s.repo.Get(name, version)
}

func (s ConfigGroupService) Delete(name string, version string) error {
	err := s.repo.Delete(name, version)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigGroupService) AddConfigToGroup(groupName string, version string, configName string, configVersion string) error {
	return s.repo.AddConfigToGroup(groupName, version, configName, configVersion)
}

func (s ConfigGroupService) RemoveConfigFromGroup(groupName string, version string, configName string, configVersion string) error {
	return s.repo.RemoveConfigFromGroup(groupName, version, configName, configVersion)
}

func (s ConfigGroupService) AddConfigWithLabelToGroup(groupName string, version string, config model.ConfigWithLabels) error {
	return s.repo.AddConfigWithLabelToGroup(groupName, version, config)
}

func (s ConfigGroupService) RemoveConfigWithLabelFromGroup(groupName string, version string, configName string, configVersion string, label model.Label) error {
	return s.repo.RemoveConfigWithLabelFromGroup(groupName, version, configName, configVersion, label)
}

func (s ConfigGroupService) RemoveConfigsWithLabelsFromGroup(groupName string, version string, labels []model.Label) error {
	return s.repo.RemoveConfigsWithLabelsFromGroup(groupName, version, labels)
}

func (s ConfigGroupService) SearchConfigsWithLabelsInGroup(groupName string, version string, labels []model.Label) ([]*model.ConfigWithLabels, error) {
	return s.repo.SearchConfigsWithLabelsInGroup(groupName, version, labels)
}
