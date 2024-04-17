package model

type ConfigWithLabels struct {
	Config
	Labels map[string]string `json:"labels"`
}

type ConfigGroup struct {
	Name    string                       `json:"name"`
	Version int                          `json:"version"`
	Configs map[string]*ConfigWithLabels `json:"configs"`
}

type ConfigGroupRepository interface {
	Add(configGroup ConfigGroup) error
	Get(name string, version int) (ConfigGroup, error)
	Delete(name string, version int) error
	AddConfigToGroup(groupName string, version int, config ConfigWithLabels) error
	RemoveConfigFromGroup(groupName string, version int, configName string) error
}
