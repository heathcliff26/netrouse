package gui

import (
	"embed"
	"log/slog"

	"fyne.io/fyne/v2"
	fApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/heathcliff26/netrouse/pkg/version"
)

// Used to change the new app function for testing
var newApp = fApp.New

//go:embed translations
var translationsFS embed.FS

type App struct {
	app         fyne.App
	main        fyne.Window
	tabLocal    *wakeTab
	tabSettings *container.TabItem
	tabs        []*wakeTab
}

func New() *App {
	err := lang.AddTranslationsFS(translationsFS, "translations")
	if err != nil {
		slog.Error("Failed to load translations", slog.Any("error", err))
	}

	app := newApp()
	main := app.NewWindow(version.Name)

	a := &App{
		app:  app,
		main: main,
	}
	a.tabLocal = newLocalTab(a.main)
	a.tabSettings = a.newSettingsTab()
	a.main.SetContent(a.newTabs())
	a.main.Resize(fyne.NewSize(450, 600))
	a.main.Show()

	return a
}

func (a *App) Run() {
	a.app.Run()
}

func (a *App) newTabs() *container.AppTabs {
	items := make([]*container.TabItem, 0, len(a.tabs)+2)
	items = append(items, a.tabLocal.tab)
	for _, tab := range a.tabs {
		items = append(items, tab.tab)
	}
	items = append(items, a.tabSettings)

	tabs := container.NewAppTabs(items...)
	tabs.SetTabLocation(container.TabLocationLeading)
	return tabs
}

func (a *App) newSettingsTab() *container.TabItem {
	return container.NewTabItemWithIcon(lang.L("Settings"), theme.SettingsIcon(), widget.NewLabel("TODO"))
}
