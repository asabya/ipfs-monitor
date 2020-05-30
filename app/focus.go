package app

import (
	"sort"

	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/config"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"go.uber.org/fx"
)

type FocusState int

const (
	widgetFocused FocusState = iota
	appBoardFocused
	neverFocused
)

type FocusTracker struct {
	App       *tview.Application
	Idx       int
	IsFocused bool
	Widgets   []block.Block

	Config *config.Config
}

func NewFocusTracker() fx.Option {
	return fx.Provide(func(app *tview.Application, widgets []block.Block, config *config.Config) *FocusTracker {
		return &FocusTracker{
			App:       app,
			Idx:       -1,
			IsFocused: false,
			Widgets:   widgets,

			Config: config,
		}
	})
}

// Next sets the focus on the next widget in the widget list. If the current widget is
// the last widget, sets focus on the first widget.
func (tracker *FocusTracker) Next() {
	if tracker.focusState() == appBoardFocused {
		return
	}

	tracker.blur(tracker.Idx)
	tracker.increment()
	tracker.focus(tracker.Idx)

	tracker.IsFocused = true
}

// None removes focus from the currently-focused widget.
func (tracker *FocusTracker) None() {
	if tracker.focusState() == appBoardFocused {
		return
	}

	tracker.blur(tracker.Idx)
}

// Prev sets the focus on the previous widget in the widget list. If the current widget is
// the last widget, sets focus on the last widget.
func (tracker *FocusTracker) Prev() {
	if tracker.focusState() == appBoardFocused {
		return
	}

	tracker.blur(tracker.Idx)
	tracker.decrement()
	tracker.focus(tracker.Idx)

	tracker.IsFocused = true
}

func (tracker *FocusTracker) Refocus() {
	tracker.focus(tracker.Idx)
}

func (tracker *FocusTracker) blur(idx int) {
	widget := tracker.focusableAt(idx)
	if widget == nil {
		return
	}

	view := widget.TextView()
	view.Blur()

	view.SetBorderColor(
		tcell.ColorNames[tracker.Config.Monitor.Colors.Border.Normal],
	)

	tracker.IsFocused = false
}

func (tracker *FocusTracker) decrement() {
	tracker.Idx--

	if tracker.Idx < 0 {
		tracker.Idx = len(tracker.focusables()) - 1
	}
}

func (tracker *FocusTracker) focus(idx int) {
	widget := tracker.focusableAt(idx)
	if widget == nil {
		return
	}

	view := widget.TextView()
	view.SetBorderColor(
		tcell.ColorNames[tracker.Config.Monitor.Colors.Border.Focused],
	)
	tracker.App.SetFocus(view)
}

func (tracker *FocusTracker) focusables() []block.Block {
	focusable := []block.Block{}

	for _, widget := range tracker.Widgets {
		if widget.Focusable() {
			focusable = append(focusable, widget)
		}
	}

	// Sort for deterministic ordering
	sort.SliceStable(focusable[:], func(i, j int) bool {
		iTop := focusable[i].CommonSettings().Top
		jTop := focusable[j].CommonSettings().Top

		if iTop < jTop {
			return true
		}
		if iTop == jTop {
			return focusable[i].CommonSettings().Left < focusable[j].CommonSettings().Left
		}
		return false
	})

	return focusable
}

func (tracker *FocusTracker) focusableAt(idx int) block.Block {
	if idx < 0 || idx >= len(tracker.focusables()) {
		return nil
	}

	return tracker.focusables()[idx]
}

func (tracker *FocusTracker) focusState() FocusState {
	if tracker.Idx < 0 {
		return neverFocused
	}

	for _, widget := range tracker.Widgets {
		if widget.TextView() == tracker.App.GetFocus() {
			return widgetFocused
		}
	}

	return appBoardFocused
}

func (tracker *FocusTracker) increment() {
	tracker.Idx++

	if tracker.Idx == len(tracker.focusables()) {
		tracker.Idx = 0
	}
}
