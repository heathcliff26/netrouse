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
		require := require.New(t)
		assert := assert.New(t)

		persistence.SetConfigFolder(t.TempDir())

		settings := persistence.Settings{
			Remotes: []persistence.RemoteServer{
				{Name: "Alpha", URL: "http://alpha.example"},
				{Name: "Beta", URL: "http://beta.example"},
			},
		}
		err := settings.Save()
		require.NoError(err, "Should save settings")

		app := New()

		require.NotNil(app)
		require.Len(app.tabs, 2)
		require.Len(app.appTabs.Items, 4)
		assert.Equal(app.tabLocal.tab, app.appTabs.Items[0])
		assert.Equal(app.tabs[0].tab, app.appTabs.Items[1])
		assert.Equal(app.tabs[1].tab, app.appTabs.Items[2])
		assert.Equal(app.tabSettings, app.appTabs.Items[3])
		assert.Equal(0, app.appTabs.SelectedIndex())
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

	require.NotNil(app.tabLocal.ctx)
	assert.Nil(app.tabs[0].ctx)

	app.selectTab(app.tabs[0].tab)
	assert.NotNil(app.tabs[0].ctx)

	app.selectTab(app.tabLocal.tab)
	assert.NotNil(app.tabLocal.ctx)
}
