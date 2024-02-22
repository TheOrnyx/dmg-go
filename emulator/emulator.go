package emulator

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/TheOrnyx/dmg-go/cartridge"
	"github.com/TheOrnyx/dmg-go/cpu"
	"github.com/TheOrnyx/dmg-go/joypad"
	"github.com/TheOrnyx/dmg-go/mmu"
	"github.com/TheOrnyx/dmg-go/ppu"
	"github.com/TheOrnyx/dmg-go/timer"
	"github.com/TheOrnyx/dmg-go/window"
)

var FrameRate float64 = 59.7
var frameDuration = time.Second / time.Duration(FrameRate)

var SaveDirLoc string = "./Saves" // TODO - replace this with like a different directory

// generateSaveDirLoc create the save directory location using
// XDG_DATA_HOME or $HOME/.local/share if env doesn't exist
func generateSaveDirLoc()  {
	system := runtime.GOOS
	if system == "linux" {
		if loc, exists := os.LookupEnv("XDG_DATA_HOME"); exists {
			SaveDirLoc = fmt.Sprintf("%v/dmg-go", loc)
		} else {
			SaveDirLoc = fmt.Sprintf("%v/.local/share/dmg-go", os.Getenv("HOME"))
		}
	} else {
		// TODO - implement things for windows and stuff
	}
}

type Emulator struct {
	CPU            *cpu.CPU
	MMU            *mmu.MMU
	PPU            *ppu.PPU
	Timer          *timer.Timer
	Renderer       window.Screen
	Joypad         *joypad.Joypad
	CycleCount     int // the cycle count in T-Cycles!
	frameStartTime time.Time
}

// NewEmulator Start a new emulator, load the rom in the given path and return the emulator instance
func NewEmulator(romPath string, renderer window.Screen) (*Emulator, error) {
	emu := new(Emulator)
	generateSaveDirLoc()
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
	fmt.Println(emu.DebugInfo())
	if true {
		err := emu.LoadSaveFile()
		if err != nil {
			// return nil, err
		}
	}
	return emu, nil
}

// RequestInterrupt request interrupt on CPU (used for PPU)
func (e *Emulator) RequestInterrupt(code byte) {
	e.CPU.RequestInterrupt(code)
}

// RunEmulator run the emulator normally
func (e *Emulator) RunEmulator() {
	e.frameStartTime = time.Now()
	running := true
	for running {
		running = !e.Step()
	}

}

// Step step the emulator by one
// Returns whether or not emu should close
func (e *Emulator) Step() bool {
	closeEmu := false
	mCycles := e.CPU.Step()
	tCycles := mCycles * 4
	e.Timer.TickT(tCycles)
	e.PPU.Step(uint16(tCycles))
	e.CycleCount += tCycles

	if float64(e.CycleCount) >= (cpu.ClockSpeed / FrameRate) { // finish frame
		e.CycleCount = 0
		inputs, close := e.Renderer.GetInput()
		closeEmu = close
		e.Joypad.HandleInput(inputs)
		e.RenderScreen()
		e.PPU.Screen.Reset()

		elapsedTime := time.Since(e.frameStartTime)
		sleepTime := frameDuration - elapsedTime
		time.Sleep(sleepTime)

		// fmt.Printf("time since last frame: %v, sleeptime: %v\n", time.Since(e.frameStartTime), sleepTime)
		e.frameStartTime = time.Now()
	}

	return closeEmu
}

// LoadSaveFile load a save file for a game if exists
func (e *Emulator) LoadSaveFile() error {
	saveLoc := fmt.Sprintf("%s/%s", SaveDirLoc, e.MMU.Cart.SaveTitle())
	if _, err := os.Stat(saveLoc); err != nil {
		return nil
	}

	save, err := os.Open(saveLoc)
	if err != nil {
		return fmt.Errorf("Failed to open save file: %v", err)
	}
	defer save.Close()

	err = e.MMU.Cart.MBC.LoadFile(save)
	if err != nil {
		return fmt.Errorf("Failed to load Save file to RAM: %v", err)
	}

	return nil
}

// CloseEmulator close the emulator and write saves if needed
func (e *Emulator) CloseEmulator() {
	e.Renderer.CloseScreen()
	if e.MMU.Cart.RAMSize == 0 {
		return
	}
	
	os.Mkdir(SaveDirLoc, 0750)
	file, err := os.Create(fmt.Sprintf("%s/%s", SaveDirLoc, e.CPU.MMU.Cart.SaveTitle()))
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	err = e.MMU.Cart.MBC.SaveFile(file)
	if err != nil {
		log.Fatalf("Failed to Save file: %v", err)
	}
}

// RenderScreen render the screen
func (e *Emulator) RenderScreen() {
	e.Renderer.ClearScreen()
	e.Renderer.RenderScreen(&e.PPU.Screen)
}

// DebugInfo print debug info about the emulator
func (e *Emulator) DebugInfo() string {
	return fmt.Sprintf("Cart Info:\n%s\n", e.MMU.Cart)
}
