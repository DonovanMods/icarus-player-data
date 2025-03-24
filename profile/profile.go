package profile

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	Metadata      shared.Metadata        `json:"-"`
}

// NewProfileData creates a new ProfileData struct
func NewProfileData() (*ProfileData, error) {
	p := ProfileData{
		MetaResources: make([]shared.MetaResources, 0, 128),
		UnlockedFlags: make([]int, 0),
		Talents:       make([]shared.Talents, 0, 128),
		Metadata:      shared.Metadata{FileName: "Profile.json"},
	}

	if err := p.Read(); err != nil {
		return nil, err
	}

	return &p, nil
}

// OpenProfileFile return an io.Reader for ProfileData
func (P *ProfileData) Read() error {
	appDataDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	playerData := filepath.Join(appDataDir, "Icarus", "Saved", "PlayerData")

	readData := func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() && filepath.Base(path) == P.Metadata.FileName {
			if err := shared.ReadDataTo(path, P); err != nil {
				return err
			}
			P.Metadata.Path = path
		}

		return nil
	}

	if err := filepath.WalkDir(playerData, readData); err != nil {
		return err
	}

	return nil
}

// Read reads the ProfileData from an io.Reader
func (P *ProfileData) ReadF(file io.Reader) error {
	if file == nil {
		return errors.New("ProfileData.Read(): input is nil - expected an io.Reader")
	}

	if err := json.NewDecoder(file).Decode(&P); err != nil {
		return err
	}

	return nil
}

// Write writes the ProfileData to a file
func (P *ProfileData) Write() error {
	if P.Metadata.Path == "" {
		return errors.New("ProfileData.Write(): Metadata.Path is empty")
	}

	file, err := os.Create(P.Metadata.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	return P.WriteF(file)
}

// Write writes the ProfileData to an io.Writer
// The function will only write if the data has been altered
func (P *ProfileData) WriteF(file io.Writer) error {
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
