package main

import (
	"embed"
	"log/slog"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

//go:embed translations/*.json
var translationsFS embed.FS

func main() {
	a := app.NewWithID("io.github.heathcliff26.netrouse")
	w := a.NewWindow("NetRouse")

	if err := lang.AddTranslationsFS(translationsFS, "translations"); err != nil {
		slog.Error("Failed to load translations", slog.Any("error", err))
	}

	w.SetContent(widget.NewLabel(lang.L("Hello World!")))
	w.ShowAndRun()
}
