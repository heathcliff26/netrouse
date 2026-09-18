package config

import (
	"log/slog"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/server/storage"
	"github.com/heathcliff26/netrouse/pkg/server/storage/file"
	"github.com/heathcliff26/netrouse/pkg/server/storage/valkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidConfigs(t *testing.T) {
	tMatrix := []struct {
		Name, Path string
		Result     Config
	}{
		{
			Name:   "EmptyConfig",
			Path:   "",
			Result: DefaultConfig(),
		},
		{
			Name: "ValidConfig1",
			Path: "testdata/valid-config-1.yaml",
			Result: Config{
				LogLevel: "debug",
				Server: ServerConfig{
					Port: 1234,
				},
				Storage: storage.NewDefaultStorageConfig(),
			},
		},
		{
			Name: "ValidConfig2",
			Path: "testdata/valid-config-2.yaml",
			Result: Config{
				LogLevel: "error",
				Server: ServerConfig{
					Port: 5678,
					SSL: SSLConfig{
						Enabled: true,
						Key:     "test.key",
						Cert:    "test.crt",
					},
				},
				Storage: storage.NewDefaultStorageConfig(),
			},
		},
		{
			Name: "ValidConfigWithValkey",
			Path: "testdata/valid-config-valkey.yaml",
			Result: Config{
				LogLevel: "info",
				Server: ServerConfig{
					Port: DEFAULT_SERVER_PORT,
				},
				Storage: storage.StorageConfig{
					Type: "valkey",
					File: file.NewDefaultFileBackendConfig(),
					Valkey: valkey.ValkeyConfig{
						Addrs:     []string{"localhost:6379"},
						Username:  "user",
						Password:  "pass",
						DB:        1,
						TLS:       true,
						Sentinel:  true,
						MasterSet: "none",
					},
				},
			},
		},
		{
			Name: "ValidConfigFileBackend",
			Path: "testdata/valid-config-file-backend.yaml",
			Result: Config{
				LogLevel: "warn",
				Server: ServerConfig{
					Port: 8080,
				},
				Storage: storage.StorageConfig{
					Type:     "file",
					Readonly: true,
					File: file.FileBackendConfig{
						Path: "/data/storage",
					},
				},
			},
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			c, err := LoadConfig(tCase.Path, false, "")

			require.NoError(t, err, "Should not return an error")
			assert.Equal(t, tCase.Result, c, "The config should match the expected result")
		})
	}
}

func TestInvalidConfig(t *testing.T) {
	tMatrix := []struct {
		Name, Path, ErrorMsg string
	}{
		{
			Name:     "NotYaml",
			Path:     "testdata/not-a-config.txt",
			ErrorMsg: "failed to unmarshal config file",
		},
		{
			Name:     "FileDoesNotExist",
			Path:     "not-a-file",
			ErrorMsg: "failed to read config file",
		},
		{
			Name:     "InvalidLogLevel",
			Path:     "testdata/invalid-config-loglevel.yaml",
			ErrorMsg: "failed to set log level",
		},
		{
			Name:     "ServerIncompleteSSLConfig1",
			Path:     "testdata/invalid-config-ssl-1.yaml",
			ErrorMsg: "incomplete SSL configuration",
		},
		{
			Name:     "ServerIncompleteSSLConfig2",
			Path:     "testdata/invalid-config-ssl-2.yaml",
			ErrorMsg: "incomplete SSL configuration",
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			require := require.New(t)

			_, err := LoadConfig(tCase.Path, false, "")

			require.Error(err, "Should return an error")
			require.Contains(err.Error(), tCase.ErrorMsg, "Should return the correct error")
		})
	}
}

func TestEnvSubstitution(t *testing.T) {
	tMatrix := []struct {
		Name   string
		Env    bool
		Config Config
	}{
		{
			Name: "Enabled",
			Env:  true,
			Config: Config{
				LogLevel: "debug",
				Server: ServerConfig{
					Port: 1234,
				},
				Storage: storage.NewDefaultStorageConfig(),
			},
		},
		{
			Name: "Disabled",
			Env:  false,
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			t.Setenv("NETROUSE_LOG_LEVEL", "debug")
			t.Setenv("NETROUSE_PORT", "1234")

			c, err := LoadConfig("testdata/env-config.yaml", tCase.Env, "")
			if tCase.Env {
				assert.NoError(err, "Should not return an error")
				assert.Equal(tCase.Config, c, "The config should match the expected result")
			} else {
				assert.Error(err, "Should return an error")
			}
		})
	}
}

func TestLogLevelOverride(t *testing.T) {
	assert := assert.New(t)

	_, err := LoadConfig("testdata/valid-config-1.yaml", false, "")
	assert.NoError(err, "Should not return an error")
	assert.Equal(logLevel.Level(), slog.LevelDebug, "Should use level from config")

	_, err = LoadConfig("testdata/valid-config-1.yaml", false, "warn")
	assert.NoError(err, "Should not return an error")
	assert.Equal(logLevel.Level(), slog.LevelWarn, "Should override the log level")
}

func TestGetPath(t *testing.T) {
	tMatrix := []struct {
		Name, Path string
		Container  bool
		Result     string
	}{
		{
			Name:   "GivenPath",
			Path:   "testpath",
			Result: "testpath",
		},
		{
			Name:   "Default",
			Result: DEFAULT_CONFIG_PATH,
		},
		{
			Name:      "Container",
			Container: true,
			Result:    DEFAULT_CONFIG_PATH_CONTAINER,
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			if tCase.Container {
				t.Setenv("container", "podman")
			}
			assert.Equal(t, tCase.Result, getPath(tCase.Path), "Should return the correct path")
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	tMatrix := []struct {
		Name        string
		Level       slog.Level
		ShouldError bool
	}{
		{"debug", slog.LevelDebug, false},
		{"info", slog.LevelInfo, false},
		{"warn", slog.LevelWarn, false},
		{"error", slog.LevelError, false},
		{"DEBUG", slog.LevelDebug, false},
		{"INFO", slog.LevelInfo, false},
		{"WARN", slog.LevelWarn, false},
		{"ERROR", slog.LevelError, false},
		{"Unknown", 0, true},
	}
	t.Cleanup(func() {
		err := setLogLevel(DEFAULT_LOG_LEVEL)
		if err != nil {
			t.Fatalf("Failed to cleanup after test: %v", err)
		}
	})

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			err := setLogLevel(tCase.Name)

			assert := assert.New(t)

			if tCase.ShouldError {
				assert.Error(err, "Should return an error")
			} else {
				assert.NoError(err, "Should not return an error")
				assert.Equal(tCase.Level, logLevel.Level(), "Should set the correct log level")
			}
		})
	}
}
