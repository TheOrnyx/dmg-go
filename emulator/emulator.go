package emulator

import (
	"fmt"
	"log"
	"os"

	"github.com/TheOrnyx/gameboy-golor/cartridge"
	"github.com/TheOrnyx/gameboy-golor/cpu"
	"github.com/TheOrnyx/gameboy-golor/mmu"
)

var InfoLog = log.New(os.Stdout, "[INFO] ", log.Ldate)
var DebugLog = log.New(os.Stdout, "[DEBUG] ", log.Ldate)
var WarnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
var FatalLog = log.New(os.Stdout, "[FaTAL] ", log.LstdFlags)

type Emulator struct {
	CPU *cpu.CPU
	MMU *mmu.MMU
}

// NewEmulator Start a new emulator, load the rom in the given path and return the emulator instance
func NewEmulator(romPath string) (*Emulator, error) {
	newEmu := new(Emulator)
	rom, err := os.ReadFile(romPath)
	if err != nil {
		return nil, err
	}

	cart, err := cartridge.LoadROM(rom)
	if err != nil {
		return nil, err
	}
	newEmu.MMU = mmu.NewMMU(cart)
	
	newEmu.CPU, _ = cpu.NewCPU(newEmu.MMU)

	return newEmu, nil
}

// DebugInfo print debug info about the emulator
func (e *Emulator) DebugInfo() string {
	return fmt.Sprintf("Cart Info:\n%s\n", e.MMU.Cart)
}
