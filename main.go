package main

import (
	"fmt"

	data "github.com/donovanmods/icarus-character-editor/lib"
)

func main() {
	// box := tview.NewBox().SetBorder(true).SetTitle("[ Icarus Character Editor ]")

	// if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
	// 	panic(err)
	// }

	if err := data.Read(); err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	fmt.Printf("\nProfile Data:\n\n")

	// Print the profile data
	fmt.Printf("User ID: %s\n", data.ProfileData.UserID)
	fmt.Println("Meta Resources:")
	for _, resource := range data.ProfileData.MetaResources {
		fmt.Printf("  %s: %d\n", resource.MetaRow, resource.Count)
	}
	fmt.Println("Unlocked Flags:", data.ProfileData.UnlockedFlags)

	// Print the talents
	fmt.Println("Talents:")
	for _, talent := range data.ProfileData.Talents {
		fmt.Printf("  %s: %d\n", talent.RowName, talent.Rank)
	}

	// Print the character data
	fmt.Printf("\nCharacter Data:\n\n")

	// Print the character data
	for _, character := range data.Characters {
		fmt.Printf("Character Slot: %d\n", character.ChrSlot)
		fmt.Printf("Character Name: %s\n", character.CharacterName)
		fmt.Printf("XP: %d\n", character.XP)
		fmt.Printf("XP Debt: %d\n", character.XP_Debt)
		fmt.Printf("Is Dead: %t\n", character.IsDead)
		fmt.Printf("Is Abandoned: %t\n", character.IsAbandoned)
		fmt.Printf("Last Prospect ID: %s\n", character.LastProspectId)
		fmt.Println()
	}

}
