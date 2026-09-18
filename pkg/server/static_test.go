package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/version"
	"github.com/heathcliff26/netrouse/static"
	"github.com/stretchr/testify/assert"
)

func TestStaticFileServer(t *testing.T) {
	t.Run("BasicRequest", func(t *testing.T) {
		assert := assert.New(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/css/bootstrap.css", nil)

		fs := StaticFileServer(static.Assets)

		fs.ServeHTTP(res, req)

		assert.Equal(version.Version(), res.Header().Get("ETag"), "Should have ETag header set")
		assert.NotEmpty(res.Header().Get("Cache-Control"), "Should have Cache-Control header set")
		assert.Equal(http.StatusOK, res.Code, "Should answer with Code 200")
	})
	t.Run("FromCache", func(t *testing.T) {
		assert := assert.New(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/css/bootstrap.css", nil)

		req.Header.Add("If-None-Match", version.Version())

		fs := StaticFileServer(static.Assets)

		fs.ServeHTTP(res, req)

		assert.Equal(version.Version(), res.Header().Get("ETag"), "Should have ETag header set")
		assert.NotEmpty(res.Header().Get("Cache-Control"), "Should have Cache-Control header set")
		assert.Equal(http.StatusNotModified, res.Code, "Should answer with Code 304")
	})
	t.Run("OutdatedCache", func(t *testing.T) {
		assert := assert.New(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/css/bootstrap.css", nil)

		req.Header.Add("If-None-Match", "not-the-version")

		fs := StaticFileServer(static.Assets)

		fs.ServeHTTP(res, req)

		assert.Equal(version.Version(), res.Header().Get("ETag"), "Should have ETag header set")
		assert.NotEmpty(res.Header().Get("Cache-Control"), "Should have Cache-Control header set")
		assert.Equal(http.StatusOK, res.Code, "Should answer with Code 200")
	})
}
