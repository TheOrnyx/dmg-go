package main

import (
	"log"

	"flag"

	"github.com/TheOrnyx/gameboy-golor/debugger"
	_ "github.com/TheOrnyx/gameboy-golor/debugger"
	emu "github.com/TheOrnyx/gameboy-golor/emulator"
	"github.com/TheOrnyx/gameboy-golor/window"
	// "github.com/TheOrnyx/gameboy-golor/window"
)

var debugMode bool = false
const UsingSDL = true
const WinScalar = 4 //the scalar used to scale up the gbc screen
const WinWidth, WinHeight = 160 * WinScalar, 144 * WinScalar

// enableDebug just for the flag to use
func enableDebug(b string) error {
	debugMode = true
	return nil
}

func main() {
	flag.BoolFunc("debug", "Activate Debug Mode", enableDebug)
	flag.Parse()
	romPath := flag.Args()[0]
	var win window.Screen
	
	if debugMode {
		win = window.CreateDebugWindow(WinWidth, WinHeight, WinScalar)
	} else {
		win = window.InitSDLWindowSystem(WinWidth, WinHeight, WinScalar)
	} 
	
	emulator, err := emu.NewEmulator(romPath, win)
	if err != nil {
		log.Fatal("Error making new emulator:", err)
	}

	defer emulator.Renderer.CloseScreen()

	if debugMode {
		debugger.DebugEmu(emulator)
	} else {
		emulator.RunEmulator()
		// for  {
		// 	emulator.Step()
		// }
	}
}
