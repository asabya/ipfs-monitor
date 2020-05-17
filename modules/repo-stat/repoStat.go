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

func NewWidget(cfg *config.Config, hc *client.HttpClient,
	app *tview.Application) block.Block {
	if !cfg.Monitor.Widgets[WidgetName].Enabled {
		return nil
	}
	rsWidget := RepoStatBlock{
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
		Client: hc,
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

	view.SetText(rsWidget.getRepoStat())

	rsWidget.View = view
	return &rsWidget
}

func (w *RepoStatBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.View.Clear()
		w.View.SetText(w.getRepoStat())
	})
}

func (w *RepoStatBlock) RefreshInterval() int {
	return w.Settings.Common.RefreshInterval
}

func (w *RepoStatBlock) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *RepoStatBlock) Name() string                   { return w.View.GetTitle() }
func (w *RepoStatBlock) TextView() *tview.TextView      { return w.View }
func (w *RepoStatBlock) CommonSettings() *config.Common { return w.Settings.Common }
func (w *RepoStatBlock) Focusable() bool                { return true }

func (w *RepoStatBlock) getRepoStat() string {
	text := ""
	req, err := http.NewRequest("GET", w.Client.Base+w.Settings.URL, nil)
	resp, err := w.Client.Client.Do(req)
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
