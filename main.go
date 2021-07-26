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

	// Create modal view
	ui.modalView = tview.NewModal()

	// when the key pressed, it will refresh screen with this func
	ui.list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		ui.resultLabel.SetText(fmt.Sprintf("Search Results : %d, Index : %d", ui.list.GetItemCount(), index))
	})

	// Create Flex container
	displayFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("SEARCH MY MOVIES"), 1, 1, false).
		AddItem(ui.input, 2, 1, true).
		AddItem(ui.resultLabel, 2, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(ui.list, 0, 1, false).
			AddItem(tview.NewTextView().SetText("Info"), 0, 1, false), 10, 1, false)

	// listFlex := tview.NewFlex().SetDirection(tview.FlexRow).
	// 	AddItem(ui.resultLabel, 0, 1, false).
	// 	AddItem(ui.list, 0, 10, false).
	// 	AddItem(tview.NewTextView().SetText("Info"), 20, 1, false)

	var moviesArray []Movie

	// Capture user input
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// Anything handled here will be executed on the main thread
		switch event.Key() {
		case tcell.KeyEnter:
			name := ui.input.GetText()
			if strings.TrimSpace(name) != "" {
				moviesArray = db.GetMovieWithQuery(name)
				// ui.resultLabel.SetText(fmt.Sprintf("Search Results : %d, Index : %d", len(moviesArray), ui.list.GetCurrentItem()))
				for _, item := range moviesArray {
					ui.list.AddItem(item.C00, "", 0, nil)
				}
				ui.input.SetText("")
				ui.app.SetRoot(displayFlex, true).SetFocus(ui.list)
			} else if len(moviesArray) > 0 {
				// Create a modal dialog
				ui.modalView.SetText(fmt.Sprintf("%s \n\n %s \n\n %s \n\n Premiered: %s \n\n Rating: %.2f \n\n %s",
					moviesArray[ui.list.GetCurrentItem()].C00,
					moviesArray[ui.list.GetCurrentItem()].C01,
					moviesArray[ui.list.GetCurrentItem()].C03,
					moviesArray[ui.list.GetCurrentItem()].Premiered,
					moviesArray[ui.list.GetCurrentItem()].Rating,
					moviesArray[ui.list.GetCurrentItem()].StrPath,
				)).
					SetBackgroundColor(tcell.ColorBlack)

				// Display and focus the dialog
				ui.app.SetRoot(ui.modalView, true).SetFocus(ui.modalView)
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

			if ch == 'Q' || ch == 'q' {
				if displayFlex.HasFocus() {
					// Clear the input field
					ui.input.SetText("")
					ui.list.Clear()
					moviesArray = nil

					// Display appGrid and focus the input field
					ui.app.SetRoot(displayFlex, true).SetFocus(ui.input)
				} else if ui.modalView.HasFocus() {
					ui.app.SetRoot(displayFlex, true).SetFocus(ui.list)
				} else {
					ui.app.Stop()
				}
			}

		case tcell.KeyEsc:
			//Exit the application
			if displayFlex.HasFocus() {
				// Clear the input field
				ui.input.SetText("")
				ui.list.Clear()
				moviesArray = nil

				// Display appGrid and focus the input field
				ui.app.SetRoot(displayFlex, true).SetFocus(ui.input)
			} else if ui.modalView.HasFocus() {
				ui.app.SetRoot(displayFlex, true).SetFocus(ui.list)
			} else {
				ui.app.Stop()
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
