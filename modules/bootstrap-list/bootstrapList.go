package bootstrap_list

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/types"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/config"
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
	bbWidget := BootstrapBlock{
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
	text, count := bbWidget.getBootstrapList()

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

	bbWidget.View = view
	return &bbWidget
}

func (w *BootstrapBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		text, count := w.getBootstrapList()
		w.View.Clear()
		w.View.SetTitle(fmt.Sprintf("%s ([green]%d[white])", w.Config.Title, count))
		w.View.SetText(text)
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

func (w *BootstrapBlock) getBootstrapList() (string, int) {
	text := ""
	req, err := http.NewRequest("GET", w.Client.Base+"bootstrap/list", nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text, 0
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var bootstraps types.BootstrapList
	err = json.Unmarshal(data, &bootstraps)
	if err != nil {
		fmt.Println(err.Error())
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text, 0
	}
	for _, v := range bootstraps.Peers {
		text += fmt.Sprintf("%s\n", v)
	}
	return text, len(bootstraps.Peers)
}
