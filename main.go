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

	// Initialize application
	app := tview.NewApplication()

	// Create input field
	input := tview.NewInputField().
		SetPlaceholder(" Please Enter...").
		SetPlaceholderTextColor(tcell.ColorYellow)

	list := tview.NewList()

	// Create search result label
	resultLabel := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Search Results")

	// list.SetSelectedFunc(func(int, string, string, rune) {
	// 	resultLabel.SetText("dddd")
	// })

	// Create Flex container
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("SEARCH MY MOVIES"), 2, 0, false).AddItem(input, 0, 1, false)

	listFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(resultLabel, 2, 0, false).AddItem(list, 0, 1, false)

	modalView := tview.NewModal()
	var moviesArray []Movie
	// Capture user input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		switch event.Key() {
		case tcell.KeyEnter:
			name := input.GetText()
			if strings.TrimSpace(name) != "" {
				moviesArray = db.GetMovieWithQuery(name)
				resultLabel.SetText(fmt.Sprintf("Search Results : %d", len(moviesArray)))
				for _, item := range moviesArray {
					list.AddItem(item.C00, "", 0, nil)
				}
				input.SetText("")
				app.SetRoot(listFlex, true).SetFocus(list)
			} else if len(moviesArray) > 0 {
				// Create a modal dialog
				modalView.SetText(fmt.Sprintf("%s \n\n %s \n\n %s \n\n Premiered: %s \n\n Rating: %.2f \n\n %s",
					moviesArray[list.GetCurrentItem()].C00,
					moviesArray[list.GetCurrentItem()].C01,
					moviesArray[list.GetCurrentItem()].C03,
					moviesArray[list.GetCurrentItem()].Premiered,
					moviesArray[list.GetCurrentItem()].Rating,
					moviesArray[list.GetCurrentItem()].StrPath,
				)).
					SetBackgroundColor(tcell.ColorBlack)

				// Display and focus the dialog
				app.SetRoot(modalView, true).SetFocus(modalView)
			}

			return nil

		case tcell.KeyEsc:
			//Exit the application
			if listFlex.HasFocus() {
				// Clear the input field
				input.SetText("")
				list.Clear()
				moviesArray = nil

				// Display appGrid and focus the input field
				app.SetRoot(flex, true).SetFocus(input)
			} else if modalView.HasFocus() {
				app.SetRoot(listFlex, true).SetFocus(list)
			} else {
				app.Stop()
			}

			return nil
		}
		return event
	})

	// Set the grid as the application root and focus the input field
	app.SetRoot(flex, true).SetFocus(input)

	// Run the application
	log.Info("tview Loading")
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("tview Closed")
}
