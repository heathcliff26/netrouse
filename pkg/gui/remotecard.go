package gui

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
)

// RemoteCard Extends the base widget to display a remote server and allow deleting it
type RemoteCard struct {
	widget.BaseWidget

	app    *App
	remote persistence.RemoteServer

	card      *widget.Card
	deleteBtn *widget.Button
	editBtn   *widget.Button
}

// Create a new card for displaying a remote server
func NewRemoteCard(app *App) *RemoteCard {
	// Use placeholder data to ensure proper rendering
	c := &RemoteCard{
		app: app,
		remote: persistence.RemoteServer{
			Name: "Name",
			URL:  "URL",
		},
	}
	c.ExtendBaseWidget(c)

	return c
}

// Function to create renderer needed to implement widget
func (c *RemoteCard) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)

	c.deleteBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(), c.delete)
	c.editBtn = widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), c.edit)

	c.card = widget.NewCard(c.remote.Name, c.remote.URL, container.NewHBox(layout.NewSpacer(), c.editBtn, c.deleteBtn))

	return widget.NewSimpleRenderer(c.card)
}

func (c *RemoteCard) SetRemote(remote persistence.RemoteServer) {
	c.remote = remote
	c.card.SetTitle(c.remote.Name)
	c.card.SetSubTitle(c.remote.URL)
	c.Refresh()
}

// Delete the remote server from the list
func (c *RemoteCard) delete() {
	for i := range c.app.settings.Remotes {
		if c.app.settings.Remotes[i] == c.remote {
			c.app.settings.Remotes = deleteItem(c.app.settings.Remotes, i)
			break
		}
	}

	for i := range c.app.tabs {
		if *c.app.tabs[i].remote == c.remote {
			tab := c.app.tabs[i]
			c.app.tabs = deleteItem(c.app.tabs, i)
			c.app.appTabs.Remove(tab.tab)
			break
		}
	}
}

// Edit the remote server
func (c *RemoteCard) edit() {
	remoteName := binding.NewString()
	_ = remoteName.Set(c.remote.Name)
	remoteURL := binding.NewString()
	_ = remoteURL.Set(c.remote.URL)
	dialog.ShowForm(lang.L("Edit Server"), lang.L("Save"), lang.L("Cancel"), []*widget.FormItem{
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

		remote := persistence.RemoteServer{
			Name: name,
			URL:  url,
		}

		for i := range c.app.settings.Remotes {
			if c.app.settings.Remotes[i] == c.remote {
				c.app.settings.Remotes[i] = remote
				break
			}
		}

		for i := range c.app.tabs {
			if *c.app.tabs[i].remote == c.remote {
				c.app.tabs[i].SetRemote(remote)
				c.app.tabs[i].tab.Text = remote.Name
				c.app.appTabs.Refresh()
				break
			}
		}

		c.SetRemote(remote)
	}, c.app.main)
}

// Helper function to delete an item from a slice
func deleteItem[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}
