package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

func main() {
	var dbName string

	_, currentFilePath, _, _ := runtime.Caller(0)
	dirpath, _ := filepath.Split(currentFilePath)

	dbName = dirpath + "MyVideos116.db"
	db := NewDb(dbName)
	defer db.Close()
	// fmt.Printf("%v", db.GetMovie())

	// DebugMovieArray(db.GetMovie())
	// DebugMovieStruct(db.GetMovie()[1])

	// ans := Input("Search Movies : ")
	// fmt.Println(ans)
	// DebugMovieArray(db.GetMovieWithQuery(strings.Join(ans, " ")))
	// moviesArray := db.GetMovieWithQuery(strings.Join(ans, " "))
	// log.Info(len(moviesArray))
	// fmt.Println(moviesArray)

	// Some encoding stuff in case the user isn't using UTF-8
	encoding.Register()
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	ui := NewUI()

	// itemsPerPage for CtrlB, CtrlD
	ui.itemsPerPage = 10

	// firstG set to 0
	ui.firstG = false

	// Initialize application
	ui.app = tview.NewApplication()

	// Create input field
	ui.input = tview.NewInputField().
		SetPlaceholder(" Please Enter...").
		SetPlaceholderTextColor(tcell.ColorYellow)

	// Create list view
	ui.list = tview.NewList()

	// Create search result label
	ui.resultLabel = tview.NewTextView().
		SetTextColor(tcell.ColorYellow)

	// Create Info view
	ui.infoView = tview.NewTextView()

	// when the key pressed, it will refresh screen with this func
	ui.list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		ui.resultLabel.SetText(fmt.Sprintf("Search Results : %d, Index : %d", ui.list.GetItemCount(), index))
		ui.infoView.SetText(fmt.Sprintf("%s \n\n %s \n\n %s \n\n Premiered: %s \n\n Rating: %.2f \n\n %s",
			ui.moviesArray[index].C00,
			ui.moviesArray[index].C01,
			ui.moviesArray[index].C03,
			ui.moviesArray[index].Premiered,
			ui.moviesArray[index].Rating,
			ui.moviesArray[index].StrPath,
		)).
			SetBackgroundColor(tcell.ColorBlack)
	})

	ui.list.SetDoneFunc(func() {
		handleQuit(ui)
	})

	// ui.input setdonefunc
	ui.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			name := ui.input.GetText()
			if strings.TrimSpace(name) != "" {
				ui.moviesArray = db.GetMovieWithQuery(name)
			} else {
				ui.moviesArray = db.GetMovies()
			}

			// item loading into ui.list
			for _, item := range ui.moviesArray {
				ui.list.AddItem(item.C00, "", 0, nil)
			}
			ui.input.SetText("")
			ui.resultLabel.SetText(fmt.Sprintf("Search Results : %d, Index : %d", ui.list.GetItemCount(), 0))
			ui.app.SetFocus(ui.list)
		} else if key == tcell.KeyEscape {
			handleQuit(ui)
		}
	})

	// Create Flex container
	displayFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("SEARCH MY MOVIES"), 1, 1, false).
		AddItem(ui.input, 2, 1, true).
		AddItem(ui.resultLabel, 2, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(ui.list, 0, 2, false).
			AddItem(ui.infoView, 0, 3, false).
			AddItem(tview.NewTextView(), 0, 1, false), 0, 1, false)
		// last TextView is a right margin for infoView

	// Capture user input
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// Anything handled here will be executed on the main thread
		switch event.Key() {
		case tcell.KeyCtrlD:
			ui.list.SetCurrentItem(ui.list.GetCurrentItem() + ui.itemsPerPage)
			if ui.list.GetCurrentItem() >= ui.list.GetItemCount() {
				ui.list.SetCurrentItem(ui.list.GetItemCount() - 1)
			}
		case tcell.KeyCtrlB:
			ui.list.SetCurrentItem(ui.list.GetCurrentItem() - ui.itemsPerPage)
			if ui.list.GetCurrentItem() < 0 {
				ui.list.SetCurrentItem(0)
			}
		case tcell.KeyRune:
			ch := event.Rune()
			if ch == 'J' || ch == 'j' {
				if ui.list.GetCurrentItem() == (ui.list.GetItemCount() - 1) {
					ui.list.SetCurrentItem(0)
				} else {
					ui.list.SetCurrentItem(ui.list.GetCurrentItem() + 1)
				}
			}

			if ch == 'K' || ch == 'k' {
				ui.list.SetCurrentItem(ui.list.GetCurrentItem() - 1)
			}

			// going to last items
			if ch == 'G' {
				ui.list.SetCurrentItem(ui.list.GetItemCount() - 1)
			}

			// goint to first items
			if ch == 'g' {
				if ui.firstG == false {
					ui.firstG = true
				} else {
					ui.list.SetCurrentItem(0)
					ui.firstG = false
				}

			} else if ui.firstG == true {
				ui.firstG = false
			}

			if ch == 'Q' || ch == 'q' {
				handleQuit(ui)
			}

		}
		return event
	})

	// Set the grid as the application root and focus the input field
	ui.app.SetRoot(displayFlex, true).SetFocus(ui.input)

	// Run the application
	log.Info("tview Loading")
	err := ui.app.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("tview Closed")
}

func handleQuit(ui *UI) {
	// Handling the Quit or ESC key
	if ui.list.HasFocus() {
		// Clear the input field
		ui.input.SetText("")
		ui.infoView.Clear()
		ui.list.Clear()
		ui.resultLabel.Clear()
		ui.moviesArray = nil

		// focus the input field
		ui.app.SetFocus(ui.input)
	} else if ui.input.HasFocus() {
		ui.app.Stop()
	}
}
