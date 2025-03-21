package data

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/donovanmods/icarus-player-data/character"
	"github.com/donovanmods/icarus-player-data/profile"
)

var (
	CharacterData *character.CharacterData
	ProfileData   *profile.ProfileData
)

func Read() error {
	readCount := 0

	appDataDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	readData := func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() {
			if filepath.Base(path) == "Characters.json" {
				log.Printf("Reading Character data from %q\n", path)

				CharacterData, err = character.NewCharacterData(path)
				if err != nil {
					return fmt.Errorf("reading Character data: %w", err)
				}

				readCount++
			}

			if filepath.Base(path) == "Profile.json" {
				log.Printf("Reading Profile data from %q", path)

				ProfileData, err = profile.NewProfileData(path)
				if err != nil {
					return fmt.Errorf("reading Profile data: %w", err)
				}

				readCount++
			}
		}

		// Stop walking after we've read both files
		if readCount == 2 {
			log.Println("Finished reading data")
			return filepath.SkipAll
		}

		return nil
	}

	playerData := filepath.Join(appDataDir, "Icarus", "Saved", "PlayerData")

	log.Println("Reading Data from", playerData)

	if err := filepath.WalkDir(playerData, readData); err != nil {
		return err
	}

	return nil
}
