package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/server/config"
	"github.com/heathcliff26/netrouse/pkg/server/storage"
	"github.com/heathcliff26/netrouse/static"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cfgServer := config.ServerConfig{
		Port: 8080,
		SSL: config.SSLConfig{
			Enabled: true,
			Cert:    "test.crt",
			Key:     "test.key",
		},
	}
	cfgStorage := storage.NewDefaultStorageConfig()
	cfgStorage.File.Path = "testdata/hosts.yaml"

	s, err := NewServer(cfgServer, cfgStorage)

	require.NoError(err, "Should create new server")
	require.NotNil(s, "Server should not be empty")

	assert.Equal(":8080", s.addr, "Server should have address set")
	assert.Equal(cfgServer.SSL, s.ssl, "Server should have SSL Config")

	indexHTML, indexChecksum, err := s.storage.GetIndexHTML()
	require.NoError(err, "Should create index.html")

	assert.Contains(indexHTML, "testName", "Should add hostname to index.html")
	assert.Contains(indexHTML, "TESTMAC", "Should add MAC to index.html")
	assert.Contains(indexHTML, "onclick=\"wake('TESTMAC', 'testName');\"", "Should add onClick javascript function call with MAC")
	assert.NotEmpty(indexChecksum, "Server should have checksum of index.html")
}

func TestServer(t *testing.T) {
	t.Run("SSL", func(t *testing.T) {
		assert := assert.New(t)

		cfgServer := config.ServerConfig{
			Port: 8080,
			SSL: config.SSLConfig{
				Enabled: true,
				Cert:    "test.crt",
				Key:     "test.key",
			},
		}
		cfgStorage := storage.NewDefaultStorageConfig()
		cfgStorage.File.Path = "testdata/hosts.yaml"

		s, err := NewServer(cfgServer, cfgStorage)
		require.NoError(t, err, "Should create server without error")

		assert.Error(s.Run(), "Server should fail to run, as the ssl certificate and key do not exist")
	})
	cfgServer := config.ServerConfig{
		Port: 8080,
	}
	cfgStorage := storage.NewDefaultStorageConfig()
	cfgStorage.File.Path = "testdata/hosts.yaml"
	s, err := NewServer(cfgServer, cfgStorage)
	require.NoError(t, err, "Should create server without error")

	go func() {
		err := s.Run()
		if err != nil {
			t.Logf("Failed to run server: %v", err)
			t.Fail()
		}
	}()
	address := "http://localhost:8080"

	t.Run("IndexHandler", func(t *testing.T) {
		assert := assert.New(t)

		indexHTML, indexChecksum, err := s.storage.GetIndexHTML()
		require.NoError(t, err, "Should get index.html")

		for _, path := range []string{"/", "/index.html"} {
			res, err := http.Get(address + path)
			t.Cleanup(func() {
				res.Body.Close()
			})

			assert.NoErrorf(err, "Should not return error for request to %s", path)
			assert.Equal(http.StatusOK, res.StatusCode, "Should return 200 for request to %s", path)
			assert.Equalf(indexChecksum, res.Header.Get("ETag"), "Should have ETag set on request to %s", path)
			assert.NotEmptyf(res.Header.Get("Cache-Control"), "Should have Cache-Control header set on request to %s", path)

			body, err := io.ReadAll(res.Body)
			assert.NoErrorf(err, "Should read body without error for request to %s", path)
			assert.Equalf(indexHTML, string(body), "Body should match index.html on path %s", path)
		}

		resNotFound, err := http.Get(address + "/something")
		t.Cleanup(func() {
			resNotFound.Body.Close()
		})
		assert.NoError(err, "Should receive valid response for random path")
		assert.Equal(http.StatusNotFound, resNotFound.StatusCode, "Should return 404 for random path")

		req, _ := http.NewRequest(http.MethodGet, address+"/", nil)
		req.Header.Add("If-None-Match", indexChecksum)

		resCache, err := (&http.Client{}).Do(req)
		t.Cleanup(func() {
			resCache.Body.Close()
		})
		assert.NoError(err, "Should receive valid response to call with cache")
		assert.Equal(http.StatusNotModified, resCache.StatusCode, "Should receive cache hit")
	})
	t.Run("API", func(t *testing.T) {
		assert := assert.New(t)

		res, err := http.Get(address + "/api/v1/wake/not-a-mac")
		t.Cleanup(func() {
			res.Body.Close()
		})

		assert.NoError(err)
		assert.Equal(http.StatusBadRequest, res.StatusCode, "Should receive a bad request response when using a malformed mac address")
	})

	assetTMatrix := map[string]string{"CSS": "css/bootstrap.css", "Icons": "icons/favicon.svg", "JS": "js/client.js"}
	for name, path := range assetTMatrix {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			file, err := static.Assets.ReadFile(path)
			assert.NoError(err, "Should read file from static")

			res, err := http.Get(address + "/" + path)
			t.Cleanup(func() {
				res.Body.Close()
			})

			assert.NoError(err, "Request should not return an error")
			assert.Equal(http.StatusOK, res.StatusCode, "Request should not return an error")
			body, _ := io.ReadAll(res.Body)
			assert.Equal(string(file), string(body), "Response should match file")
		})
	}
}
