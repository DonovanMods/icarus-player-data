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

func (P *ProfileData) GetCountFor(key string) string {
	return strconv.Itoa(P.metaMap()[key])
}

func (P *ProfileData) SetCountFor(key string, count int) {
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

func (P *ProfileData) metaMap() map[string]int {
	m := make(map[string]int)

	for _, meta := range P.Profile.MetaResources {
		m[meta.MetaRow] = meta.Count
	}

	return m
}
