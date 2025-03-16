package main

import (
	"log"

	"github.com/donovanmods/icarus-character-editor/lib/data"
	"github.com/rivo/tview"
)

func main() {
	if err := data.Read(); err != nil {
		log.Fatal("error reading player data:", err)
	}

	// Create a new TUI application
	app := tview.NewApplication()

	characterDataView := tview.NewTextView().
		SetDynamicColors(true). // Enable dynamic coloring of text
		SetRegions(true).       // Allows regions for interaction (not used here)
		SetWordWrap(true)       // Enables word wrapping to fit the TextView size
	characterDataView.SetBorder(true).SetBorderPadding(1, 1, 1, 1).SetTitle("[ Character Data ]")

	// Create a TextView that will display the character list in the TUI
	characterListView := tview.NewList().SetHighlightFullLine(true).SetWrapAround(false)
	characterListView.SetBorder(true).SetBorderPadding(1, 1, 1, 1).SetTitle("[ Characters ]")

	// Iterate through characters and add each item to the character list
	for i, item := range data.CharacterData.Characters {
		characterListView.AddItem(item.Name, "", rune(i+49), nil)
	}

	// Add a quit option
	characterListView.AddItem("Quit", "Exit the Program", 'q', func() {
		app.Stop()
	})

	// Set the function to be called when a character is selected
	characterListView.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		data.CharacterData.Print(characterDataView, index, shortcut)
	})

	characterListView.SetCurrentItem(0)                 // Set the first item as the current item
	data.CharacterData.Print(characterDataView, 0, 'a') // Print the first character data

	// Create a layout using Flex to display the character list and the form side by side
	flex := tview.NewFlex().
		AddItem(characterListView, 0, 1, true). // Left side: character list
		AddItem(characterDataView, 0, 4, false) // Right side: character data

	// Start the TUI application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
