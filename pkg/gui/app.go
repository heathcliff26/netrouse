package gui

import (
	"embed"
	"log/slog"
	"runtime"

	"fyne.io/fyne/v2"
	fApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
	"github.com/heathcliff26/netrouse/pkg/version"
)

// Used to change the new app function for testing
var newApp = fApp.New

//go:embed translations
var translationsFS embed.FS

type App struct {
	app  fyne.App
	main fyne.Window

	tabLocal    *wakeTab
	tabSettings *container.TabItem
	tabs        []*wakeTab
	appTabs     *container.AppTabs

	settings *persistence.Settings
}

func New() *App {
	err := lang.AddTranslationsFS(translationsFS, "translations")
	if err != nil {
		slog.Error("Failed to load translations", slog.Any("error", err))
	}

	app := newApp()
	main := app.NewWindow(version.Name)

	persistence.Init()
	s, err := persistence.LoadSettings()
	if err != nil {
		slog.Error("Failed to load settings", slog.Any("error", err))
	}

	a := &App{
		app:      app,
		main:     main,
		tabs:     make([]*wakeTab, 0, len(s.Remotes)),
		settings: s,
	}
	a.tabLocal = newLocalTab(a.main)
	a.tabSettings = a.newSettingsTab()
	a.appTabs = a.newTabs()
	for _, remote := range a.settings.Remotes {
		a.addRemote(&remote)
	}
	a.main.SetContent(a.appTabs)
	a.main.Resize(a.settings.WindowSize.ToFyne())
	a.main.SetFullScreen(a.settings.FullScreen)
	a.main.Show()

	a.app.Lifecycle().SetOnStarted(a.onStarted)
	a.app.Lifecycle().SetOnStopped(a.onStopped)

	return a
}

func (a *App) Run() {
	a.app.Run()
}

func (a *App) newTabs() *container.AppTabs {
	items := make([]*container.TabItem, 0, len(a.tabs)+2)
	items = append(items, a.tabLocal.tab)
	items = append(items, a.tabSettings)

	tabs := container.NewAppTabs(items...)
	tabs.SetTabLocation(container.TabLocationLeading)

	tabs.OnSelected = a.selectTab
	return tabs
}

func (a *App) addRemote(remote *persistence.RemoteServer) {
	tab := newTabFromRemote(a.main, remote)
	a.tabs = append(a.tabs, tab)

	index := a.appTabs.SelectedIndex()
	a.appTabs.Items = append(a.appTabs.Items[:len(a.appTabs.Items)-1], tab.tab, a.tabSettings)

	offset := 0
	if index == len(a.appTabs.Items)-2 {
		offset = 1
	}
	a.appTabs.SelectIndex(index + offset)
	a.appTabs.Refresh()
}

func (a *App) selectTab(item *container.TabItem) {
	slog.Debug("Selected new tab", slog.String("tab", item.Text))
	if a.tabLocal.tab == item {
		a.tabLocal.selected()
	} else {
		a.tabLocal.unselected()
	}
	for _, t := range a.tabs {
		if t.tab == item {
			t.selected()
		} else {
			t.unselected()
		}
	}
}

func (a *App) newSettingsTab() *container.TabItem {
	title := widget.NewLabel(lang.L("Settings"))
	title.TextStyle = fyne.TextStyle{Bold: true}

	tab := container.NewTabItemWithIcon(lang.L("Settings"), theme.SettingsIcon(), nil)

	remoteList := widget.NewList(
		func() int {
			return len(a.settings.Remotes)
		},
		func() fyne.CanvasObject {
			return NewRemoteCard(a)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*RemoteCard).SetRemote(a.settings.Remotes[i])
		},
	)
	remoteName := binding.NewString()
	remoteURL := binding.NewString()
	addRemoteButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		dialog.ShowForm(lang.L("Add Server"), lang.L("Add"), lang.L("Cancel"), []*widget.FormItem{
			widget.NewFormItem(lang.L("Name"), widget.NewEntryWithData(remoteName)),
			widget.NewFormItem(lang.L("URL"), widget.NewEntryWithData(remoteURL)),
		}, func(b bool) {
			if !b {
				return
			}
			name, err := remoteName.Get()
			if err != nil {
				slog.Error("Failed to get remote Name from binding", "error", err)
				return
			}
			url, err := remoteURL.Get()
			if err != nil {
				slog.Error("Failed to get remote URL from binding", "error", err)
				return
			}
			url = urlAddSchema(url)
			remote := persistence.RemoteServer{
				Name: name,
				URL:  url,
			}
			a.settings.Remotes = append(a.settings.Remotes, remote)
			a.addRemote(&remote)
			tab.Content.Resize(tab.Content.MinSize())

			_ = remoteName.Set("")
			_ = remoteURL.Set("")
		}, a.main)
	})
	remoteContainer := widget.NewCard(lang.L("Server"), "", container.NewBorder(nil, addRemoteButton, nil, nil, remoteList))

	resetWindowBtn := widget.NewButton(lang.L("Reset Window"), a.resetWindow)
	resetWindowBtn.Hidden = fyne.CurrentDevice().IsMobile()

	aboutBtn := widget.NewButton(lang.L("About"), a.showAbout)

	tab.Content = container.NewBorder(title, container.NewVBox(resetWindowBtn, aboutBtn), nil, nil, remoteContainer)

	return tab
}

func (a *App) resetWindow() {
	a.settings.WindowSize = persistence.DefaultSettings().WindowSize
	a.main.SetFullScreen(false)
	a.main.Resize(a.settings.WindowSize.ToFyne())
}

func (a *App) onStarted() {
	currentTab := a.settings.SelectedTab
	if currentTab >= len(a.appTabs.Items) {
		currentTab = 0
	}
	a.appTabs.SelectIndex(currentTab)
	a.selectTab(a.appTabs.Selected())
}

func (a *App) onStopped() {
	a.settings.SelectedTab = a.appTabs.SelectedIndex()
	a.settings.WindowSize = persistence.SizeFromFyne(a.main.Canvas().Size())
	a.settings.FullScreen = a.main.FullScreen()
	err := a.settings.Save()
	if err != nil {
		slog.Error("Failed to save settings", slog.Any("error", err))
	}
}

func (a *App) showAbout() {
	data := [][]string{
		{"Version:", a.app.Metadata().Version},
		{"Commit:", version.Commit()},
		{"Go:", runtime.Version()},
	}

	description := container.NewVBox()
	values := container.NewVBox()
	for _, row := range data {
		description.Add(widget.NewLabel(row[0]))
		values.Add(widget.NewLabel(row[1]))
	}

	versionTable := container.NewHBox(description, values)

	dialog.ShowCustom(lang.L("About"), lang.L("OK"), versionTable, a.main)
}
