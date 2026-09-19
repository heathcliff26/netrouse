package persistence

import (
	"errors"
	"io/fs"
	"os"

	"go.yaml.in/yaml/v3"
)

type Settings struct {
	Remotes []RemoteServer `yaml:"remotes"`
}

type RemoteServer struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Load settings from the settings file
func LoadSettings() (settings *Settings, err error) {
	settings = &Settings{}

	buf, err := os.ReadFile(SettingsFile())
	if errors.Is(err, fs.ErrNotExist) {
		err = nil
		return
	} else if err != nil {
		return
	}

	err = yaml.Unmarshal(buf, settings)
	if err != nil {
		return
	}
	return
}

func (s *Settings) Save() error {
	buf, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(SettingsFile(), buf, 0644)
}
