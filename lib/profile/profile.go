package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/donovanmods/icarus-player-data/lib/shared"
	"github.com/rivo/tview"
)

const (
	Credits       = "Credits"
	Refund        = "Refund"
	PurpleExotics = "Exotic1"
	RedExotics    = "Exotic2"
)

type profile struct {
	UserID        string                 `json:"UserID"`
	MetaResources []shared.MetaResources `json:"MetaResources"`
	UnlockedFlags []int                  `json:"UnlockedFlags"`
	Talents       []shared.Talents       `json:"Talents"`
}

type ProfileData struct {
	Profile profile
	Path    string
	Dirty   bool
}

func NewProfileData(path string) (*ProfileData, error) {
	p := ProfileData{
		Profile: profile{},
		Path:    path,
		Dirty:   false,
	}

	if err := p.Read(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (P *ProfileData) Print() tview.Primitive {
	saveCount := func(field string, text string) {
		if text == "" {
			return
		}

		count, err := strconv.Atoi(text)
		if err != nil {
			log.Print(fmt.Errorf("unable to convert %s to int: %w", text, err))
		}

		P.setCountFor(field, count)
		P.Dirty = true
	}

	form := tview.NewForm()
	form.SetBorder(false).SetBorderPadding(1, 1, 1, 1)

	form.AddTextView("UserID", P.Profile.UserID, 40, 2, true, false)
	form.AddInputField("Credits", P.getCountFor(Credits), 10, nil, func(text string) {
		saveCount(Credits, text)
	})
	form.AddInputField("Refund", P.getCountFor(Refund), 10, nil, func(text string) {
		saveCount(Refund, text)
	})
	form.AddInputField("Purple Exotics", P.getCountFor(PurpleExotics), 10, nil, func(text string) {
		saveCount(PurpleExotics, text)
	})
	form.AddInputField("Red Exotics", P.getCountFor(RedExotics), 10, nil, func(text string) {
		saveCount(RedExotics, text)
	})

	// table := tview.NewTable().SetSelectable(true, true).SetBorders(false)
	// table.SetBorderPadding(1, 1, 1, 1)

	// table.SetCell(0, 0, tview.NewTableCell("UserID:").SetTextColor(tcell.ColorGreen).SetSelectable(false))
	// table.SetCell(0, 1, tview.NewTableCell(C.Profile.UserID).SetTextColor(tcell.ColorWhite).SetSelectable(false))

	// table.SetCell(2, 0, tview.NewTableCell("Credits:").SetTextColor(tcell.ColorGreen).SetSelectable(false))
	// table.SetCell(2, 1, tview.NewTableCell(C.getMetaCountFor(Credits)).SetTextColor(tcell.ColorYellow).SetSelectable(true))
	// table.SetCell(2, 2, tview.NewTableCell("Refund:").SetTextColor(tcell.ColorGreen).SetSelectable(false))
	// table.SetCell(2, 3, tview.NewTableCell(C.getMetaCountFor(Refund)).SetTextColor(tcell.ColorYellow).SetSelectable(true))

	// table.SetCell(4, 1, tview.NewTableCell("Purple").SetTextColor(tcell.ColorPurple).SetAlign(tview.AlignRight).SetSelectable(false))
	// table.SetCell(4, 2, tview.NewTableCell("Red").SetTextColor(tcell.ColorRed).SetAlign(tview.AlignRight).SetSelectable(false))

	// table.SetCell(5, 0, tview.NewTableCell("Exotics:").SetTextColor(tcell.ColorBlue).SetSelectable(false))
	// table.SetCell(5, 1, tview.NewTableCell(C.getMetaCountFor(Exotics)).SetTextColor(tcell.ColorPurple).SetAlign(tview.AlignRight).SetSelectable(true))
	// table.SetCell(5, 2, tview.NewTableCell(C.getMetaCountFor(ExoticsRed)).SetTextColor(tcell.ColorRed).SetAlign(tview.AlignRight).SetSelectable(true))

	return form
}

func (P *ProfileData) Read() error {
	if P.Path == "" {
		return errors.New("path is empty")
	}

	file, err := os.Open(P.Path)
	if err != nil {
		return fmt.Errorf("ProfileData.Read(): %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&P.Profile); err != nil {
		return err
	}

	return nil
}

func (P *ProfileData) Write(file io.Writer) error {
	if !P.Dirty {
		return nil
	}

	log.Printf("Writing Profile data to %q\n", file)

	jdata, err := json.Marshal(P.Profile)
	if err != nil {
		return err
	}

	_, err = file.Write(jdata)
	if err != nil {
		return err
	}

	return nil
}

func (P *ProfileData) metaMap() map[string]int {
	m := make(map[string]int)

	for _, meta := range P.Profile.MetaResources {
		m[meta.MetaRow] = meta.Count
	}

	return m
}

func (P *ProfileData) getCountFor(key string) string {
	return strconv.Itoa(P.metaMap()[key])
}

func (P *ProfileData) setCountFor(key string, count int) {
	if _, ok := P.metaMap()[key]; !ok {
		P.Profile.MetaResources = append(P.Profile.MetaResources, shared.MetaResources{MetaRow: key, Count: count})
		return
	}

	for i, meta := range P.Profile.MetaResources {
		if meta.MetaRow == key {
			P.Profile.MetaResources[i].Count = count
		}
	}
}
