package data

import (
	"encoding/json"
	"os"
)

type MetaResources struct {
	MetaRow string `json:"MetaRow"`
	Count   int    `json:"Count"`
}

type Talents struct {
	RowName string `json:"RowName"`
	Rank    int    `json:"Rank"`
}

type Profile struct {
	UserID        string          `json:"UserID"`
	MetaResources []MetaResources `json:"MetaResources"`
	UnlockedFlags []int           `json:"UnlockedFlags"`
	Talents       []Talents       `json:"Talents"`
}

var (
	ProfileData Profile
	ProfilePath string
)

func ReadProfileData(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&ProfileData); err != nil {
		return err
	}

	// Save the filepath we read from
	ProfilePath = path

	return nil
}
