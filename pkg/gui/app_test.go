package gui

import (
	"testing"

	fApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	oldFolder := persistence.ConfigFolder()
	newApp = test.NewApp
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
		newApp = fApp.New
	})

	t.Run("DefaultApp", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		persistence.SetConfigFolder(t.TempDir())

		app := New()

		require.NotNil(app)
		assert.NotNil(app.app)
		assert.NotNil(app.main)
		assert.NotNil(app.tabLocal)
		assert.NotNil(app.tabSettings)
		assert.NotNil(app.tabs)
		assert.NotNil(app.appTabs)
		assert.Len(app.appTabs.Items, 2)
		assert.Equal(app.tabLocal.tab, app.appTabs.Items[0])
		assert.Equal(app.tabSettings, app.appTabs.Items[1])
		assert.Equal(0, app.appTabs.SelectedIndex())
	})
	t.Run("LoadsConfiguredRemotes", func(t *testing.T) {
		assert := assert.New(t)

		persistence.SetConfigFolder(t.TempDir())

		app := New()
		app.settings.SelectedTab = 1

		app.onStarted()
		assert.Equal(1, app.appTabs.SelectedIndex(), "Should select the correct tab")
	})
	t.Run("SelectTabOnStartup", func(t *testing.T) {
		require := require.New(t)

		persistence.SetConfigFolder(t.TempDir())

		settings := persistence.Settings{
			Remotes: []persistence.RemoteServer{
				{Name: "Alpha", URL: "http://alpha.example"},
				{Name: "Beta", URL: "http://beta.example"},
			},
		}
		err := settings.Save()
		require.NoError(err, "Should save settings")
	})
	t.Run("SaveSettingsOnClose", func(t *testing.T) {
		assert := assert.New(t)
		persistence.SetConfigFolder(t.TempDir())
		a := New()

		a.settings.Remotes = []persistence.RemoteServer{
			{Name: "Alpha", URL: "http://alpha.example"},
			{Name: "Beta", URL: "http://beta.example"},
		}

		a.onStopped()

		settings, err := persistence.LoadSettings()
		assert.NoError(err, "Should load settings")
		assert.Equal(a.settings, settings, "Should match saved settings")
	})
}

func TestAddRemoteAndSelectTab(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	oldFolder := persistence.ConfigFolder()
	newApp = test.NewApp
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
		newApp = fApp.New
	})

	persistence.SetConfigFolder(t.TempDir())
	app := New()

	app.selectTab(app.tabSettings)

	remote := &persistence.RemoteServer{Name: "Remote One", URL: "http://remote.example"}
	app.addRemote(remote)

	t.Cleanup(app.tabLocal.unselected)
	t.Cleanup(app.tabs[0].unselected)

	require.Len(app.tabs, 1)
	require.Len(app.appTabs.Items, 3)
	assert.Equal(app.tabLocal.tab, app.appTabs.Items[0])
	assert.Equal(app.tabs[0].tab, app.appTabs.Items[1])
	assert.Equal(app.tabSettings, app.appTabs.Items[2])

	assert.Error(app.tabLocal.ctx.Err(), "Should have no context selected")
	assert.Error(app.tabs[0].ctx.Err(), "Should have no context selected")

	app.selectTab(app.tabs[0].tab)
	assert.NoError(app.tabs[0].ctx.Err(), "Should have context selected")
	assert.Error(app.tabLocal.ctx.Err(), "Should have no context selected")

	app.selectTab(app.tabLocal.tab)
	assert.NotNil(app.tabLocal.ctx)
	assert.NoError(app.tabLocal.ctx.Err(), "Should have context selected")
	assert.Error(app.tabs[0].ctx.Err(), "Should have no context selected")
}

func TestResetWindow(t *testing.T) {
	assert := assert.New(t)

	oldFolder := persistence.ConfigFolder()
	newApp = test.NewApp
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
		newApp = fApp.New
	})

	persistence.SetConfigFolder(t.TempDir())
	app := New()

	app.settings.WindowSize = persistence.Size{Width: 100, Height: 100}
	app.main.Resize(app.settings.WindowSize.ToFyne())
	app.main.SetFullScreen(true)

	app.resetWindow()

	assert.Equal(persistence.DefaultSettings().WindowSize, app.settings.WindowSize, "Should reset size in settings")
	assert.Equal(persistence.DefaultSettings().WindowSize.ToFyne(), app.main.Canvas().Size(), "Should reset window size")
	assert.False(app.main.FullScreen(), "Should not be fullscreen")
}
