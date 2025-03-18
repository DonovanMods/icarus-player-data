package main

import (
	"fmt"
	"log"

	"github.com/donovanmods/icarus-character-editor/lib/data"
	"github.com/rivo/tview"
)

func main() {
	startApp()
}

func startApp() {
	if err := data.Read(); err != nil {
		log.Fatal("error reading player data:", err)
	}

	// Create a new TUI application
	app := tview.NewApplication()

	// Create a TextView that will display the character list in the TUI
	mainMenu := tview.NewList().SetHighlightFullLine(true).SetWrapAround(false)
	mainMenu.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("[ Characters ]")
	mainMenu.SetCurrentItem(0) // Set the first item as the current item

	dataView := tview.NewFlex()
	dataView.SetBorder(true).SetBorderPadding(0, 0, 0, 0)

	// Iterate through characters and add each item to the character list
	for i, item := range data.CharacterData.Characters {
		mainMenu.AddItem(item.Name, "", rune(i+49), nil)
	}

	// Add a profile option
	mainMenu.AddItem("Player Profile", "", 'p', nil)

	// Add a quit option
	mainMenu.AddItem("Exit the Program", "", 'q', func() {
		app.Stop()
	})

	// Set the function to be called when a character is selected
	mainMenu.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if shortcut == rune('q') {
			dataView.Clear().SetTitle("[ Quit ]")
			dataView.AddItem(quitView(), 0, 1, false)

			return
		}

		if shortcut == rune('p') {
			dataView.Clear().SetTitle("[ Player Profile ]")
			dataView.AddItem(data.ProfileData.Print(), 0, 1, false)
			return
		}

		// Print the selected character data
		dataView.Clear().SetTitle("[ Character Data ]")
		dataView.AddItem(data.CharacterData.Print(index), 0, 1, false)
	})

	// Print the first character data by default
	dataView.Clear().SetTitle("[ Character Data ]")
	dataView.AddItem(data.CharacterData.Print(0), 0, 1, false)

	// Create a layout using Flex to display the character list and the form side by side
	flex := tview.NewFlex().
		AddItem(mainMenu, 0, 1, true). // Left side: character list
		AddItem(dataView, 0, 4, false) // Right side: character data

	// Start the TUI application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func quitView() tview.Primitive {
	view := tview.NewTextView()
	view.
		SetDynamicColors(true).
		SetBorderPadding(1, 1, 1, 1)

	fmt.Fprintln(view, "[green]Exit the Character Editor without Saving[-]")

	return view
}
