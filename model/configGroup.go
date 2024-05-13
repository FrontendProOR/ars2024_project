package model

type ConfigWithLabels struct {
	Config
	Labels map[string]string `json:"labels"`
}

type ConfigGroup struct {
	Name    string              `json:"name"`
	Version int                 `json:"version"`
	Configs []*ConfigWithLabels `json:"configs"`
}

type ConfigGroupRepository interface {
	Add(configGroup ConfigGroup) error
	Get(name string, version int) (ConfigGroup, error)
	Delete(name string, version int) error
	AddConfigToGroup(groupName string, version int, configName string, configVersion int) error
	RemoveConfigFromGroup(groupName string, version int, configName string, configVersion int) error
	GetConfigsByLabels(groupName string, versionInt int, labels map[string]string) (Config, error)
	RemoveConfigsByLabels(groupName string, versionInt int, labels map[string]string) error
}
