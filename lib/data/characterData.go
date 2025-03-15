package data

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type Cosmetics struct {
	CustomizationHead           int64 `json:"Customization_Head"`
	CustomizationHair           int64 `json:"Customization_Hair"`
	CustomizationHairColor      int64 `json:"Customization_HairColor"`
	CustomizationBody           int64 `json:"Customization_Body"`
	CustomizationBodyColor      int64 `json:"Customization_BodyColor"`
	CustomizationSkinTone       int64 `json:"Customization_SkinTone"`
	CustomizationHeadTattoo     int64 `json:"Customization_HeadTattoo"`
	CustomizationHeadScar       int64 `json:"Customization_HeadScar"`
	CustomizationHeadFacialHair int64 `json:"Customization_HeadFacialHair"`
	CustomizationCapLogo        int64 `json:"Customization_CapLogo"`
	IsMale                      bool  `json:"IsMale"`
	CustomizationVoice          int64 `json:"Customization_Voice"`
	CustomizationEyeColor       int64 `json:"Customization_EyeColor"`
}

type Character struct {
	CharacterName  string          `json:"CharacterName"`
	ChrSlot        int             `json:"ChrSlot"`
	XP             int64           `json:"XP"`
	XP_Debt        int64           `json:"XP_Debt"`
	IsDead         bool            `json:"IsDead"`
	IsAbandoned    bool            `json:"IsAbandoned"`
	LastProspectId string          `json:"LastProspectId"`
	Location       string          `json:"Location"`
	UnlockedFlags  []int           `json:"UnlockedFlags"`
	MetaResources  []MetaResources `json:"MetaResources"`
	Cosmetic       Cosmetics       `json:"Cosmetic"`
	Talents        []Talents       `json:"Talents"`
	TimeLastPlayed uint64          `json:"TimeLastPlayed"`
}

type charactersJson struct {
	JsonData []string `json:"Characters.json"`
}

var (
	CharacterData []Character
	CharacterPath string
)

func ReadCharacterData(path string) error {
	var characterJson charactersJson
	var character Character

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
		CharacterData = append(CharacterData, character)
	}

	// Save the filepath we read from
	CharacterPath = path

	return nil
}

func WriteCharacterData(path string) error {
	jdata := charactersJson{}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Printf("Writing Character data to %q\n", path)

	sort.Slice(CharacterData, func(i, j int) bool {
		return CharacterData[i].ChrSlot < CharacterData[j].ChrSlot
	})

	for _, c := range CharacterData {
		cdata, err := json.Marshal(c)
		if err != nil {
			return err
		}
		jdata.JsonData = append(jdata.JsonData, string(cdata))
	}

	// Marshal the JSON data into the charactersJson struct
	if err := json.NewEncoder(file).Encode(jdata); err != nil {
		return err
	}

	return nil
}
