package swarm_peers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/Sab94/ipfs-monitor/types"
	"github.com/Sab94/ipfs-monitor/widget"
	"github.com/gdamore/tcell"
	logging "github.com/ipfs/go-log"
	"github.com/rivo/tview"
)

type SwarmPeersBlock widget.Widget

var log = logging.Logger("modules/swarmpeers")

const (
	WidgetName = "swarmpeers"
	URL        = "swarm/peers"
)

func NewWidget(cfg *config.Config, httpClient *client.HttpClient,
	app *tview.Application) block.Block {
	if !cfg.Monitor.Widgets[WidgetName].Enabled {
		return nil
	}

	w := widget.NewWidget(cfg, httpClient, app, WidgetName, URL)
	spWidget := SwarmPeersBlock(w)
	spWidget.Render()
	return &spWidget
}

func (w *SwarmPeersBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.Render()
	})
}

func (w *SwarmPeersBlock) RefreshInterval() int {
	return w.Settings.Common.RefreshInterval
}

func (w *SwarmPeersBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *SwarmPeersBlock) Name() string                   { return w.View.GetTitle() }
func (w *SwarmPeersBlock) TextView() *tview.TextView      { return w.View }
func (w *SwarmPeersBlock) CommonSettings() *config.Common { return w.Settings.Common }
func (w *SwarmPeersBlock) Focusable() bool                { return true }

func (w *SwarmPeersBlock) Render() {
	text := ""
	var swarmPeers types.SwarmPeers
	data := []byte{}
	req, err := http.NewRequest("POST", w.Client.Base+"swarm/peers", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprintf("[red]Unable to connect to a running ipfs daemon, %s",
			err.Error())
		goto set
	}
	data, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(data, &swarmPeers)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		goto set
	}
	for _, v := range swarmPeers.Peers {
		text += fmt.Sprintf("%s/%s\n", v.Addr, v.Peer)
	}
set:
	w.View.Clear()
	w.View.SetTitle(fmt.Sprintf("%s ([green]%d[white])", w.Settings.Common.Title, len(swarmPeers.Peers)))
	w.View.SetText(text)
}
