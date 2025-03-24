package profile

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/donovanmods/icarus-player-data/lib/shared"
)

const (
	Credits       = "Credits"
	Refund        = "Refund"
	PurpleExotics = "Exotic1"
	RedExotics    = "Exotic2"
)

type ProfileData struct {
	UserID        string                 `json:"UserID"`
	MetaResources []shared.MetaResources `json:"MetaResources"`
	UnlockedFlags []int                  `json:"UnlockedFlags"`
	Talents       []shared.Talents       `json:"Talents"`
}

// NewProfileData creates a new ProfileData struct
func NewProfileData(r io.Reader) (*ProfileData, error) {
	p := ProfileData{}

	if err := p.Read(r); err != nil {
		return nil, err
	}

	return &p, nil
}

// Read reads the ProfileData from an io.Reader
func (P *ProfileData) Read(file io.Reader) error {
	if file == nil {
		return errors.New("ProfileData.Read(): input is nil - expected an io.Reader")
	}

	if err := json.NewDecoder(file).Decode(&P); err != nil {
		return err
	}

	return nil
}

// Write writes the ProfileData to an io.Writer
// The function will only write if the data has been altered
func (P *ProfileData) Write(file io.Writer) error {
	if file == nil {
		return errors.New("ProfileData.Write(): input is nil - expected an io.Writer")
	}

	log.Printf("Writing Profile data to %q\n", file)

	jdata, err := json.Marshal(P)
	if err != nil {
		return err
	}

	_, err = file.Write(jdata)
	if err != nil {
		return err
	}

	return nil
}

// GetCountFor returns the count for a given key
func (P *ProfileData) GetCountFor(key string) string {
	return strconv.Itoa(P.metaMap()[key])
}

// SetCountFor sets the count for a given key
func (P *ProfileData) SetCountFor(key string, count int) {
	if _, ok := P.metaMap()[key]; !ok {
		P.MetaResources = append(P.MetaResources, shared.MetaResources{MetaRow: key, Count: count})
		return
	}

	for i, meta := range P.MetaResources {
		if meta.MetaRow == key {
			P.MetaResources[i].Count = count
		}
	}
}

func (P *ProfileData) metaMap() map[string]int {
	m := make(map[string]int)

	for _, meta := range P.MetaResources {
		m[meta.MetaRow] = meta.Count
	}

	return m
}
