package main

import (
	"fmt"
	"os"
)

type Cartridge struct {
	ROM []byte
}



// LoadROM load and initialize a ROM based on cart path
func LoadROM(path string) (*Cartridge, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load rom file: %s", err)
	}

	return &Cartridge{ROM: buffer}, nil
}
