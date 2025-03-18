package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	Credits    = "Credits"
	Refund     = "Refund"
	Exotics    = "Exotic1"
	ExoticsRed = "Exotic2"
)

type metaResources struct {
	MetaRow string `json:"MetaRow"`
	Count   int    `json:"Count"`
}

type talents struct {
	RowName string `json:"RowName"`
	Rank    int    `json:"Rank"`
}

type profile struct {
	UserID        string          `json:"UserID"`
	MetaResources []metaResources `json:"MetaResources"`
	UnlockedFlags []int           `json:"UnlockedFlags"`
	Talents       []talents       `json:"Talents"`
}

type profileData struct {
	Profile profile
	Path    string
}

func createProfileData(path string) (*profileData, error) {
	p := profileData{
		Profile: profile{},
		Path:    path,
	}

	if err := p.Read(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (C *profileData) Print() tview.Primitive {
	table := tview.NewTable().SetBorders(false)
	table.SetBorderPadding(1, 1, 1, 1)

	table.SetCell(0, 0, tview.NewTableCell("UserID:").SetTextColor(tcell.ColorGreen))
	table.SetCell(0, 1, tview.NewTableCell(C.Profile.UserID).SetTextColor(tcell.ColorWhite))

	table.SetCell(2, 0, tview.NewTableCell("Credits:").SetTextColor(tcell.ColorGreen))
	table.SetCell(2, 1, tview.NewTableCell(C.getMetaCountFor(Credits)).SetTextColor(tcell.ColorYellow))
	table.SetCell(2, 2, tview.NewTableCell("Refund:").SetTextColor(tcell.ColorGreen))
	table.SetCell(2, 3, tview.NewTableCell(C.getMetaCountFor(Refund)).SetTextColor(tcell.ColorYellow))

	table.SetCell(4, 1, tview.NewTableCell("Purple").SetTextColor(tcell.ColorPurple).SetAlign(tview.AlignRight))
	table.SetCell(4, 2, tview.NewTableCell("Red").SetTextColor(tcell.ColorRed).SetAlign(tview.AlignRight))

	table.SetCell(5, 0, tview.NewTableCell("Exotics:").SetTextColor(tcell.ColorBlue))
	table.SetCell(5, 1, tview.NewTableCell(C.getMetaCountFor(Exotics)).SetTextColor(tcell.ColorPurple).SetAlign(tview.AlignRight))
	table.SetCell(5, 2, tview.NewTableCell(C.getMetaCountFor(ExoticsRed)).SetTextColor(tcell.ColorRed).SetAlign(tview.AlignRight))

	return table
}

func (P *profileData) Read() error {
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

func (P *profileData) Write(file io.Writer) error {
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

func (P *profileData) metaMap() map[string]int {
	m := make(map[string]int)

	for _, meta := range P.Profile.MetaResources {
		m[meta.MetaRow] = meta.Count
	}

	return m
}

func (P *profileData) getMetaCountFor(key string) string {
	return strconv.Itoa(P.metaMap()[key])
}
