package gui

import (
	"os"
	"path/filepath"

	"github.com/heathcliff26/netrouse/pkg/client"
	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
	"github.com/heathcliff26/netrouse/pkg/ping"
	"github.com/heathcliff26/netrouse/pkg/server/storage/file"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/heathcliff26/netrouse/pkg/wol"
)

var _ client.Client = &localClient{}

type localClient struct {
	store *file.FileBackend
}

func newLocalClient() (client.Client, error) {
	path := persistence.HostsFile()
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return &localClient{
			store: file.NewEmptyBackend(),
		}, err
	}
	cfg := file.FileBackendConfig{
		Path: persistence.HostsFile(),
	}
	store, err := file.NewFileBackend(cfg)
	if err != nil {
		store = file.NewEmptyBackend()
	}
	return &localClient{
		store: store,
	}, err
}

// Add a new host, overwrite existing host name if it already exists.
func (lc *localClient) AddHost(host types.Host) error {
	return lc.store.AddHost(host)
}

// Return all hosts
func (lc *localClient) GetHosts() ([]types.Host, error) {
	return lc.store.GetHosts()
}

// Remove a host, ignore if the host does not exist
func (lc *localClient) RemoveHost(mac string) error {
	return lc.store.RemoveHost(mac)
}

// Return the current status of all hosts
func (lc *localClient) Status() ([]types.HostStatus, error) {
	hosts, err := lc.store.GetHosts()
	if err != nil {
		return nil, err
	}
	hostsToCheck := make([]types.Host, 0, len(hosts))
	for _, host := range hosts {
		if host.Address != "" {
			hostsToCheck = append(hostsToCheck, host)
		}
	}
	return ping.PingHosts(hostsToCheck), nil
}

// Send a magic packet to wake a host
func (lc *localClient) Wake(mac string) error {
	packet, err := wol.CreatePacket(mac)
	if err != nil {
		return err
	}
	return packet.Send("")
}
