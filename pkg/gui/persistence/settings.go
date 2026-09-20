package persistence

import (
	"errors"
	"io/fs"
	"os"

	"fyne.io/fyne/v2"
	"go.yaml.in/yaml/v3"
)

const (
	DefaultWindowWidth  = 450
	DefaultWindowHeight = 600
)

type Settings struct {
	SelectedTab int            `yaml:"selectedTab"`
	FullScreen  bool           `yaml:"fullScreen"`
	WindowSize  Size           `yaml:"windowSize"`
	Remotes     []RemoteServer `yaml:"remotes"`
}

type RemoteServer struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Size struct {
	Width  float32 `yaml:"width"`
	Height float32 `yaml:"height"`
}

func DefaultSettings() *Settings {
	return &Settings{
		WindowSize: Size{
			Width:  DefaultWindowWidth,
			Height: DefaultWindowHeight,
		},
	}
}

// Load settings from the settings file
func LoadSettings() (settings *Settings, err error) {
	settings = DefaultSettings()

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

func SizeFromFyne(s fyne.Size) Size {
	return Size{
		Width:  s.Width,
		Height: s.Height,
	}
}

func (s Size) ToFyne() fyne.Size {
	return fyne.NewSize(s.Width, s.Height)
}
