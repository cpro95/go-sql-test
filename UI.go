package main

import (
	"github.com/rivo/tview"
)

// UI of Application
type UI struct {
	app          *tview.Application
	input        *tview.InputField
	list         *tview.List
	resultLabel  *tview.TextView
	infoView     *tview.TextView
	moviesArray  []Movie
	itemsPerPage int
	firstG       bool
}

// NewUI return UI struct
func NewUI() *UI {
	ui := &UI{}
	return ui
}
