package repo_stat

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

type RepoStatBlock widget.Widget

const (
	WidgetName = "repostat"
	URL        = "repo/stat"
)

func NewWidget(cfg *config.Config, hc *client.HttpClient) block.Block {
	rsWidget := RepoStatBlock{
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

	view.SetText(rsWidget.getRepoStat(hc))

	rsWidget.View = view
	return &rsWidget
}

func (rs *RepoStatBlock) Refresh() {
	fmt.Println("Refreshing Repo stat")
}
func (bb *RepoStatBlock) Refreshing() bool {
	return false
}
func (bb *RepoStatBlock) RefreshInterval() int {
	return 10
}

func (w *RepoStatBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *RepoStatBlock) Name() string                   { return w.View.GetTitle() }
func (w *RepoStatBlock) TextView() *tview.TextView      { return w.View }
func (w *RepoStatBlock) CommonSettings() *config.Common { return w.Settings.Common }

func (w *RepoStatBlock) getRepoStat(client *client.HttpClient) string {
	text := ""
	req, err := http.NewRequest("GET", client.Base+w.Settings.URL, nil)
	resp, err := client.Client.Do(req)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var repoStat types.RepoStat
	err = json.Unmarshal(data, &repoStat)
	if err != nil {
		fmt.Println(err.Error())
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		return text
	}
	text += fmt.Sprintf("Path : [green]%s[white]\n", repoStat.RepoPath)
	text += fmt.Sprintf("Size : [green]%d[white]\n", repoStat.RepoSize)
	text += fmt.Sprintf("StorageMax : [green]%d[white]", repoStat.StorageMax)
	return text
}
