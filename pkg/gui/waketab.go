package gui

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

var (
	hostStatusIcon = theme.RadioButtonFillIcon()

	hostStatusOnline  = theme.NewPrimaryThemedResource(hostStatusIcon)
	hostStatusOffline = theme.NewErrorThemedResource(hostStatusIcon)
	hostStatusUnknown = theme.NewDisabledResource(hostStatusIcon)
)

type wakeTab struct {
	tab    *container.TabItem
	remote *persistence.RemoteServer
	hosts  []*hostWidget
	window fyne.Window
	client client.Client

	lock sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	title               *widget.Label
	errFetch, errStatus *canvas.Text
	hostsContainer      *fyne.Container
}

func newTabFromRemote(window fyne.Window, remote *persistence.RemoteServer) *wakeTab {
	tab := &wakeTab{
		remote: remote,
		window: window,
		client: client.NewAPIClient(remote.URL),
	}
	tab.tab = container.NewTabItemWithIcon(remote.Name, theme.ComputerIcon(), nil)
	tab.init()
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
		dialog.ShowInformation(lang.L("Error"), lang.L("error.createClient"), window)
	}
	tab.client = client
	tab.init()
	return tab
}

func (t *wakeTab) init() {
	var titleStr string
	if t.remote == nil {
		titleStr = lang.L("Local Devices")
	} else {
		titleStr = remoteTitle(t.remote.Name)
	}
	t.title = widget.NewLabel(titleStr)
	t.title.TextStyle = fyne.TextStyle{Bold: true}

	addButton := widget.NewButtonWithIcon(lang.L("Add"), theme.ContentAddIcon(), t.addHost)
	refreshButton := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), t.fetchHosts)
	refreshButton.Hidden = t.remote == nil

	t.errFetch = newErrorText(lang.L("error.fetchHosts"))
	t.errStatus = newErrorText(lang.L("error.getStatus"))
	t.hostsContainer = container.NewVBox()

	t.tab.Content = container.NewBorder(
		container.NewHBox(layout.NewSpacer(), t.title, layout.NewSpacer()),
		container.NewBorder(nil, nil, nil, refreshButton, addButton),
		nil,
		nil,
		container.NewVBox(t.hostsContainer, t.errFetch, t.errStatus),
	)

	// Ensure context is never nil, but start with cancelled ctx
	t.ctx, t.cancel = context.WithCancel(context.Background())
	t.cancel()
}

func (t *wakeTab) update() {
	hosts := make([]fyne.CanvasObject, 0, len(t.hosts))
	for _, host := range t.hosts {
		hosts = append(hosts, host.card)
	}
	t.hostsContainer.Objects = hosts
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

		err = t.client.AddHost(t.ctx, host)
		if errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			slog.Error("Failed to add host", slog.String("client", t.tab.Text), "error", err)
			dialog.ShowInformation(lang.L("Error"), lang.L("error.addHost"), t.window)
			return
		}
		t.fetchHosts()
	}, t.window)
	d.Show()
}

func (t *wakeTab) removeHost(mac string) {
	err := t.client.RemoveHost(t.ctx, mac)
	if errors.Is(err, context.Canceled) {
		return
	}
	if err != nil {
		slog.Error("Failed to remove host", slog.String("client", t.tab.Text), "error", err)
		dialog.ShowInformation(lang.L("Error"), lang.L("error.removeHost"), t.window)
		return
	}
	t.fetchHosts()
}

func (t *wakeTab) fetchHosts() {
	hosts, err := t.client.GetHosts(t.ctx)

	t.lock.Lock()
	defer t.lock.Unlock()
	defer t.update()

	if errors.Is(err, context.Canceled) {
		return
	}
	if err != nil {
		slog.Error("Failed to fetch hosts", slog.String("client", t.tab.Text), "error", err)
		t.errFetch.Show()
		return
	}
	t.errFetch.Hide()

	newHosts := make([]*hostWidget, 0, len(hosts))
	for _, host := range hosts {
		i := slices.IndexFunc(t.hosts, func(w *hostWidget) bool {
			return w.host.MAC == host.MAC
		})
		if i == -1 {
			newHosts = append(newHosts, newHostWidget(t, host))
		} else {
			newHosts = append(newHosts, t.hosts[i])
			t.hosts[i].updateHost(host)
		}
	}
	t.hosts = newHosts

	go t.updateStatus()
}

func (t *wakeTab) updateStatus() {
	slog.Info("Update status", slog.String("tab", t.tab.Text))
	status, err := t.client.Status(t.ctx)

	if errors.Is(err, context.Canceled) {
		return
	}

	fyne.DoAndWait(func() {
		t.lock.Lock()
		defer t.lock.Unlock()

		if err != nil {
			slog.Error("Failed to update status", slog.String("tab", t.tab.Text), "error", err)
			t.errStatus.Show()
			return
		}
		t.errStatus.Hide()

		for _, s := range status {
			for _, host := range t.hosts {
				if host.host.MAC == s.MAC {
					host.updateStatus(s)
				}
			}
		}
	})
}

func (t *wakeTab) selected() {
	if t.ctx.Err() == nil {
		return
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	t.fetchHosts()

	go func() {
		ctx := t.ctx
		slog.Debug("Start periodic status updates", slog.String("tab", t.tab.Text))
		tick := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ctx.Done():
				slog.Debug("Stop periodic status updates", slog.String("tab", t.tab.Text))
				tick.Stop()
				return
			case <-tick.C:
			}
			t.updateStatus()
		}
	}()
}

func (t *wakeTab) unselected() {
	t.cancel()
}

func (t *wakeTab) SetRemote(remote persistence.RemoteServer) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.remote = &remote
	t.client = client.NewAPIClient(remote.URL)
	t.title.SetText(remoteTitle(remote.Name))
}

type hostWidget struct {
	host types.Host

	card    *widget.Card
	status  *widget.Icon
	address *widget.Label

	deleteBtn *widget.Button
	wakeBtn   *widget.Button

	statusContainer *container.ThemeOverride
}

func newHostWidget(parent *wakeTab, host types.Host) *hostWidget {
	items := make([]fyne.CanvasObject, 0, 2)

	status := widget.NewIcon(hostStatusUnknown)
	address := widget.NewLabel(host.Address)
	statusContainer := container.NewThemeOverride(status, customthemes.NewStatusIconTheme())
	statusContainer.Hidden = host.Address == ""
	items = append(items, container.NewHBox(statusContainer, address))

	delete := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		parent.removeHost(host.MAC)
	})
	wake := widget.NewButton(lang.L("Wake"), nil)
	wake.OnTapped = func() {
		wake.Disable()
		defer wake.Enable()

		err := parent.client.Wake(parent.ctx, host.MAC)
		if errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			slog.Error("Failed to wake host", slog.String("mac", host.MAC), slog.String("error", err.Error()))
			dialog.ShowInformation(lang.L("Error"), lang.L("error.wake"), parent.window)
		}
	}
	items = append(items, container.NewBorder(nil, nil, nil, delete, wake))

	content := container.NewVBox(items...)
	card := widget.NewCard(host.Name, host.MAC, content)

	return &hostWidget{
		host:    host,
		card:    card,
		status:  status,
		address: address,

		deleteBtn: delete,
		wakeBtn:   wake,

		statusContainer: statusContainer,
	}
}

func (w *hostWidget) updateHost(new types.Host) {
	if w.host.Name != new.Name {
		w.card.SetTitle(new.Name)
		w.host.Name = new.Name
	}

	if w.host.Address != new.Address {
		w.address.SetText(new.Address)
		w.statusContainer.Hidden = new.Address == ""
		w.status.SetResource(hostStatusUnknown)
		w.host.Address = new.Address
	}
}

func (w *hostWidget) updateStatus(status types.HostStatus) {
	switch {
	case status.Online:
		w.status.SetResource(hostStatusOnline)
	case status.Error != "":
		slog.Info("Failed to fetch status", slog.String("host", status.Address), slog.String("mac", status.MAC), slog.String("error", status.Error))
		w.status.SetResource(hostStatusUnknown)
	default:
		w.status.SetResource(hostStatusOffline)
	}
}

func newErrorText(msg string) *canvas.Text {
	text := canvas.NewText(lang.L("Error")+": "+msg, customthemes.Red())
	text.TextStyle.Bold = true
	text.Alignment = fyne.TextAlignCenter
	text.Hidden = true
	return text
}

func remoteTitle(name string) string {
	return lang.L("Devices") + ": " + name
}
