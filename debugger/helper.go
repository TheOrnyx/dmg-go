package debugger

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// drawBox draw a box at given location with given width and height
func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// draw corners if needed
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	
}

// drawText draws given text into specified position on screen
func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			// row++
			// col = x1
		}
		if row > y2 {
			break
		}
	}
}

// updateDimensions update the maxX and maxY vars to match the current terminal size
func updateDimensions(s tcell.Screen) {
	maxX, maxY = s.Size()
	maxX -= 1
	maxY -= 1 // because the max is one off the screen
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

// initTcell initialize the tcell screen and return it
func initTcell() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("Failed to create TCell screen: %v", err)
	}

	err = s.Init()
	if err != nil {
		return nil, fmt.Errorf("Failed to init TCell screen: %v", err)
	}

	return s, nil
}

// switchPanel switch panel to given panel and wrap
func (d *Debugger) switchPanel(new int)  {
	if new >= PanelCount {
		d.ActivePanel = CPU
	} else if new < 0 {
		d.ActivePanel = PanelCount-1
	} else {
		d.ActivePanel = new
	}
}
