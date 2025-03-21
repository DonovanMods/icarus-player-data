package character

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/donovanmods/icarus-player-data/lib/shared"
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

type characterJson struct {
	JsonData []string `json:"Characters.json"`
}

/*
// Character struct
*/
type Character struct {
	Name           string                   `json:"CharacterName"`
	Slot           int                      `json:"ChrSlot"`
	XP             uint64                   `json:"XP"`
	XP_Debt        uint64                   `json:"XP_Debt"`
	IsDead         bool                     `json:"IsDead"`
	IsAbandoned    bool                     `json:"IsAbandoned"`
	LastProspectId string                   `json:"LastProspectId"`
	Location       string                   `json:"Location"`
	UnlockedFlags  []int                    `json:"UnlockedFlags"`
	MetaResources  [](shared.MetaResources) `json:"MetaResources"`
	Cosmetic       Cosmetics                `json:"Cosmetic"`
	Talents        [](shared.Talents)       `json:"Talents"`
	TimeLastPlayed uint64                   `json:"TimeLastPlayed"`
}

func (C *Character) Level() int {
	xpTable := shared.BuildExperienceTable()

	return xpTable.Level(C.XP)
}

/*
// characterData struct
*/

type CharacterData struct {
	Characters []Character
	path       string
}

func NewCharacterData(path string) (*CharacterData, error) {
	c := CharacterData{
		Characters: make([]Character, 0, 10),
		path:       path,
	}

	if err := c.Read(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (C *CharacterData) Read() error {
	var characterJson characterJson
	var character Character

	if C.path == "" {
		return errors.New("path is empty")
	}

	file, err := os.Open(C.path)
	if err != nil {
		return fmt.Errorf("CharacterData.Read(): %w", err)
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
		C.Characters = append(C.Characters, character)
	}

	return nil
}

func (C *CharacterData) Write(file io.Writer) error {
	jdata := characterJson{}

	log.Printf("Writing Character data to %q\n", file)

	sort.Slice(C.Characters, func(i, j int) bool {
		return C.Characters[i].Slot < C.Characters[j].Slot
	})

	for _, c := range C.Characters {
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
