// Package model defines the Config struct and its repository interface.
//
// Config represents a configuration entity with a name, version, and parameters.
// ConfigRepository outlines the required methods for a config repository.
package model

// Config represents a configuration entity
// @Description A configuration entity with a name, version, and parameters
type Config struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Params  map[string]string `json:"params"`
}

// ConfigRepository outlines the required methods for a config repository
// @Description Interface defining methods for a config repository
type ConfigRepository interface {
	Add(config Config) error
	Get(name string, version string) (Config, error)
	Delete(name string, version string) error
}
