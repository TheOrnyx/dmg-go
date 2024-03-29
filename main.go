package main

import (
	"fmt"
	"log"
	"os"

	"flag"

	"github.com/TheOrnyx/dmg-go/debugger"
	_ "github.com/TheOrnyx/dmg-go/debugger"
	emu "github.com/TheOrnyx/dmg-go/emulator"
	"github.com/TheOrnyx/dmg-go/window"
	// "github.com/TheOrnyx/dmg-go/window"
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of dmg-go: %s [flags] [rom path]", os.Args[0])

		flag.PrintDefaults()
	}
	
	flag.BoolFunc("debug", "Use debug mode", enableDebug)
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

	defer emulator.CloseEmulator()
	
	if debugMode {
		debugger.DebugEmu(emulator)
	} else {
		emulator.RunEmulator()
	}
}
