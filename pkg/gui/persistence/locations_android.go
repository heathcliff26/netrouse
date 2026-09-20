//go:build android

package persistence

import (
	"log/slog"

	"fyne.io/fyne/v2"
)

func initConfigFolder() {
	app := fyne.CurrentApp()
	if app == nil {
		slog.Error("No app active, not getting storage location")
		configFolder = "./"
		return
	}
	configFolder = app.Storage().RootURI().Path()
}
