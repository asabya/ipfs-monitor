package widget

import (
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Widget struct {
	View     *tview.TextView
	Settings config.Settings
	Client   *client.HttpClient
	Config   config.WidgetConfigs
	App      *tview.Application
}

func (w *Widget) Refresh()             {}
func (w *Widget) Refreshing() bool     { return false }
func (w *Widget) RefreshInterval() int { return w.Settings.Common.RefreshInterval }

func (w *Widget) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *Widget) Name() string                   { return w.View.GetTitle() }
func (w *Widget) TextView() *tview.TextView      { return w.View }
func (w *Widget) CommonSettings() *config.Common { return w.Settings.Common }
