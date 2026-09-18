package storage

import (
	"github.com/heathcliff26/netrouse/pkg/server/storage/file"
	"github.com/heathcliff26/netrouse/pkg/server/storage/valkey"
)

const (
	DEFAULT_READONLY     = false
	DEFAULT_BACKEND_TYPE = "file"
)

type StorageConfig struct {
	Type        string                 `yaml:"type"`
	Readonly    bool                   `yaml:"readonly,omitempty"`
	SeededHosts string                 `yaml:"seeded-hosts,omitempty"`
	File        file.FileBackendConfig `yaml:"file,omitempty"`
	Valkey      valkey.ValkeyConfig    `yaml:"valkey,omitempty"`
}

func NewDefaultStorageConfig() StorageConfig {
	return StorageConfig{
		Type:     DEFAULT_BACKEND_TYPE,
		Readonly: DEFAULT_READONLY,
		File:     file.NewDefaultFileBackendConfig(),
	}
}
