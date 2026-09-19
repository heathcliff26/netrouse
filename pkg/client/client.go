package client

import (
	"bytes"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	v1 "github.com/heathcliff26/netrouse/pkg/server/api/v1"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
)

const (
	apiPath = "/api/v1"
)

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

type apiClient struct {
	http.Client
	endpoint string
}

// Create a new client for api
func NewAPIClient(url string) Client {
	url = strings.TrimSuffix(url, "/")
	if !strings.HasSuffix(url, apiPath) {
		url += apiPath
	}
	return &apiClient{
		endpoint: url,
		Timeout:  time.Second,
	}
}

// Add a new host, overwrite existing host name if it already exists.
func (a *apiClient) AddHost(host types.Host) error {
	body, err := json.Marshal(host)
	if err != nil {
		return fmt.Errorf("failed to parse host: %w", err)
	}
	return a.sendRequest(http.MethodPut, "/hosts", bytes.NewReader(body), nil)
}

// Return all hosts
func (a *apiClient) GetHosts() ([]types.Host, error) {
	var hosts []types.Host
	err := a.sendRequest(http.MethodGet, "/hosts", nil, &hosts)
	return hosts, err
}

// Remove a host, ignore if the host does not exist
func (a *apiClient) RemoveHost(mac string) error {
	return a.sendRequest(http.MethodDelete, fmt.Sprintf("/hosts/%s", mac), nil, nil)
}

// Return the current status of all hosts
func (a *apiClient) Status() ([]types.HostStatus, error) {
	var status []types.HostStatus
	err := a.sendRequest(http.MethodGet, "/hosts/status", nil, &status)
	return status, err
}

// Send a magic packet to wake a host
func (a *apiClient) Wake(mac string) error {
	return a.sendRequest(http.MethodGet, fmt.Sprintf("/wake/%s", mac), nil, nil)
}

func (a *apiClient) sendRequest(method, target string, body io.Reader, v interface{}) error {
	req, err := http.NewRequest(method, a.endpoint+target, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	res, err := a.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		if v == nil {
			return nil
		}
		return json.UnmarshalRead(res.Body, v)
	}

	var msg v1.Response
	err = json.UnmarshalRead(res.Body, &msg)
	if err != nil {
		return fmt.Errorf("request returned with '%s', failed to parse body: %w", res.Status, err)
	}
	return fmt.Errorf("server responded with '%s': %s", res.Status, msg.Reason)
}
