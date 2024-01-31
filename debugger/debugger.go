package debugger

import (
	"fmt"
	"strings"
	"time"

	"github.com/TheOrnyx/gameboy-golor/emulator"
	"github.com/gdamore/tcell/v2"
)

type Debugger struct {
	Emu *emulator.Emulator
}

// DebugEmulatorDoctor run and debug emulator outputting in doctor format
func DebugEmulatorDoctor(emu *emulator.Emulator)  {
	emu.CPU.ResetDebug()
	count := 1
	maxTests := 55000000 // set to 0 or below for infinite running >:3
	debugger := Debugger{emu}
	var serialOutput string
	
	for !strings.Contains(serialOutput, "Passed") && count != maxTests {
		fmt.Printf("%v\n", emu.CPU.StringDoctor())
		emu.CPU.Step()
		serialWritten, data := debugger.checkSerialLink()
		if serialWritten {
			// fmt.Printf("%v", string(data))
			serialOutput += string(data)
		}
		
		// time.Sleep(time.Nanosecond*2)

		count += 1
	}
	fmt.Println(serialOutput)
}

// DebugEmulator run and debug an emulator
func DebugEmulator(emu *emulator.Emulator)  {
	emu.CPU.ResetDebug()
	count := 1
	// fmt.Println(emu.MMU.Cart.ROM)
	fmt.Println("Beginning debug...")
	fmt.Println(emu.DebugInfo())
	running := true
	debugger := Debugger{emu}
	var serialOutput string
	
	for running && !strings.Contains(serialOutput, "Passed") {
		fmt.Println("---------------------------------------------")
		fmt.Println("Current Step Count:", count)
		serialWritten, data := debugger.checkSerialLink()
		if serialWritten {
			serialOutput += string(data)
		}
		fmt.Printf("%v\n", emu.CPU)
		emu.CPU.Step()
		fmt.Scanln()
		
		if true {
			time.Sleep(time.Nanosecond*5)
		} else {
			fmt.Scanln()
		}
		count += 1
	}
	// fmt.Println("\n",serialOutput)
}

// printDoctorString print the string in gameboy doctor form
func printDoctorString(emu *emulator.Emulator)  {
	
}

// checkSerialLink check whether data has been sent to the serial link and return it and the status
func (d *Debugger) checkSerialLink() (bool, byte) {
	if d.Emu.MMU.ReadByte(0xFF02) != 0 {
		char := d.Emu.MMU.ReadByte(0xFF01)
		d.Emu.MMU.WriteByte(0xFF02, 0)
		return true, char
		// fmt.Println("Found write")
		// fmt.Printf("%v", char)
	}
	return false, 0
}

// drawText draws given text into specified position on screen
func drawText(s *tcell.Screen)  {
	
}
