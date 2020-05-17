package swarm_peers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	logging "github.com/ipfs/go-log"

	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/types"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/Sab94/ipfs-monitor/widget"
	"github.com/gdamore/tcell"
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
	spWidget := SwarmPeersBlock{
		Settings: config.Settings{
			Common: &config.Common{
				PositionSettings: cfg.Monitor.Widgets[WidgetName].PositionSettings,
				Bordered:         false,
				Enabled:          false,
				RefreshInterval:  cfg.Monitor.Widgets[WidgetName].RefreshInterval,
				Title:            cfg.Monitor.Widgets[WidgetName].Title,
			},
			URL: URL,
		},
		Client: httpClient,
		Config: cfg.Monitor.Widgets[WidgetName],
		App:    app,
	}
	text, count := spWidget.getSwarmPeers()

	view := tview.NewTextView()
	view.SetTitle(fmt.Sprintf("%s ([green]%d[white])", cfg.Monitor.Widgets[WidgetName].Title, count))
	view.SetBackgroundColor(tcell.ColorNames[cfg.Monitor.Colors.Background])
	view.SetBorder(true)
	view.SetBorderColor(tcell.ColorNames[cfg.Monitor.Colors.Border.Normal])
	view.SetDynamicColors(true)
	view.SetTextColor(tcell.ColorNames[cfg.Monitor.Colors.Text])
	view.SetTitleColor(tcell.ColorNames[cfg.Monitor.Colors.Text])
	view.SetWrap(false)
	view.SetScrollable(true)
	view.SetText(text)

	spWidget.View = view
	return &spWidget
}

func (w *SwarmPeersBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		text, count := w.getSwarmPeers()
		w.View.Clear()
		w.View.SetTitle(fmt.Sprintf("%s ([green]%d[white])", w.Config.Title, count))
		w.View.SetText(text)
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

func (w *SwarmPeersBlock) getSwarmPeers() (string, int) {
	text := ""
	req, err := http.NewRequest("GET", w.Client.Base+"swarm/peers", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprintf("[red]Unable to connect to a running ipfs daemon, %s",
			err.Error())
		return text, 0
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var swarmPeers types.SwarmPeers
	err = json.Unmarshal(data, &swarmPeers)
	if err != nil {
		fmt.Println(err.Error())
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text, 0
	}
	for _, v := range swarmPeers.Peers {
		text += fmt.Sprintf("%s/%s\n", v.Addr, v.Peer)
	}
	return text, len(swarmPeers.Peers)
}
