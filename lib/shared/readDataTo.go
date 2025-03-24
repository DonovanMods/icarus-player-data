package shared

import (
	"fmt"
	"os"
)

func ReadDataTo(path string, p PlayerData) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer r.Close()

	err = p.ReadF(r)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	return nil
}
