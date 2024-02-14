package debugger

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TheOrnyx/gameboy-golor/emulator"
	"github.com/gdamore/tcell/v2"
)

// The panel constants
const (
	CPU = iota
	MMU
	PPU
	PanelCount = 3
)

var (
	maxX int // the max X to draw the debugger to
	maxY int // the max Y to draw the debugger to
)

const (
	leftX2 = 6 // the max X for the left box to draw to
	topY2  = 2 // the max Y for the top info box to draw to
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

type Debugger struct {
	Emu         *emulator.Emulator
	ActivePanel int // the active panel being used
	Screen      tcell.Screen
}

// DebugEmulatorDoctor run and debug emulator outputting in doctor format
func DebugEmulatorDoctor(emu *emulator.Emulator) {
	emu.CPU.ResetDebug()
	count := 1
	maxTests := 5500000 // set to 0 or below for infinite running >:3
	debugger := Debugger{Emu: emu}
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
	fmt.Println(emu.PPU.VRAM.RAM)
	fmt.Println(emu.PPU.OAM.Data)
}

// DebugEmulator run and debug an emulator
func DebugEmulator(emu *emulator.Emulator) {
	emu.CPU.ResetDebug()
	// fmt.Println(emu.MMU.Cart.ROM)
	fmt.Println("Beginning debug...")
	fmt.Println(emu.DebugInfo())
	running := true
	// debugger := Debugger{emu}

	for running {
		fmt.Println("---------------------------------------------")
		fmt.Printf("%v\nPPU Info: %v\n", emu.CPU, emu.PPU)
		emu.CPU.Step()
		// emu.PPU.Step()
		// emu.RenderScreen()

		if true {
			time.Sleep(time.Nanosecond * 230)
			// time.Sleep(time.Second/6)
		}
	}
	// fmt.Println("\n",serialOutput)
}

// DebugEmu debug and run an emulator with the TUI
func DebugEmu(emu *emulator.Emulator) {
	emu.CPU.ResetDebug()
	d := Debugger{Emu: emu, ActivePanel: 1}
	s, err := initTcell()
	if err != nil {
		log.Fatal(err)
	}
	d.Screen = s
	d.Screen.SetStyle(defStyle)
	d.Screen.Clear()
	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	running := true

	for running {
		updateDimensions(d.Screen)
		d.Screen.Show()
		d.Screen.Clear()

		ev := d.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				running = false

			case tcell.KeyTAB:
				d.switchPanel(d.ActivePanel + 1)
			case tcell.KeyBacktab:
				d.switchPanel(d.ActivePanel - 1)

			default: // search for runes instead
				switch ev.Rune() {
				case ' ':
					d.Emu.Step()
				}
			}
		}

		d.DrawUI()
	}
}

// DrawUI draw the UI for the debugger
func (d *Debugger) DrawUI() {
	// Draw the boxes first
	drawBox(d.Screen, 0, 0, maxX, topY2, defStyle)        // the top info box
	drawBox(d.Screen, 0, topY2+1, leftX2, maxY, defStyle) // the left box
	// drawBox(d.Screen, leftX2+2, topY2+1, maxX, maxY, defStyle) // the main box
	
	// Draw the menu labels
	midY := maxY / 2
	startY := midY - 1 // the Y value to start with
	drawText(d.Screen, 2, startY, 6, startY, defStyle, "CPU")
	drawText(d.Screen, 2, startY+1, 6, startY+1, defStyle, "MMU")
	drawText(d.Screen, 2, startY+2, 6, startY+2, defStyle, "PPU")
	drawText(d.Screen, 1, startY+d.ActivePanel, 2, startY+d.ActivePanel, defStyle.Bold(true), ">")

	//Draw top Info box information
	infoText := fmt.Sprintf("Cart Name: %v     Cart Type: %v", d.Emu.MMU.Cart.Title, d.Emu.MMU.Cart.MBCType)
	drawText(d.Screen, 1, 1, maxX-1, 1, defStyle, infoText)

	switch d.ActivePanel {
	case CPU:
		sepY := maxY - 10
		drawBox(d.Screen, leftX2+2, topY2+1, maxX, sepY, defStyle) // the instructions box
		drawBox(d.Screen, leftX2+2, sepY+1, maxX, maxY, defStyle)  // the instructions box
		d.drawCPUInstrPanel(sepY)
		d.drawCPUDataPanel(sepY)
	case MMU:
		drawBox(d.Screen, leftX2+2, topY2+1, maxX, maxY, defStyle) // the main box
		d.drawMMUPanel()
	}
}

// drawCPUInstrPanel draw the cpu instructions panel
func (d *Debugger) drawCPUInstrPanel(sepY int) {
	startY := 4
	endY := sepY - 2
	startX := 10

	clampedPrev := d.Emu.CPU.PrevInstructions
	if len(clampedPrev) > endY-startY {
		clampedPrev = clampedPrev[len(clampedPrev)-(endY-startY):]
	}

	for i := range clampedPrev {
		drawText(d.Screen, startX, startY+i, maxX-2, endY, defStyle, clampedPrev[i].String())
	}

	drawText(d.Screen, startX, endY+1, maxX-2, endY+1, defStyle.Background(tcell.ColorBlack), d.Emu.CPU.GetInstrDebug())
}

// drawCPUDataPanel draw the cpu panel for the cpu Data
func (d *Debugger) drawCPUDataPanel(sepY int) {
	startY := sepY + 2
	endY := maxY - 1
	startX := 10
	cpu := d.Emu.CPU

	pc, pcOne, pcTwo, pcThree := cpu.GetPCMem()
	firstLine := fmt.Sprintf("PC:0x%04X SP:0x%04X RegF:%v", cpu.PC, cpu.SP, cpu.Reg.F.String())
	secondLine := fmt.Sprintf("Regs:%v", cpu.Reg.String())
	thirdLine := fmt.Sprintf("PCMem: %v, %v, %v, %v", pc, pcOne, pcTwo, pcThree)
	drawText(d.Screen, startX, startY, maxX-1, endY, defStyle, firstLine)
	drawText(d.Screen, startX, startY+1, maxX-1, endY, defStyle, secondLine)
	drawText(d.Screen, startX, startY+2, maxX-1, endY, defStyle, thirdLine)
}

// drawMMUPanel draw the MMU panel
func (d *Debugger) drawMMUPanel() {
	startY := 4
	endY := maxY - 1
	startX := 10

	clampedPrev := d.Emu.MMU.DebugRecords
	if len(clampedPrev) > endY-startY {
		clampedPrev = clampedPrev[len(clampedPrev)-(endY-startY):]
	}

	for i := range clampedPrev {
		drawText(d.Screen, startX, startY+i, maxX-2, endY, defStyle, clampedPrev[i])
	}
}
