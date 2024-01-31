package main

import (
	"log"

	"github.com/TheOrnyx/gameboy-golor/debugger"
	emu "github.com/TheOrnyx/gameboy-golor/emulator"
	// "github.com/TheOrnyx/gameboy-golor/window"
)

const UsingSDL = true
const WinScalar = 4 //the scalar used to scale up the gbc screen
const WinWidth, WinHeight = 160 * WinScalar, 144 * WinScalar

func main() {
	// fmt.Println("Starting...")
	// window.StartSDLWindowSystem(WinWidth, WinHeight)
	emulator, err := emu.NewEmulator("./Data/Roms/blargg-test-roms/cpu_instrs/cpu_instrs.gb")
	if err != nil {
		log.Fatal("Error making new emulator:", err)
	}

	debugger.DebugEmulatorDoctor(emulator)
	// debugger.DebugEmulator(emulator)
	// var running bool = true
	// for running {
	// 	emulator.CPU.Step()
	// 	emu.DebugLog.Println(emulator.CPU)
	// 	if emulator.MMU.ReadByte(0xff02) != 0 {
	// 		char := emulator.MMU.ReadByte(0xff01)
	// 		fmt.Println("Found write")
	// 		fmt.Printf("%v", char)
	// 		// emulator.MMU.WriteByte(0xff02, 0)
	// 	}
	// 	time.Sleep(time.Second/2)
	// }
	
	// emu.InfoLog.Println("Program finished, exiting...")
}
