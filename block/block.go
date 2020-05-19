package block

import (
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Block is the interface that is implemented by modules
type Block interface {
	BorderColor() tcell.Color
	Name() string
	TextView() *tview.TextView
	Focusable() bool
	Render()
	CommonSettings() *config.Common
}

// Schedulable is a module that can be refreshed after
// a fixed time interval
type Schedulable interface {
	Refresh()
	RefreshInterval() int
}
