package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/heathcliff26/netrouse/pkg/server/storage/file"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/heathcliff26/netrouse/pkg/server/storage/valkey"
	"github.com/heathcliff26/netrouse/pkg/version"
	"github.com/heathcliff26/netrouse/static"

	"go.yaml.in/yaml/v3"
)

type Storage struct {
	backend  types.StorageBackend
	readonly bool
}

func NewStorage(cfg StorageConfig) (*Storage, error) {
	var backend types.StorageBackend
	var err error
	switch cfg.Type {
	case "file":
		backend, err = file.NewFileBackend(cfg.File)
	case "valkey":
		backend, err = valkey.NewValkeyBackend(cfg.Valkey)
	default:
		return nil, fmt.Errorf("unknown storage backend type: %s", cfg.Type)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create storage backend: %w", err)
	}

	s := &Storage{
		backend:  backend,
		readonly: cfg.Readonly,
	}

	if !s.readonly {
		s.readonly, err = s.backend.Readonly()
		if err != nil {
			return nil, fmt.Errorf("failed to check if storage backend is readonly: %w", err)
		}
	}

	if cfg.SeededHosts != "" {
		if s.readonly {
			return nil, fmt.Errorf("cannot seed hosts in readonly mode")
		}

		f, err := os.ReadFile(cfg.SeededHosts)
		if err != nil {
			return nil, fmt.Errorf("failed to read seeded hosts file: %w", err)
		}
		var seededHosts types.HostsFile
		err = yaml.Unmarshal(f, &seededHosts)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal seeded hosts file: %w", err)
		}

		for _, host := range seededHosts.Hosts {
			slog.Debug("Adding seeded host", "mac", host.MAC, "name", host.Name)
			err := s.backend.AddHost(host)
			if err != nil {
				return nil, fmt.Errorf("failed to add seeded host '%s': %w", host.MAC, err)
			}
		}
	}

	return s, nil
}

// Create a new storage backend from the given backend
// This is used for testing purposes
func NewStorageFromBackend(backend types.StorageBackend, readonly bool) *Storage {
	return &Storage{
		backend: backend,
	}
}

type indexValues struct {
	Readonly bool
	Hosts    []types.Host
	Version  string
	Name     string
}

// Generate the index.html file from the template and the current hosts.
// Returns the generated HTML and its checksum.
func (s *Storage) GetIndexHTML() (string, string, error) {
	hosts, err := s.backend.GetHosts()
	if err != nil {
		return "", "", fmt.Errorf("failed to get hosts: %w", err)
	}

	values := indexValues{
		Readonly: s.readonly,
		Hosts:    hosts,
		Version:  version.Version(),
		Name:     version.Name,
	}

	tmpl, err := template.New("index.html").Parse(string(static.IndexTemplate))
	if err != nil {
		return "", "", fmt.Errorf("unexpected error when creating a template from static html: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return "", "", fmt.Errorf("failed to apply values to html template: %w", err)
	}
	indexHTML := buf.String()

	checksum := sha256.Sum256([]byte(indexHTML))

	return indexHTML, hex.EncodeToString(checksum[:]), nil
}

// Return if the storage is readonly
func (s *Storage) Readonly() bool {
	return s.readonly
}

// Get all hosts from the storage
func (s *Storage) GetHosts() ([]types.Host, error) {
	hosts, err := s.backend.GetHosts()
	if err != nil {
		return nil, fmt.Errorf("failed to get hosts: %w", err)
	}
	return hosts, nil
}

// Add a new host and update the index.html
func (s *Storage) AddHost(host types.Host) error {
	if s.readonly {
		return fmt.Errorf("storage is readonly")
	}

	err := s.backend.AddHost(host)
	if err != nil {
		return fmt.Errorf("failed to add host: %w", err)
	}

	return nil
}

// Remove a host and update the index.html
func (s *Storage) RemoveHost(mac string) error {
	if s.readonly {
		return fmt.Errorf("storage is readonly")
	}

	err := s.backend.RemoveHost(mac)
	if err != nil {
		return fmt.Errorf("failed to remove host: %w", err)
	}

	return nil
}
