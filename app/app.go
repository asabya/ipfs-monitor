package app

import (
	"context"
	"fmt"

	"github.com/Sab94/ipfs-monitor/modules"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"go.uber.org/fx"
)

// TerminalMonitor is the terminal app. It has all the modules
// to be rendered. tview app and http.client to call ipfs apis
type TerminalMonitor struct {
	App        *tview.Application
	HttpClient *client.HttpClient
	Blocks     []block.Block
}

// Start fx app with TerminalMonitor
func Start(ctx context.Context) (*TerminalMonitor, error) {
	tApp := &TerminalMonitor{}

	app := fx.New(
		fx.Provide(config.CreateOrLoadConfigFile),
		client.NewHttpClient(),
		fx.Provide(
			NewTviewGrid,
			NewTviewApp,
			modules.BootstrapAllModules,
		),
		fx.Extract(tApp),
		fx.Invoke(func(grid *tview.Grid) {
			for _, v := range tApp.Blocks {
				grid.AddItem(v.TextView(), v.CommonSettings().Top, v.CommonSettings().Left,
					v.CommonSettings().Height, v.CommonSettings().Width, 0, 0, false)
			}
		}),
	)
	tApp.App.SetInputCapture(tApp.keyboardIntercept)
	if err := app.Start(ctx); err != nil {
		return nil, err
	}

	return tApp, nil
}

// NewTviewGrid creates a tview.Grid
func NewTviewGrid(cfg *config.Config) *tview.Grid {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(tcell.ColorNames[cfg.Tapp.Grid.Background])
	grid.SetColumns(cfg.Tapp.Grid.Columns...)
	grid.SetRows(cfg.Tapp.Grid.Rows...)
	grid.SetBorder(false)

	return grid
}

// NewTviewApp creates tview.Application
func NewTviewApp(cfg *config.Config, grid *tview.Grid) *tview.Application {
	pages := tview.NewPages()
	pages.AddPage("grid", grid, true, true)
	pages.Box.SetBackgroundColor(tcell.ColorNames[cfg.Tapp.Grid.Background])

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
		fmt.Println("KeyCtrlR")
		return nil
	case tcell.KeyTab:
		fmt.Println("Tab")
	case tcell.KeyBacktab:
		fmt.Println("Backtab")
		return nil
	case tcell.KeyEsc:
		t.App.Stop()
	}
	return event
}
