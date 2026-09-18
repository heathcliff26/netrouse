package persistence

import (
	"log/slog"
	"path/filepath"
)

var configFolder string

const (
	hostsFileName = "hosts.yaml"
)

func init() {
	err := initConfigFolder()
	if err != nil {
		slog.Error("Failed to find config location, defaulting to current directory", "error", err)
		configFolder = "."
	}
}

// Path to hosts file for local hosts
func HostsFile() string {
	return filepath.Join(configFolder, hostsFileName)
}

// Path to config folder
func ConfigFolder() string {
	return configFolder
}

// Set the config folder path, used for testing
func SetConfigFolder(path string) {
	configFolder = path
}
