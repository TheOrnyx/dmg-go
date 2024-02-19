package emulator

import (
	"fmt"
	"os"
	"time"

	"github.com/TheOrnyx/gameboy-golor/cartridge"
	"github.com/TheOrnyx/gameboy-golor/cpu"
	"github.com/TheOrnyx/gameboy-golor/joypad"
	"github.com/TheOrnyx/gameboy-golor/mmu"
	"github.com/TheOrnyx/gameboy-golor/ppu"
	"github.com/TheOrnyx/gameboy-golor/timer"
	"github.com/TheOrnyx/gameboy-golor/window"
)

var FrameRate float64 = 59.7
var frameDuration = time.Second / time.Duration(FrameRate)


type Emulator struct {
	CPU *cpu.CPU
	MMU *mmu.MMU
	PPU *ppu.PPU
	Timer *timer.Timer
	Renderer window.Screen
	Joypad *joypad.Joypad
	CycleCount int // the cycle count in T-Cycles!
	frameStartTime time.Time
}

// NewEmulator Start a new emulator, load the rom in the given path and return the emulator instance
func NewEmulator(romPath string, renderer window.Screen) (*Emulator, error) {
	emu := new(Emulator)
	rom, err := os.ReadFile(romPath)
	if err != nil {
		return nil, err
	}

	cart, err := cartridge.LoadROM(rom)
	if err != nil {
		return nil, err
	}

	emu.Timer = timer.NewTimer(emu.RequestInterrupt)
	emu.Renderer = renderer
	emu.Joypad = joypad.NewJoypad(emu.RequestInterrupt)
	emu.Joypad.ResetInput()
	emu.PPU = ppu.NewPPU(emu.Timer, emu.RequestInterrupt)
	emu.MMU = mmu.NewMMU(cart, emu.Timer, emu.PPU, emu.Joypad)
	emu.CPU, _ = cpu.NewCPU(emu.MMU, emu.Timer)
	emu.CPU.ResetDebug()
	emu.frameStartTime = time.Now()
	return emu, nil
}

// RequestInterrupt request interrupt on CPU (used for PPU)
func (e *Emulator) RequestInterrupt(code byte)  {
	// fmt.Println("Requesting interrupt:", code)
	e.CPU.RequestInterrupt(code)
}

// RunEmulator run the emulator normally
func (e *Emulator) RunEmulator()  {
	e.frameStartTime = time.Now()
	for  {
		e.Step()
	}

}

// Step step the emulator by one
func (e *Emulator) Step()  {
	mCycles := e.CPU.Step()
	tCycles := mCycles*4
	e.Timer.TickT(tCycles)
	e.PPU.Step(uint16(tCycles))
	e.CycleCount += tCycles

	if float64(e.CycleCount) >= (cpu.ClockSpeed/ FrameRate) { // finish frame
		e.CycleCount = 0
		inputs := e.Renderer.GetInput()
		e.Joypad.HandleInput(inputs)
		e.RenderScreen()
		e.PPU.Screen.Reset()
		
		elapsedTime := time.Since(e.frameStartTime)
		sleepTime := frameDuration - elapsedTime
		time.Sleep(sleepTime)
		
		// fmt.Printf("time since last frame: %v, sleeptime: %v\n", time.Since(e.frameStartTime), sleepTime)
		e.frameStartTime = time.Now()
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
