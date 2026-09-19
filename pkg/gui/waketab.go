package gui

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/heathcliff26/netrouse/pkg/client"
	"github.com/heathcliff26/netrouse/pkg/gui/customthemes"
	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/heathcliff26/netrouse/pkg/utils"
)

var hostStatusIcon = theme.RadioButtonFillIcon()

type wakeTab struct {
	tab    *container.TabItem
	remote *persistence.RemoteServer
	hosts  []*hostWidget
	window fyne.Window
	client client.Client

	lock sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc
}

func newTabFromRemote(window fyne.Window, remote *persistence.RemoteServer) *wakeTab {
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
	// TODO: Better error handling. Likely fetchHosts should never fail
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
		go t.updateStatus()
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
	go t.updateStatus()
}

func (t *wakeTab) fetchHosts() error {
	hosts, err := t.client.GetHosts()
	if err != nil {
		return err
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	t.hosts = make([]*hostWidget, 0, len(hosts))
	for _, host := range hosts {
		t.hosts = append(t.hosts, newHostWidget(t, host))
	}
	t.update()
	return nil
}

func (t *wakeTab) updateStatus() {
	slog.Info("Update status", slog.String("tab", t.tab.Text))
	status, err := t.client.Status()
	if err != nil {
		slog.Error("Failed to update status", slog.String("tab", t.tab.Text), "error", err)
		return
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	for _, s := range status {
		for _, host := range t.hosts {
			if host.host.MAC == s.MAC {
				host.updateStatus(s)
			}
		}
	}
}

func (t *wakeTab) selected() {
	if t.ctx == nil {
		t.ctx, t.cancel = context.WithCancel(context.Background())
	} else {
		select {
		case <-t.ctx.Done():
			t.ctx, t.cancel = context.WithCancel(context.Background())
		default:
			return
		}
	}
	err := t.fetchHosts()
	if err != nil {
		dialog.ShowError(err, t.window)
		return
	}

	go func() {
		slog.Debug("Start periodic status updates", slog.String("tab", t.tab.Text))
		tick := time.NewTicker(30 * time.Second)
		for {
			t.updateStatus()
			select {
			case <-t.ctx.Done():
				slog.Debug("Stop periodic status updates", slog.String("tab", t.tab.Text))
				tick.Stop()
				return
			case <-tick.C:
			}
		}
	}()
}

func (t *wakeTab) unselected() {
	if t.cancel == nil {
		return
	}
	t.cancel()
}

func (t *wakeTab) SetRemote(remote persistence.RemoteServer) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.remote = &remote
	t.client = client.NewAPIClient(remote.URL)
	t.update()
}

type hostWidget struct {
	host   types.Host
	object fyne.CanvasObject
	status *widget.Icon
}

func newHostWidget(parent *wakeTab, host types.Host) *hostWidget {
	status := widget.NewIcon(theme.NewDisabledResource(hostStatusIcon))
	items := make([]fyne.CanvasObject, 0, 2)
	if host.Address != "" {
		statusContainer := container.NewThemeOverride(status, customthemes.NewStatusIconTheme())
		items = append(items, container.NewHBox(statusContainer, widget.NewLabel(host.Address)))
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

func (w *hostWidget) updateStatus(status types.HostStatus) {
	switch {
	case status.Online:
		w.status.SetResource(theme.NewPrimaryThemedResource(hostStatusIcon))
	case status.Error != "":
		slog.Info("Failed to fetch status", slog.String("host", status.Address), slog.String("mac", status.MAC), slog.String("error", status.Error))
		w.status.SetResource(theme.NewDisabledResource(hostStatusIcon))
	default:
		w.status.SetResource(theme.NewErrorThemedResource(hostStatusIcon))
	}
}
