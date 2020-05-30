package repo_stat

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

	w := widget.NewWidget(cfg, hc, app, WidgetName, URL)
	rsWidget := RepoStatBlock(w)
	rsWidget.Render()
	return &rsWidget
}

func (w *RepoStatBlock) Refresh() {
	w.App.QueueUpdateDraw(func() {
		w.Render()
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

func (w *RepoStatBlock) Render()  {
	text := ""
	var repoStat types.RepoStat
	data := []byte{}

	req, err := http.NewRequest("POST", w.Client.Base+w.Settings.URL, nil)
	resp, err := w.Client.Client.Do(req)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		goto set
	}
	data, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &repoStat)
	if err != nil {
		text += fmt.Sprint("[red]Unable to connect to a running ipfs daemon")
		goto set
	}
	text += fmt.Sprintf("%10s: [green]%s[white]\n", "Path", repoStat.RepoPath)
	text += fmt.Sprintf("%10s: [green]%s[white]\n", "Size", humanize.Bytes(repoStat.RepoSize))
	text += fmt.Sprintf("%10s: [green]%s[white]", "StorageMax", humanize.Bytes(repoStat.StorageMax))
set:
	w.View.Clear()
	w.View.SetTitle(w.Settings.Common.Title)
	w.View.SetText(text)
}
