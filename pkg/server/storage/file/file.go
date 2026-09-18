package file

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"go.yaml.in/yaml/v3"
)

type FileBackend struct {
	path    string
	storage *types.HostsFile
	lock    sync.RWMutex
}

func NewFileBackend(cfg FileBackendConfig) (*FileBackend, error) {
	fb := &FileBackend{
		path: cfg.Path,
		storage: &types.HostsFile{
			Hosts: []types.Host{},
		},
	}

	f, err := os.ReadFile(cfg.Path)
	if os.IsNotExist(err) {
		slog.Info("File not found, creating new file", slog.String("path", cfg.Path))
		err := fb.save()
		if err != nil {
			return nil, fmt.Errorf("failed to create new storage file: %w", err)
		}
		return fb, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read storage file: %w", err)
	}

	err = yaml.Unmarshal(f, fb.storage)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal storage file: %w", err)
	}

	slog.Debug("Storage file loaded", slog.String("path", cfg.Path))

	// Ensure that all MAC addresses are uppercase and unique
	macSet := make(map[string]bool, len(fb.storage.Hosts))
	uniqueHosts := make([]types.Host, 0, len(fb.storage.Hosts))
	changed := false
	for _, host := range fb.storage.Hosts {
		uppercaseMAC := strings.ToUpper(host.MAC)
		if host.MAC != uppercaseMAC {
			changed = true
		}
		if !macSet[uppercaseMAC] {
			macSet[uppercaseMAC] = true
			host.MAC = uppercaseMAC
			uniqueHosts = append(uniqueHosts, host)
		} else {
			changed = true
		}
	}
	fb.storage.Hosts = uniqueHosts

	if changed {
		err := fb.save()
		if err != nil {
			return nil, fmt.Errorf("failed to save storage file after ensuring unique, uppercase MAC addresses: %w", err)
		}
	}

	return fb, nil
}

// Add a new host, overwrite existing host name if it already exists.
// Ensures that the MAC address is unique and uppercase.
func (fb *FileBackend) AddHost(host types.Host) error {
	fb.lock.Lock()
	defer fb.lock.Unlock()

	host.MAC = strings.ToUpper(host.MAC)
	for i, h := range fb.storage.Hosts {
		if h.MAC == host.MAC {
			fb.storage.Hosts[i] = host
			return fb.save()
		}
	}

	fb.storage.Hosts = append(fb.storage.Hosts, host)
	return fb.save()
}

// Remove a host, ignore if the host does not exist
func (fb *FileBackend) RemoveHost(mac string) error {
	fb.lock.Lock()
	defer fb.lock.Unlock()

	uppercaseMAC := strings.ToUpper(mac)
	for i, host := range fb.storage.Hosts {
		if host.MAC == uppercaseMAC {
			fb.storage.Hosts = append(fb.storage.Hosts[:i], fb.storage.Hosts[i+1:]...)
			return fb.save()
		}
	}
	return nil
}

// Return the host name for a given MAC address, return empty if not found
func (fb *FileBackend) GetHost(mac string) (types.Host, error) {
	fb.lock.RLock()
	defer fb.lock.RUnlock()

	uppercaseMAC := strings.ToUpper(mac)
	for _, host := range fb.storage.Hosts {
		if host.MAC == uppercaseMAC {
			return host, nil
		}
	}
	return types.Host{}, nil
}

// Return all hosts
func (fb *FileBackend) GetHosts() ([]types.Host, error) {
	fb.lock.RLock()
	defer fb.lock.RUnlock()

	return append([]types.Host{}, fb.storage.Hosts...), nil
}

// Check if the storage backend is readonly
func (fb *FileBackend) Readonly() (bool, error) {
	//#nosec G302 -- The file does not contain sensitive data, so it can be world readable. Additionally final permissions are determined by the umask.
	f, err := os.OpenFile(fb.path, os.O_RDWR, 0644)
	if err != nil {
		return true, nil
	}
	_ = f.Close()
	return false, nil
}

func (fb *FileBackend) save() error {
	data, err := yaml.Marshal(fb.storage)
	if err != nil {
		return fmt.Errorf("failed to marshal storage data: %w", err)
	}

	// #nosec G306 -- The file does not contain sensitive data, so it can be world readable. Additionally final permissions are determined by the umask.
	err = os.WriteFile(fb.path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write storage file: %w", err)
	}

	slog.Debug("Storage file saved", slog.String("path", fb.path))
	return nil
}
