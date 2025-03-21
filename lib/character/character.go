package character

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/donovanmods/icarus-player-data/lib/shared"
	"github.com/rivo/tview"
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

func (C *Character) nameString() string {
	status := make([]string, 0, 2)
	statusString := ""

	if C.IsDead {
		status = append(status, "[red::bi]DEAD[-::-]")
	}

	if C.IsAbandoned {
		status = append(status, "[purple::bi]Abandoned[-::-]")
	}

	if len(status) > 0 {
		statusString = fmt.Sprintf("(%s)", strings.Join(status, " & "))
	}

	return fmt.Sprintf("[yellow::b]%s[-::-] %s\n\n", C.Name, statusString)
}

func (C *Character) xpString() string {
	return fmt.Sprintf("Level: %-3d (XP: %d%s)\n\n", C.Level(), C.XP, C.xpDebtString())
}

func (C *Character) xpDebtString() string {
	if C.XP_Debt > 0 {
		return fmt.Sprintf("; Debt: %d", C.XP_Debt)
	}

	return ""
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

func (C *CharacterData) Print(index int) tview.Primitive {
	subView := tview.NewTextView()
	subView.SetDynamicColors(true).SetBorderPadding(1, 1, 1, 1)

	if index < 0 || index >= len(C.Characters) {
		fmt.Fprintln(subView, "Invalid Character")
		return subView
	}

	char := &C.Characters[index]

	// Iterate through characters and print each item to the TextView
	fmt.Fprint(subView, char.nameString())
	fmt.Fprint(subView, char.xpString())
	fmt.Fprintf(subView, "Known Talents: %d\n", len(char.Talents))

	return subView
}
