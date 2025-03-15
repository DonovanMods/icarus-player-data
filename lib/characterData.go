package data

import (
	"encoding/json"
	"os"
)

type CharacterData struct {
	ChrSlot        int    `json:"ChrSlot"`
	CharacterName  string `json:"CharacterName"`
	XP             int64  `json:"XP"`
	XP_Debt        int64  `json:"XP_Debt"`
	IsDead         bool   `json:"IsDead"`
	IsAbandoned    bool   `json:"IsAbandoned"`
	LastProspectId string `json:"LastProspectId"`
}

type charactersJson struct {
	JsonData []string `json:"Characters.json"`
}

var Characters []CharacterData

func ReadCharacterData(path string) error {
	var characterJson charactersJson
	var character CharacterData

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file and unmarshal the JSON data into the Characters slice
	if err := json.NewDecoder(file).Decode(&characterJson); err != nil {
		return err
	}

	// Unmarshal the JSON data into the Characters slice
	for _, c := range characterJson.JsonData {
		if err := json.Unmarshal([]byte(c), &character); err != nil {
			return err
		}
		Characters = append(Characters, character)
	}

	return nil
}
