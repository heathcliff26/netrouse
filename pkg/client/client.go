package client

import "github.com/heathcliff26/netrouse/pkg/server/storage/types"

var _ Client = &apiClient{}

type Client interface {
	// Return all hosts
	GetHosts() ([]types.Host, error)
	// Add a new host, overwrite existing host name if it already exists.
	AddHost(host types.Host) error
	// Remove a host, ignore if the host does not exist
	RemoveHost(mac string) error
	// Return the current status of all hosts
	Status() ([]types.HostStatus, error)
	// Send a magic packet to wake a host
	Wake(mac string) error
}

type apiClient struct{}

// Create a new client for api
func NewAPIClient(url string) Client {
	return &apiClient{}
}

// Add a new host, overwrite existing host name if it already exists.
func (a *apiClient) AddHost(host types.Host) error {
	panic("TODO")
}

// Return all hosts
func (a *apiClient) GetHosts() ([]types.Host, error) {
	panic("TODO")
}

// Remove a host, ignore if the host does not exist
func (a *apiClient) RemoveHost(mac string) error {
	panic("TODO")
}

// Return the current status of all hosts
func (a *apiClient) Status() ([]types.HostStatus, error) {
	panic("TODO")
}

// Send a magic packet to wake a host
func (a *apiClient) Wake(mac string) error {
	panic("TODO")
}
