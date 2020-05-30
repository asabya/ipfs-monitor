package bootstrap_list

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
	"github.com/rivo/tview"
)

type BootstrapBlock widget.Widget

const (
	WidgetName = "bootstraplist"
	URL        = "bootstrap/list"
)

func NewWidget(cfg *config.Config, httpClient *client.HttpClient,
	app *tview.Application) block.Block {
	if !cfg.Monitor.Widgets[WidgetName].Enabled {
		return nil
	}

	w := widget.NewWidget(cfg, httpClient, app, WidgetName, URL)
	bbWidget := BootstrapBlock(w)
	bbWidget.Render()

	return &bbWidget
}

func (w *BootstrapBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.Render()
	})
}

func (w *BootstrapBlock) RefreshInterval() int {
	return w.Settings.Common.RefreshInterval
}

func (w *BootstrapBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *BootstrapBlock) Name() string                   { return w.View.GetTitle() }
func (w *BootstrapBlock) TextView() *tview.TextView      { return w.View }
func (w *BootstrapBlock) CommonSettings() *config.Common { return w.Settings.Common }
func (w *BootstrapBlock) Focusable() bool                { return true }

func (w *BootstrapBlock) Render() {
	w.View.Clear()
	text := ""
	var bootstraps types.BootstrapList
	var data = []byte{}
	req, err := http.NewRequest("POST", w.Client.Base+"bootstrap/list", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		goto set
	}
	data, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &bootstraps)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		goto set
	}
	for _, v := range bootstraps.Peers {
		text += fmt.Sprintf("%s\n", v)
	}
set:
	w.View.SetTitle(fmt.Sprintf("%s ([green]%d[white])", w.Settings.Common.Title, len(bootstraps.Peers)))
	w.View.SetText(text)
	return
}
