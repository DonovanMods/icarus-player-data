package character

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

// Character data
type Character struct {
	Name           string                 `json:"CharacterName"`
	Slot           int                    `json:"ChrSlot"`
	XP             uint64                 `json:"XP"`
	XP_Debt        uint64                 `json:"XP_Debt"`
	IsDead         bool                   `json:"IsDead"`
	IsAbandoned    bool                   `json:"IsAbandoned"`
	LastProspectId string                 `json:"LastProspectId"`
	Location       string                 `json:"Location"`
	UnlockedFlags  []int                  `json:"UnlockedFlags"`
	MetaResources  []shared.MetaResources `json:"MetaResources"`
	Cosmetic       Cosmetics              `json:"Cosmetic"`
	Talents        []shared.Talents       `json:"Talents"`
	TimeLastPlayed uint64                 `json:"TimeLastPlayed"`
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
	Metadata   shared.Metadata
}

// NewCharacterData creates a new CharacterData struct
func NewCharacterData() (*CharacterData, error) {
	c := CharacterData{
		Characters: make([]Character, 0, 16),
		Metadata:   shared.Metadata{FileName: "Characters.json"},
	}

	if err := c.Read(); err != nil {
		return nil, err
	}

	return &c, nil
}

// Read attempts to find and read the CharacterData
func (C *CharacterData) Read() error {
	appDataDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	readData := func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() && filepath.Base(path) == C.Metadata.FileName {
			if err := shared.ReadDataTo(path, C); err != nil {
				return err
			}
			C.Metadata.Path = path
		}

		return nil
	}

	playerData := filepath.Join(appDataDir, "Icarus", "Saved", "PlayerData")

	if err := filepath.WalkDir(playerData, readData); err != nil {
		return err
	}

	return nil
}

// Read reads the CharacterData from an io.Reader
func (C *CharacterData) ReadF(file io.Reader) error {
	var characterJson characterJson
	var character Character

	if file == nil {
		return errors.New("CharacterData.Read(): input is nil - expected an io.Reader")
	}

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

// Write writes the ProfileData to a file
func (C *CharacterData) Write() error {
	if C.Metadata.Path == "" {
		return errors.New("CharacterData.Write(): Metadata.Path is empty")
	}

	file, err := os.Create(C.Metadata.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	return C.WriteF(file)
}

// Write writes the CharacterData to an io.Writer
func (C *CharacterData) WriteF(file io.Writer) error {
	jdata := characterJson{}

	if file == nil {
		return errors.New("CharacterData.Write(): input is nil - expected an io.WriteCloser")
	}

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
