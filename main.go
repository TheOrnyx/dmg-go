package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TheOrnyx/gameboy-golor/cartridge"
	"github.com/TheOrnyx/gameboy-golor/cpu"
	"github.com/TheOrnyx/gameboy-golor/mmu"
	// "github.com/TheOrnyx/gameboy-golor/window"
)

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

// PrintDebugInfo print debug info about the emulator
func (e *Emulator) PrintDebugInfo()  {
	DebugLog.Println("Rom name:", e.MMU.Cart.Title)
}

const UsingSDL = true
const WinScalar = 4 //the scalar used to scale up the gbc screen
const WinWidth, WinHeight = 160 * WinScalar, 144 * WinScalar

var InfoLog = log.New(os.Stdout, "[INFO] ", log.Ldate)
var DebugLog = log.New(os.Stdout, "[DEBUG] ", log.Ldate)
var WarnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
var FatalLog = log.New(os.Stdout, "[FaTAL] ", log.LstdFlags)

func main() {
	InfoLog.Println("Starting...")
	// window.StartSDLWindowSystem(WinWidth, WinHeight)
	emulator, err := NewEmulator("./Data/Roms/blargg-test-roms/cpu_instrs/cpu_instrs.gb")
	if err != nil {
		FatalLog.Fatal("Error making new emulator:", err)
	}
	emulator.CPU.ResetDebug()
	emulator.PrintDebugInfo()

	var running bool = true
	for running {
		emulator.CPU.Step()
		DebugLog.Println(emulator.CPU)
		if emulator.MMU.ReadByte(0xff02) != 0 {
			char := emulator.MMU.ReadByte(0xff01)
			fmt.Println("Found write")
			fmt.Printf("%v", char)
			// emulator.MMU.WriteByte(0xff02, 0)
		}
		time.Sleep(time.Second/2)
	}
	
	InfoLog.Println("Program finished, exiting...")
}
