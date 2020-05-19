package bw_stat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/tabwriter"

	logging "github.com/ipfs/go-log"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/Sab94/ipfs-monitor/types"
	"github.com/Sab94/ipfs-monitor/widget"
	"github.com/gdamore/tcell"
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
	wrtr := new(tabwriter.Writer)
	var buf bytes.Buffer
	var bwStat types.BWStat
	wrtr.Init(&buf, 6, 8, 8, '\t', 0)
	data := []byte{}
	req, err := http.NewRequest("GET", w.Client.Base+"stats/bw", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		fmt.Fprintf(wrtr, "[red]Unable to connect to a running ipfs daemon, %s",
			err.Error())
		goto set
	}
	data, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &bwStat)
	if err != nil {
		fmt.Fprintf(wrtr, "[red]Unable to connect to a running ipfs daemon")
		goto set
	}

	fmt.Fprintf(wrtr, "%12s: [green]%.3f  [white]%12s: [green]%.3f\n",
		"Rate In", bwStat.RateIn, "Rate Out", bwStat.RateOut)
	fmt.Fprintf(wrtr, "%12s: [green]%d  [white]%12s: [green]%d\n",
		"Data Got", bwStat.TotalIn, "Data Sent", bwStat.TotalOut)
set:
	wrtr.Flush()
	w.View.Clear()
	w.View.SetTitle(w.Settings.Common.Title)
	w.View.SetText(buf.String())
}
