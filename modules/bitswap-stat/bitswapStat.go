package bitswap_stat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/tabwriter"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/Sab94/ipfs-monitor/types"
	"github.com/Sab94/ipfs-monitor/widget"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type BitswapStatBlock widget.Widget

const (
	WidgetName = "bitswapstat"
	URL        = "bitswap/stat"
)

func NewWidget(cfg *config.Config, httpClient *client.HttpClient,
	app *tview.Application) block.Block {
	bsWidget := BitswapStatBlock{
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
	view.SetText(bsWidget.getBitswapStat())

	bsWidget.View = view
	return &bsWidget
}

func (w *BitswapStatBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.View.Clear()
		w.View.SetText(w.getBitswapStat())
	})
}
func (w *BitswapStatBlock) Refreshing() bool {
	return false
}
func (w *BitswapStatBlock) RefreshInterval() int {
	return w.Settings.Common.RefreshInterval
}

func (w *BitswapStatBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *BitswapStatBlock) Name() string                   { return w.View.GetTitle() }
func (w *BitswapStatBlock) TextView() *tview.TextView      { return w.View }
func (w *BitswapStatBlock) CommonSettings() *config.Common { return w.Settings.Common }
func (w *BitswapStatBlock) Focusable() bool                { return true }

func (w *BitswapStatBlock) getBitswapStat() string {
	text := ""
	req, err := http.NewRequest("GET", w.Client.Base+"bitswap/stat", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprintf("[red]Unable to connect to a running ipfs daemon, %s",
			err.Error())
		return text
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var bitswapStat types.BitswapStat
	err = json.Unmarshal(data, &bitswapStat)
	if err != nil {
		fmt.Println(err.Error())
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text
	}
	wrtr := new(tabwriter.Writer)
	var buf bytes.Buffer

	wrtr.Init(&buf, 6, 8, 8, '\t', 0)
	fmt.Fprintf(wrtr, "%12s: [green]%d\t[white]%12s: [green]%d\t[white]%12s: [green]%d\n",
		"Blocks Got", bitswapStat.BlocksReceived, "Blocks Sent",
		bitswapStat.BlocksSent, "Dup Blocks", bitswapStat.DupBlksReceived)
	fmt.Fprintf(wrtr, "%12s: [green]%d\t[white]%12s: [green]%d\t[whitw]%12s: [green]%d\n",
		"Data Got", bitswapStat.DataReceived, "Data Sent", bitswapStat.DataSent, "Dup Dats",
		bitswapStat.DupDataReceived)
	wrtr.Flush()

	return buf.String()
}
