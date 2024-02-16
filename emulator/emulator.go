package emulator

import (
	"fmt"
	"os"

	"github.com/TheOrnyx/gameboy-golor/cartridge"
	"github.com/TheOrnyx/gameboy-golor/cpu"
	"github.com/TheOrnyx/gameboy-golor/mmu"
	"github.com/TheOrnyx/gameboy-golor/ppu"
	"github.com/TheOrnyx/gameboy-golor/timer"
)

const FrameRate = 60

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
	CycleCount int // the cycle count in M-Cycles
}

// NewEmulator Start a new emulator, load the rom in the given path and return the emulator instance
func NewEmulator(romPath string, renderer Screen) (*Emulator, error) {
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
	newEmu.Renderer = renderer

	newEmu.PPU = ppu.NewPPU(newEmu.Timer, newEmu.RequestInterrupt)
	newEmu.MMU = mmu.NewMMU(cart, newEmu.Timer, newEmu.PPU)
	newEmu.CPU, _ = cpu.NewCPU(newEmu.MMU, newEmu.Timer)

	return newEmu, nil
}

// RequestInterrupt request interrupt on CPU (used for PPU)
func (e *Emulator) RequestInterrupt(code byte)  {
	e.CPU.RequestInterrupt(code)
}

// RunEmulator run the emulator normally
func (e *Emulator) RunEmulator()  {
	
}

// Step step the emulator by one
func (e *Emulator) Step()  {
	mCycles := e.CPU.Step()
	tCycles := mCycles*4
	e.Timer.TickM(mCycles)
	e.PPU.Step(uint16(tCycles))
	e.CycleCount += mCycles

	if e.CycleCount >= (cpu.ClockSpeed/ FrameRate) { // finish frame
		e.CycleCount = 0
		e.RenderScreen()
		e.PPU.Screen.Reset()
	}
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
