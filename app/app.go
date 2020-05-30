package app

import (
	"context"
	"time"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/Sab94/ipfs-monitor/modules"
	"github.com/gdamore/tcell"
	logging "github.com/ipfs/go-log"
	"github.com/rivo/tview"
	"go.uber.org/fx"
)

var log = logging.Logger("subsystem name")

// TerminalMonitor is the terminal app for ipfs. For now
// it has some of the ipfs apis implemented and rendered
type TerminalMonitor struct {
	App          *tview.Application
	HttpClient   *client.HttpClient
	Blocks       []block.Block
	FocusTracker *FocusTracker
}

// Bootstrap TerminalMonitor using fx
func Bootstrap(ctx context.Context) (*TerminalMonitor, error) {
	logging.SetLogLevel("*", "Debug")

	monitor := &TerminalMonitor{}
	app := fx.New(
		fx.Provide(config.CreateOrLoadConfigFile),
		client.NewHttpClient(),
		fx.Provide(
			NewTviewGrid,
			NewTviewApp,
			modules.BootstrapAllModules,
		),
		NewFocusTracker(),
		fx.Extract(monitor),
		fx.Invoke(func(grid *tview.Grid) {
			for _, v := range monitor.Blocks {
				grid.AddItem(v.TextView(), v.CommonSettings().Top, v.CommonSettings().Left,
					v.CommonSettings().Height, v.CommonSettings().Width, 0, 0, false)
			}
			monitor.scheduleWidgets()
		}),
	)
	monitor.App.SetInputCapture(monitor.keyboardIntercept)

	if err := app.Start(ctx); err != nil {
		return nil, err
	}
	return monitor, nil
}

// NewTviewGrid creates a tview.Grid
func NewTviewGrid(cfg *config.Config) *tview.Grid {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(tcell.ColorNames[cfg.Monitor.Grid.Background])
	grid.SetColumns(cfg.Monitor.Grid.Columns...)
	grid.SetRows(cfg.Monitor.Grid.Rows...)
	grid.SetBorder(false)

	return grid
}

// NewTviewApp creates tview.Application
func NewTviewApp(cfg *config.Config, grid *tview.Grid) *tview.Application {
	pages := tview.NewPages()
	pages.AddPage("grid", grid, true, true)
	pages.Box.SetBackgroundColor(tcell.ColorNames[cfg.Monitor.Grid.Background])

	tviewApp := tview.NewApplication()
	tviewApp.SetRoot(pages, true)

	return tviewApp
}

// keyboardIntercept listens for keyboard inputs
// based on key combination different functions are called
func (t *TerminalMonitor) keyboardIntercept(event *tcell.EventKey) *tcell.EventKey {
	// These keys are global keys used by the app. Widgets should not implement these keys
	switch event.Key() {
	case tcell.KeyCtrlR:
		t.refreshAllWidgets()
		return nil
	case tcell.KeyTab:
		t.FocusTracker.Next()
	case tcell.KeyBacktab:
		t.FocusTracker.Prev()
		return nil
	case tcell.KeyEsc:
		t.App.Stop()
	}
	return event
}

func (t *TerminalMonitor) refreshAllWidgets() {
	for _, widget := range t.Blocks {
		schedulable, ok := widget.(block.Schedulable)
		if ok {
			go schedulable.Refresh()
		}
	}
}

func (t *TerminalMonitor) scheduleWidgets() {
	for _, widget := range t.Blocks {
		schedulable, ok := widget.(block.Schedulable)
		if ok {
			log.Debugf("scheduleWidgets : %s %v", widget.Name(), ok)
			go Schedule(schedulable)
		}
	}
}

// Schedule queues widgets for data refresh on a timer
func Schedule(widget block.Schedulable) {
	if widget.RefreshInterval() <= 0 {
		return
	}
	interval := time.Duration(widget.RefreshInterval()) * time.Second
	timer := time.NewTicker(interval)

	for {
		select {
		case <-timer.C:
			widget.Refresh()
		}
	}
}
