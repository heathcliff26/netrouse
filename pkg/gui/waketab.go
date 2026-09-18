package gui

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/heathcliff26/netrouse/pkg/client"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/heathcliff26/netrouse/pkg/utils"
)

const (
	hostStatusUnknown = iota
	hostStatusOnline
	hostStatusOffline
)

const (
	hostStatusUnknownStr = "⚪"
	hostStatusOnlineStr  = "🟢"
	hostStatusOfflineStr = "🔴"
)

type wakeTab struct {
	tab    *container.TabItem
	remote *RemoteServer
	hosts  []*hostWidget
	window fyne.Window
	client client.Client
}

func newTabFromRemote(window fyne.Window, remote *RemoteServer) *wakeTab {
	tab := &wakeTab{
		remote: remote,
		window: window,
		client: client.NewAPIClient(remote.URL),
	}
	tab.tab = container.NewTabItemWithIcon(remote.Name, theme.ComputerIcon(), nil)
	tab.update()
	return tab
}

func newLocalTab(window fyne.Window) *wakeTab {
	tab := &wakeTab{
		tab:    container.NewTabItemWithIcon(lang.L("Local"), theme.HomeIcon(), nil),
		window: window,
	}
	client, err := newLocalClient()
	if err != nil {
		slog.Error("failed to create local client", "error", err)
		dialog.ShowError(err, window)
	}
	tab.client = client
	err = tab.fetchHosts()
	if err != nil {
		slog.Error("failed to fetch hosts", "error", err)
		dialog.ShowError(err, window)
	}
	return tab
}

func (t *wakeTab) update() {
	var titleStr string
	if t.remote == nil {
		titleStr = lang.L("Local Devices")
	} else {
		titleStr = lang.L("Devices") + ": " + t.remote.Name
	}
	title := widget.NewLabel(titleStr)
	title.TextStyle = fyne.TextStyle{Bold: true}
	addButton := widget.NewButtonWithIcon(lang.L("Add"), theme.ContentAddIcon(), t.addHost)

	hosts := make([]fyne.CanvasObject, 0, len(t.hosts))
	for _, host := range t.hosts {
		hosts = append(hosts, host.object)
	}

	t.tab.Content = container.NewBorder(
		container.NewHBox(layout.NewSpacer(), title, layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(), addButton, layout.NewSpacer()),
		nil,
		nil,
		container.NewVBox(hosts...),
	)
}

func (t *wakeTab) addHost() {
	name := binding.NewString()
	nameEntry := widget.NewEntryWithData(name)
	nameEntry.Validator = func(name string) error {
		if !utils.ValidateHostname(name) {
			return fmt.Errorf("invalid name")
		}
		return nil
	}
	nameItem := widget.NewFormItem(lang.L("Name"), nameEntry)

	mac := binding.NewString()
	macEntry := widget.NewEntryWithData(mac)
	macEntry.Validator = func(mac string) error {
		if !utils.ValidateMACAddress(mac) {
			return fmt.Errorf("invalid MAC address")
		}
		return nil
	}
	macItem := widget.NewFormItem(lang.L("MAC"), macEntry)

	addr := binding.NewString()
	addrEntry := widget.NewEntryWithData(addr)
	addrItem := widget.NewFormItem(lang.L("Address"), addrEntry)

	d := dialog.NewForm(lang.L("Add New Host"), lang.L("Add"), lang.L("Cancel"), []*widget.FormItem{nameItem, macItem, addrItem}, func(b bool) {
		if !b {
			return
		}
		var err error
		var host types.Host

		host.Name, err = name.Get()
		if err != nil {
			slog.Error("Failed to get host Name from binding", "error", err)
			return
		}
		host.MAC, err = mac.Get()
		if err != nil {
			slog.Error("Failed to get host MAC from binding", "error", err)
			return
		}
		host.Address, err = addr.Get()
		if err != nil {
			slog.Error("Failed to get host Address from binding", "error", err)
			return
		}

		err = t.client.AddHost(host)
		if err != nil {
			dialog.ShowError(err, t.window)
			return
		}
		err = t.fetchHosts()
		if err != nil {
			dialog.ShowError(err, t.window)
			return
		}
	}, t.window)
	d.Show()
}

func (t *wakeTab) removeHost(mac string) {
	err := t.client.RemoveHost(mac)
	if err != nil {
		dialog.ShowError(err, t.window)
		return
	}
	err = t.fetchHosts()
	if err != nil {
		dialog.ShowError(err, t.window)
		return
	}
}

func (t *wakeTab) fetchHosts() error {
	hosts, err := t.client.GetHosts()
	if err != nil {
		return err
	}

	t.hosts = make([]*hostWidget, 0, len(hosts))
	for _, host := range hosts {
		t.hosts = append(t.hosts, newHostWidget(t, host))
	}
	t.update()
	return nil
}

type hostWidget struct {
	host   types.Host
	object fyne.CanvasObject
	status *widget.Label
}

func newHostWidget(parent *wakeTab, host types.Host) *hostWidget {
	status := widget.NewLabel(hostStatusUnknownStr)
	items := make([]fyne.CanvasObject, 0, 2)
	if host.Address != "" {
		items = append(items, container.NewHBox(status, widget.NewLabel(host.Address)))
	}
	delete := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		parent.removeHost(host.MAC)
	})
	wake := widget.NewButton(lang.L("Wake"), nil)
	wake.OnTapped = func() {
		wake.Disable()
		defer wake.Enable()

		err := parent.client.Wake(host.MAC)
		if err != nil {
			dialog.ShowError(err, parent.window)
		}
	}
	// TODO: See if we can stretch wake to fill the whole row
	items = append(items, container.NewHBox(layout.NewSpacer(), wake, delete, layout.NewSpacer()))
	content := container.NewVBox(items...)
	card := widget.NewCard(host.Name, host.MAC, content)
	return &hostWidget{
		host:   host,
		object: card,
		status: status,
	}
}

type RemoteServer struct {
	Name string
	URL  string
}
