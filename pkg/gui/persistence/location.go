package persistence

import (
	"path/filepath"
)

var configFolder string

const (
	hostsFileName    = "hosts.yaml"
	settingsFileName = "settings.yaml"
)

func Init() {
	if configFolder == "" {
		initConfigFolder()
	}
}

// Path to hosts file for local hosts
func HostsFile() string {
	return filepath.Join(configFolder, hostsFileName)
}

// Path to the settings file
func SettingsFile() string {
	return filepath.Join(configFolder, settingsFileName)
}

// Path to config folder
func ConfigFolder() string {
	return configFolder
}

// Set the config folder path, used for testing
func SetConfigFolder(path string) {
	configFolder = path
}
