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

func NewWidget(cfg *config.Config, httpClient *client.HttpClient,
	app *tview.Application, widgetName, url string) Widget {

	view := tview.NewTextView()
	view.SetBackgroundColor(tcell.ColorNames[cfg.Monitor.Colors.Background])
	view.SetBorder(true)
	view.SetBorderColor(tcell.ColorNames[cfg.Monitor.Colors.Border.Normal])
	view.SetDynamicColors(true)
	view.SetTextColor(tcell.ColorNames[cfg.Monitor.Colors.Text])
	view.SetTitleColor(tcell.ColorNames[cfg.Monitor.Colors.Text])
	view.SetWrap(false)
	view.SetScrollable(true)

	return Widget{
		Settings: config.Settings{
			Common: &config.Common{
				PositionSettings: cfg.Monitor.Widgets[widgetName].PositionSettings,
				Bordered:         false,
				Enabled:          false,
				RefreshInterval:  cfg.Monitor.Widgets[widgetName].RefreshInterval,
				Title:            cfg.Monitor.Widgets[widgetName].Title,
			},
			URL: url,
		},
		Client: httpClient,
		Config: cfg.Monitor.Widgets[widgetName],
		App:    app,
		View: view,
	}
}

func (w *Widget) BorderColor() tcell.Color       { return w.View.GetBorderColor() }
func (w *Widget) Name() string                   { return w.View.GetTitle() }
func (w *Widget) TextView() *tview.TextView      { return w.View }
func (w *Widget) CommonSettings() *config.Common { return w.Settings.Common }
func (w *Widget) Focusable() bool { return true }
