package bw_stat

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
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell"
	logging "github.com/ipfs/go-log"
	"github.com/rivo/tview"
)

type BWStatBlock widget.Widget

const (
	WidgetName = "bwstat"
	URL        = "stats/bw"
)

var log = logging.Logger("modules/bwstat")

func NewWidget(cfg *config.Config, httpClient *client.HttpClient,
	app *tview.Application) block.Block {
	if !cfg.Monitor.Widgets[WidgetName].Enabled {
		return nil
	}

	w := widget.NewWidget(cfg, httpClient, app, WidgetName, URL)
	bwWidget := BWStatBlock(w)
	bwWidget.Render()

	return &bwWidget
}

func (w *BWStatBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.Render()
	})
}

func (w *BWStatBlock) RefreshInterval() int {
	return w.Settings.Common.RefreshInterval
}

func (w *BWStatBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *BWStatBlock) Name() string                   { return w.View.GetTitle() }
func (w *BWStatBlock) TextView() *tview.TextView      { return w.View }
func (w *BWStatBlock) CommonSettings() *config.Common { return w.Settings.Common }
func (w *BWStatBlock) Focusable() bool                { return true }

func (w *BWStatBlock) Render() {
	text := ""
	data := []byte{}
	var bwStat types.BWStat

	req, err := http.NewRequest("POST", w.Client.Base+"stats/bw", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprintf("[red]Unable to connect to a running ipfs daemon, %s",
			err.Error())
		goto set
	}
	data, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &bwStat)
	if err != nil {
		text += fmt.Sprintf("[red]Unable to connect to a running ipfs daemon")
		goto set
	}

	text += fmt.Sprintf("%14s: [green]%-9s  [white]%14s: [green]%-9s\n",
		"Rate In", humanize.Bytes(uint64(bwStat.RateIn))+"/s", "Rate Out", humanize.Bytes(uint64(bwStat.RateOut))+"/s")
	text += fmt.Sprintf("%14s: [green]%-7s  [white]%16s: [green]%-7s\n",
		"Data Got", humanize.Bytes(bwStat.TotalIn), "Data Sent", humanize.Bytes(bwStat.TotalOut))
set:
	w.View.Clear()
	w.View.SetTitle(w.Settings.Common.Title)
	w.View.SetText(text)
}
