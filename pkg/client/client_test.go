package client

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIClient(t *testing.T) {
	assert := assert.New(t)

	client := NewAPIClient("http://localhost:8080")
	c := client.(*apiClient)

	assert.Equal("http://localhost:8080/api/v1", c.endpoint, "Should add api path")
	assert.Equal(time.Second, c.Timeout, "Should set timeout")

	client = NewAPIClient("http://localhost:8080/api/v1")
	c = client.(*apiClient)
	assert.Equal("http://localhost:8080/api/v1", c.endpoint, "Should only add api path when needed")

	client = NewAPIClient("http://localhost:8080/api/v1/")
	c = client.(*apiClient)
	assert.Equal("http://localhost:8080/api/v1", c.endpoint, "Should trim trailing slash")
}

func TestAddHost(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	host := types.Host{
		Name:    "test-host",
		MAC:     "AA:BB:CC:DD:EE:FF",
		Address: "host.example.org",
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		assert.Equal(http.MethodPut, request.Method, "Should use PUT")
		assert.Equal("/hosts", request.URL.Path, "Should use hosts endpoint")

		var body types.Host
		err := json.UnmarshalRead(request.Body, &body)
		assert.NoError(err, "Should parse request body")
		assert.Equal(host, body, "Should send host")
		response.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := &apiClient{endpoint: server.URL}
	require.NoError(client.AddHost(t.Context(), host), "Should add host")
}

func TestGetHosts(t *testing.T) {
	require := require.New(t)
	hosts := []types.Host{
		{
			Name:    "test-host",
			MAC:     "AA:BB:CC:DD:EE:FF",
			Address: "host.example.org",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		require.Equal(http.MethodGet, request.Method, "Should use GET")
		require.Equal("/hosts", request.URL.Path, "Should use hosts endpoint")
		response.Header().Set("Content-Type", "application/json")
		err := json.MarshalWrite(response, hosts)
		require.NoError(err, "Should send response")
	}))
	t.Cleanup(server.Close)

	client := &apiClient{endpoint: server.URL}
	result, err := client.GetHosts(t.Context())
	require.NoError(err, "Should get hosts")
	require.Equal(hosts, result)
}

func TestRemoveHost(t *testing.T) {
	require := require.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		require.Equal(http.MethodDelete, request.Method, "Should use DELETE")
		require.Equal("/hosts/AA:BB:CC:DD:EE:FF", request.URL.Path, "Should include MAC address")
		response.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := &apiClient{endpoint: server.URL}
	require.NoError(client.RemoveHost(t.Context(), "AA:BB:CC:DD:EE:FF"), "Should remove host")
}

func TestStatus(t *testing.T) {
	require := require.New(t)
	status := []types.HostStatus{
		{
			MAC:     "AA:BB:CC:DD:EE:FF",
			Address: "host.example.org",
			Online:  true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		require.Equal(http.MethodGet, request.Method, "Should use GET")
		require.Equal("/hosts/status", request.URL.Path, "Should use status endpoint")
		err := json.MarshalWrite(response, status)
		require.NoError(err, "Should send response")
	}))
	t.Cleanup(server.Close)

	client := &apiClient{endpoint: server.URL}
	result, err := client.Status(t.Context())
	require.NoError(err, "Should get host status")
	require.Equal(status, result)
}

func TestWake(t *testing.T) {
	require := require.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		require.Equal(http.MethodGet, request.Method, "Should use GET")
		require.Equal("/wake/AA:BB:CC:DD:EE:FF", request.URL.Path, "Should include MAC address")
		response.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := &apiClient{endpoint: server.URL}
	require.NoError(client.Wake(t.Context(), "AA:BB:CC:DD:EE:FF"), "Should wake host")
}

func TestSendRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		require := require.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
			response.WriteHeader(http.StatusOK)
			_, err := response.Write([]byte(`{"message":"awake"}`))
			require.NoError(err, "Should send response")
		}))
		t.Cleanup(server.Close)

		client := &apiClient{endpoint: server.URL}
		var response struct {
			Message string `json:"message"`
		}

		err := client.sendRequest(t.Context(), http.MethodGet, "/hosts", nil, &response)
		require.NoError(err, "Should send request")
		require.Equal("awake", response.Message, "Should return message")
	})
	t.Run("ErrorCreateRequest", func(t *testing.T) {
		assert := assert.New(t)
		client := &apiClient{endpoint: "foo"}

		err := client.sendRequest(t.Context(), http.MethodGet, "\n", nil, nil)
		assert.ErrorContains(err, "failed to create request")
	})
	t.Run("ErrorTransport", func(t *testing.T) {
		assert := assert.New(t)
		client := &apiClient{endpoint: "https://localhost:6666"}

		err := client.sendRequest(t.Context(), http.MethodGet, "/hosts", nil, nil)
		assert.ErrorContains(err, "connection refused")
	})
	t.Run("ErrorParseBody", func(t *testing.T) {
		require := require.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
			response.WriteHeader(http.StatusBadRequest)
			_, err := response.Write([]byte("not json"))
			require.NoError(err, "Should send response")
		}))
		t.Cleanup(server.Close)

		client := &apiClient{endpoint: server.URL}

		err := client.sendRequest(t.Context(), http.MethodGet, "/hosts", nil, nil)
		require.ErrorContains(err, "request returned with '400 Bad Request', failed to parse body")
	})
	t.Run("ErrorResponse", func(t *testing.T) {
		require := require.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
			response.WriteHeader(http.StatusBadRequest)
			_, err := response.Write([]byte(`{"reason":"invalid host"}`))
			require.NoError(err, "Should send response")
		}))
		t.Cleanup(server.Close)

		client := &apiClient{endpoint: server.URL}

		err := client.sendRequest(t.Context(), http.MethodGet, "/hosts", nil, nil)
		require.ErrorContains(err, "server responded with '400 Bad Request': invalid host")
	})
}
