package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/donovanmods/icarus-character-editor/lib/data"
)

func main() {
	if err := data.Read(); err != nil {
		log.Fatal("Error reading data:", err)
	}

	// box := tview.NewBox().SetBorder(true).SetTitle("[ Icarus Character Editor ]")

	// if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
	// 	panic(err)
	// }

	if false {
		printData() // used for debugging
	}
}

func printData() {
	fmt.Printf("\nProfile Data:\n\n")

	pData, err := json.MarshalIndent(data.ProfileData, "", "\t")
	if err != nil {
		log.Fatal("Error marshaling profile data:", err)
	}

	fmt.Println(string(pData))

	fmt.Printf("\nCharacter Data:\n\n")

	cData, err := json.MarshalIndent(data.CharacterData, "", "\t")
	if err != nil {
		log.Fatal("Error marshaling character data:", err)
	}

	fmt.Println(string(cData))
}
