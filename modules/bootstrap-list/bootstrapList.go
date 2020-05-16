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

func NewWidget(cfg *config.Config, httpClient *client.HttpClient) block.Block {
	bbWidget := BootstrapBlock{
		Settings: config.Settings{
			Common: &config.Common{
				PositionSettings: cfg.Tapp.Widgets[WidgetName].PositionSettings,
				Bordered:         false,
				Enabled:          false,
				RefreshInterval:  0,
				Title:            cfg.Tapp.Widgets[WidgetName].Title,
			},
			URL: URL,
		},
	}
	text, count := getBootstrapList(httpClient)

	view := tview.NewTextView()
	view.SetTitle(fmt.Sprintf("%s (%d)", cfg.Tapp.Widgets[WidgetName].Title, count))
	view.SetBackgroundColor(tcell.ColorNames[cfg.Tapp.Colors.Background])
	view.SetBorder(true)
	view.SetBorderColor(tcell.ColorNames[cfg.Tapp.Colors.Border.Normal])
	view.SetDynamicColors(true)
	view.SetTextColor(tcell.ColorNames[cfg.Tapp.Colors.Text])
	view.SetTitleColor(tcell.ColorNames[cfg.Tapp.Colors.Text])
	view.SetWrap(false)
	view.SetScrollable(true)
	view.SetText(text)

	bbWidget.View = view
	return &bbWidget
}

func (bb *BootstrapBlock) Refresh() {
	fmt.Println("Refreshing bootstrapblock")
}
func (bb *BootstrapBlock) Refreshing() bool {
	return false
}
func (bb *BootstrapBlock) RefreshInterval() int {
	return 5
}

func (w *BootstrapBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *BootstrapBlock) Name() string                   { return w.View.GetTitle() }
func (w *BootstrapBlock) TextView() *tview.TextView      { return w.View }
func (w *BootstrapBlock) CommonSettings() *config.Common { return w.Settings.Common }

func getBootstrapList(client *client.HttpClient) (string, int) {
	text := ""
	req, err := http.NewRequest("GET", client.Base+"bootstrap/list", nil)
	resp, err := client.Client.Do(req)
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
