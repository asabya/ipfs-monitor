package peer_info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sab94/ipfs-monitor/types"

	"github.com/Sab94/ipfs-monitor/client"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/Sab94/ipfs-monitor/widget"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type PeerInfo widget.Widget

const (
	WidgetName = "peerinfo"
	URL        = "id"
)

func NewWidget(cfg *config.Config, httpClient *client.HttpClient) block.Block {
	piWidget := PeerInfo{
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

	view := tview.NewTextView()
	view.SetTitle(cfg.Tapp.Widgets[WidgetName].Title)
	view.SetBackgroundColor(tcell.ColorNames[cfg.Tapp.Colors.Background])
	view.SetBorder(true)
	view.SetBorderColor(tcell.ColorNames[cfg.Tapp.Colors.Border.Normal])
	view.SetDynamicColors(true)
	view.SetTextColor(tcell.ColorNames[cfg.Tapp.Colors.Text])
	view.SetTitleColor(tcell.ColorNames[cfg.Tapp.Colors.Text])
	view.SetWrap(false)
	view.SetScrollable(true)

	view.SetText(piWidget.getPeerInfo(httpClient))

	piWidget.View = view
	return &piWidget
}

func (w *PeerInfo) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *PeerInfo) Name() string                   { return w.View.GetTitle() }
func (w *PeerInfo) TextView() *tview.TextView      { return w.View }
func (w *PeerInfo) CommonSettings() *config.Common { return w.Settings.Common }

func (w *PeerInfo) getPeerInfo(client *client.HttpClient) string {
	text := ""
	req, err := http.NewRequest("GET", client.Base+w.Settings.URL, nil)
	resp, err := client.Client.Do(req)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text
	}
	data, _ := ioutil.ReadAll(resp.Body)

	var identity types.Identity
	err = json.Unmarshal(data, &identity)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text
	}
	text += fmt.Sprintf("ID : [green]%s[white]\n", identity.ID)
	text += fmt.Sprintf("Version : [green]%s[white]", identity.ProtocolVersion)
	return text
}
