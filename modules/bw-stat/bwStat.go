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
	bwWidget := BWStatBlock{
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

	view := tview.NewTextView()
	view.SetTitle(cfg.Monitor.Widgets[WidgetName].Title)
	view.SetBackgroundColor(tcell.ColorNames[cfg.Monitor.Colors.Background])
	view.SetBorder(true)
	view.SetBorderColor(tcell.ColorNames[cfg.Monitor.Colors.Border.Normal])
	view.SetDynamicColors(true)
	view.SetTextColor(tcell.ColorNames[cfg.Monitor.Colors.Text])
	view.SetTitleColor(tcell.ColorNames[cfg.Monitor.Colors.Text])
	view.SetWrap(false)
	view.SetScrollable(true)
	view.SetText(bwWidget.getBitswapStat())

	bwWidget.View = view
	return &bwWidget
}

func (w *BWStatBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.View.Clear()
		w.View.SetText(w.getBitswapStat())
	})
}
func (w *BWStatBlock) Refreshing() bool {
	return false
}
func (w *BWStatBlock) RefreshInterval() int {
	return w.Settings.Common.RefreshInterval
}

func (w *BWStatBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *BWStatBlock) Name() string                   { return w.View.GetTitle() }
func (w *BWStatBlock) TextView() *tview.TextView      { return w.View }
func (w *BWStatBlock) CommonSettings() *config.Common { return w.Settings.Common }
func (w *BWStatBlock) Focusable() bool                { return true }

func (w *BWStatBlock) getBitswapStat() string {
	text := ""
	req, err := http.NewRequest("GET", w.Client.Base+"stats/bw", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprintf("[red]Unable to connect to a running ipfs daemon, %s",
			err.Error())
		return text
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var bwStat types.BWStat
	err = json.Unmarshal(data, &bwStat)
	if err != nil {
		fmt.Println(err.Error())
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text
	}
	wrtr := new(tabwriter.Writer)
	var buf bytes.Buffer

	// minwidth, tabwidth, padding, padchar, flags
	wrtr.Init(&buf, 6, 8, 10, '\t', 0)
	fmt.Fprintf(wrtr, "Rate In : %f\t Rate Out : %f\n",
		bwStat.RateIn, bwStat.RateOut)
	fmt.Fprintf(wrtr, "Data Got : %d\t Data Sent : %d\n",
		bwStat.TotalIn, bwStat.TotalOut)
	wrtr.Flush()
	return buf.String()
}
