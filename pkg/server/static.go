package server

import (
	"embed"
	"net/http"
	"strings"

	"github.com/heathcliff26/netrouse/pkg/version"
	"github.com/heathcliff26/simple-fileserver/pkg/filesystem"
)

func StaticFileServer(root embed.FS) http.Handler {
	indexlessFS := filesystem.NewIndexlessFilesystem(http.FS(root))
	fs := http.FileServer(indexlessFS)

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("ETag", version.Version())
		res.Header().Set("Cache-Control", "public, max-age=3600")

		if match := req.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, version.Version()) {
				res.WriteHeader(http.StatusNotModified)
				return
			}
		}

		fs.ServeHTTP(res, req)
	})
}
