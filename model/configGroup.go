// Package model defines the ConfigGroup struct and its repository interface.
//
// ConfigGroup holds a name, version, and a list of ConfigWithLabels.
// ConfigWithLabels is a Config with an additional Labels field.
// Label represents a key-value pair.
// ConfigGroupRepository outlines the required methods for a config group repository.
package model

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ConfigWithLabels struct {
	Config
	Labels []Label `json:"labels"`
}

type ConfigGroup struct {
	Name    string              `json:"name"`
	Version string              `json:"version"`
	Configs []*ConfigWithLabels `json:"configs"`
}

type ConfigGroupRepository interface {
	Add(configGroup ConfigGroup) error
	Get(name string, version string) (ConfigGroup, error)
	Delete(name string, version string) error
	AddConfigToGroup(groupName string, version string, configName string, configVersion string) error
	RemoveConfigFromGroup(groupName string, version string, configName string, configVersion string) error
	AddConfigWithLabelToGroup(groupName string, version string, config ConfigWithLabels) error
	SearchConfigsWithLabelsInGroup(groupName string, version string, labels []Label) ([]*ConfigWithLabels, error)
	RemoveConfigsWithLabelsFromGroup(groupName string, version string, labels []Label) error
}
