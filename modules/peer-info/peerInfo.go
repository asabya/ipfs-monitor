package peer_info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sab94/ipfs-monitor/types"

	"github.com/Sab94/ipfs-monitor/client"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/Sab94/ipfs-monitor/widget"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type PeerInfo widget.Widget

const (
	WidgetName = "peerinfo"
	URL        = "id"
)

func NewWidget(cfg *config.Config, httpClient *client.HttpClient,
	app *tview.Application) block.Block {
	if !cfg.Monitor.Widgets[WidgetName].Enabled {
		return nil
	}

	w := widget.NewWidget(cfg, httpClient, app, WidgetName, URL)
	piWidget := PeerInfo(w)
	piWidget.Render()

	return &piWidget
}

func (w *PeerInfo) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *PeerInfo) Name() string                   { return w.View.GetTitle() }
func (w *PeerInfo) TextView() *tview.TextView      { return w.View }
func (w *PeerInfo) CommonSettings() *config.Common { return w.Settings.Common }
func (w *PeerInfo) Focusable() bool                { return true }

func (w *PeerInfo) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.Render()
	})
}

func (w *PeerInfo) RefreshInterval() int {
	return w.Settings.Common.RefreshInterval
}

func (w *PeerInfo) Render() {
	text := ""
	var identity types.Identity
	data := []byte{}
	req, err := http.NewRequest("POST", w.Client.Base+w.Settings.URL, nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		goto set
	}
	data, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &identity)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		goto set
	}
	text += fmt.Sprintf("ID : [green]%s[white]\n", identity.ID)
	text += fmt.Sprintf("[green]%s[white] | [green]%s[white]",
		identity.AgentVersion, identity.ProtocolVersion)
set:
	w.View.Clear()
	w.View.SetTitle(w.Settings.Common.Title)
	w.View.SetText(text)
}
