package main

import (
	"github.com/rivo/tview"
)

// UI of Application
type UI struct {
	app         *tview.Application
	input       *tview.InputField
	list        *tview.List
	resultLabel *tview.TextView
	modalView   *tview.Modal
}

// NewUI return UI struct
func NewUI() *UI {
	ui := &UI{}
	return ui
}
