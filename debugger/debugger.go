package debugger

import (
	"fmt"

	"github.com/TheOrnyx/gameboy-golor/emulator"
	"github.com/gdamore/tcell/v2"
)

type Debugger struct {
	Emu *emulator.Emulator
}

// DebugEmulator run and debug an emulator
func DebugEmulator(emu *emulator.Emulator)  {
	emu.CPU.ResetDebug()
	// fmt.Println(emu.MMU.Cart.ROM)
	fmt.Println("Beginning debug...")
	fmt.Println(emu.DebugInfo())
	running := true
	// debugger := Debugger{emu}
	
	for running {
		fmt.Println("---------------------------------------------")
		fmt.Println("")
		// serialWritten, data := debugger.checkSerialLink()
		// if serialWritten {
		// 	fmt.Println("Found serial data write")
		// 	fmt.Printf("%v", data)
		// }
		emu.CPU.Step()
		fmt.Printf("%v\n", emu.CPU)
		fmt.Scanln()
	}
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
