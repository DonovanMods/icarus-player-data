package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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
