package bitswap_stat

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

var log = logging.Logger("bitswapStat")

type BitswapStatBlock widget.Widget

const (
	WidgetName = "bitswapstat"
	URL        = "bitswap/stat"
)

func NewWidget(cfg *config.Config, httpClient *client.HttpClient,
	app *tview.Application) block.Block {
	if !cfg.Monitor.Widgets[WidgetName].Enabled {
		return nil
	}

	w := widget.NewWidget(cfg, httpClient, app, WidgetName, URL)
	bsWidget := BitswapStatBlock(w)
	bsWidget.Render()
	return &bsWidget
}

func (w *BitswapStatBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.Render()
	})
}

func (w *BitswapStatBlock) RefreshInterval() int {
	return w.Settings.Common.RefreshInterval
}

func (w *BitswapStatBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *BitswapStatBlock) Name() string                   { return w.View.GetTitle() }
func (w *BitswapStatBlock) TextView() *tview.TextView      { return w.View }
func (w *BitswapStatBlock) CommonSettings() *config.Common { return w.Settings.Common }
func (w *BitswapStatBlock) Focusable() bool                { return true }

func (w *BitswapStatBlock) Render() {
	text := ""
	var data = []byte{}
	var bitswapStat types.BitswapStat

	req, err := http.NewRequest("POST", w.Client.Base+"bitswap/stat", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprintf( "[red]Unable to connect to a running ipfs daemon, %s", err.Error())
		goto set
	}
	data, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &bitswapStat)
	if err != nil {
		text += fmt.Sprintf("[red]Unable to connect to a running ipfs daemon, %s", err.Error())
		goto set
	}

	text += fmt.Sprintf("%12s: [green]%-7d[white]%12s: [green]%-7d[white]%12s: [green]%-7d\n",
		"Blocks Got", bitswapStat.BlocksReceived, "Blocks Sent",
		bitswapStat.BlocksSent, "Dup Blocks", bitswapStat.DupBlksReceived)
	text += fmt.Sprintf(	"%12s: [green]%-7d[white]%12s: [green]%-7d[whitw]%12s: [green]%-7d\n",
		"Data Got", bitswapStat.DataReceived, "Data Sent", bitswapStat.DataSent, "Dup Dats",
		bitswapStat.DupDataReceived)

set:
	w.View.Clear()
	w.View.SetTitle(w.Settings.Common.Title)
	w.View.SetText(text)
}
