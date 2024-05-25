// Package model defines the Config struct and its repository interface.
//
// Config holds a name, version, and parameters.
// ConfigRepository outlines the required methods for a config repository.
package model

type Config struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Params  map[string]string `json:"params"`
}

type ConfigRepository interface {
	Add(config Config) error
	Get(name string, version string) (Config, error)
	Delete(name string, version string) error
}
