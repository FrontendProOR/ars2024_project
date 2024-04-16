package model

type Config struct {
	ID      string            // Unique identifier
	Name    string            // Name of the configuration
	Version string            // Version of the configuration
	Params  map[string]string // Key-value pairs of configuration parameters
}

type ConfigGroup struct {
	ID      string            // Unique identifier
	Name    string            // Name of the configuration group
	Version string            // Version of the configuration group
	Configs []Config          // List of configurations
	Labels  map[string]string // Key-value pairs of labels
}

type ConfigRepository interface {
	AddConfig(config Config) error
	GetConfig(id string) (Config, error)
	DeleteConfig(id string) error
	AddConfigGroup(group ConfigGroup) error
	GetConfigGroup(id string) (ConfigGroup, error)
	DeleteConfigGroup(id string) error
	UpdateConfigGroupAddConfig(groupId string, config Config) error
	UpdateConfigGroupRemoveConfig(groupId string, configId string) error
}
