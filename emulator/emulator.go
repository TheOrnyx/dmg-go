package emulator

import (
	"fmt"
	"log"
	"os"

	"github.com/TheOrnyx/gameboy-golor/cartridge"
	"github.com/TheOrnyx/gameboy-golor/cpu"
	"github.com/TheOrnyx/gameboy-golor/mmu"
	"github.com/TheOrnyx/gameboy-golor/ppu"
	"github.com/TheOrnyx/gameboy-golor/timer"
	"github.com/TheOrnyx/gameboy-golor/window"
)

var InfoLog = log.New(os.Stdout, "[INFO] ", log.Ldate)
var DebugLog = log.New(os.Stdout, "[DEBUG] ", log.Ldate)
var WarnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
var FatalLog = log.New(os.Stdout, "[FaTAL] ", log.LstdFlags)

type Screen interface {
	ClearScreen()
	RenderScreen(screen *ppu.Screen)
	CloseScreen()
}

type Emulator struct {
	CPU *cpu.CPU
	MMU *mmu.MMU
	PPU *ppu.PPU
	Timer *timer.Timer
	Renderer Screen
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

	newEmu.Timer = new(timer.Timer)

	newEmu.PPU = ppu.NewPPU(newEmu.Timer)
	newEmu.MMU = mmu.NewMMU(cart, newEmu.Timer, newEmu.PPU)
	newEmu.CPU, _ = cpu.NewCPU(newEmu.MMU, newEmu.Timer)

	if false {
		newEmu.Renderer = window.InitSDLWindowSystem(320, 288)
	}

	return newEmu, nil
}

// RunEmulator run the emulator normally
func (e *Emulator) RunEmulator()  {
	
}

// Step step the emulator by one
func (e *Emulator) Step()  {
	e.CPU.Step()
}

// RenderScreen render the screen
func (e *Emulator) RenderScreen()  {
	e.Renderer.ClearScreen()
	e.Renderer.RenderScreen(&e.PPU.Screen)
}

// DebugInfo print debug info about the emulator
func (e *Emulator) DebugInfo() string {
	return fmt.Sprintf("Cart Info:\n%s\n", e.MMU.Cart)
}
