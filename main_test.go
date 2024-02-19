package main

import (
	"flag"
	"testing"
	"time"

	"github.com/TheOrnyx/gameboy-golor/emulator"
	"github.com/TheOrnyx/gameboy-golor/window"
)

// TestEmu test the emulator and benchmark it
func TestEmu(t *testing.T)  {
	flag.Parse()
	romPath := flag.Args()[0]
	win := window.InitSDLWindowSystem(WinWidth, WinHeight, WinScalar)

	emulator, err := emulator.NewEmulator(romPath, win)
	if err != nil {
		t.Fatalf("Failed to create emulator:%v", err)
	}

	defer emulator.Renderer.CloseScreen()
	timeNow := time.Now()
	
	for time.Since(timeNow) < time.Second * 30 {
		emulator.Step()
	}
}
