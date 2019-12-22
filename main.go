package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

func main() {
	db := NewDb("./MyVideos107.db")
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

	// Initialize application
	app := tview.NewApplication()

	// Create input field
	input := tview.NewInputField().
		SetPlaceholder(" Please Enter...").
		SetPlaceholderTextColor(tcell.ColorYellow)

	list := tview.NewList()

	// Create submit button
	// btn := tview.NewButton("Submit")

	// Create empty Box to pad each side of appGrid
	// bx := tview.NewBox()

	// Create Grid containing the application's widgets
	// appGrid := tview.NewGrid().
	// 	SetColumns(-1, 24, 16, -1).
	// 	SetRows(-1, 2, 3, -1).
	// 	// AddItem(bx, 0, 0, 3, 1, 0, 0, false).
	// 	// AddItem(bx, 0, 1, 1, 1, 0, 0, false).
	// 	// AddItem(bx, 0, 3, 3, 1, 0, 0, false).
	// 	// AddItem(bx, 3, 1, 1, 1, 0, 0, false).
	// 	AddItem(label, 1, 1, 1, 1, 0, 0, false).
	// 	AddItem(input, 1, 2, 1, 1, 0, 0, false).
	// 	AddItem(btn, 2, 1, 1, 2, 0, 0, false)

	// Create search result label
	resultLabel := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Search Results")

	// Create Flex container
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("SEARCH MY MOVIES"), 2, 0, false).AddItem(input, 0, 1, false)

	listFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(resultLabel, 2, 0, false).AddItem(list, 0, 1, false)

	// submittedName is toggled each time Enter is pressed
	var submittedName bool

	// Capture user input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		switch event.Key() {
		case tcell.KeyEnter:
			submittedName = !submittedName

			if submittedName {
				name := input.GetText()
				if strings.TrimSpace(name) == "" {
					name = "Anonymous"
				}

				// Create a modal dialog
				// m := tview.NewModal().
				// 	SetText(fmt.Sprintf("Greetings, %s!", name)).
				// 	AddButtons([]string{"Hello"})

				// Display and focus the dialog

				moviesArray := db.GetMovieWithQuery(name)
				resultLabel.SetText(fmt.Sprintf("Search Results : %d", len(moviesArray)))
				for _, item := range moviesArray {
					list.AddItem(item.C00, "", 0, nil)
				}

				app.SetRoot(listFlex, true).SetFocus(list)
			} else {
				// Clear the input field
				input.SetText("")
				list.Clear()

				// Display appGrid and focus the input field
				app.SetRoot(flex, true).SetFocus(input)
			}

			return nil
		case tcell.KeyEsc:
			//Exit the application
			app.Stop()
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
