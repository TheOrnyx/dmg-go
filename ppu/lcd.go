package ppu

import "log"

// LCDReg the lcd IO registers
// TODO - implement CGB palette stuff
type LCDReg struct {
	Control byte // LCD Control byte (0xFF40)
	Stat    byte // LCD Status byte aka STAT (0xFF41)
	LY      byte // LCD Y coordinate (0xFF44)
	LYC     byte // LY compare - compares value of LYC and LY reg (0xFF45)

	// Palettes (non-cgb mode)
	BGP  byte // BG palette data in non-cgb mode (0xFF47)
	OBP0 byte // obj palette 0 data (0xFF48)
	OBP1 byte // obj palette 1 data (0xFF49)

	// Position and Scrolling
	// NOTE - check these are the right way around
	SCY byte // Viewport Y pos (0xFF42)
	SCX byte // Viewport X Pos (0xFF43)
	WY  byte // Window Y pos (0xFF4A)
	WX  byte // Window X pos (0xFF4B) - NOTE WX is position plus 7 so WX=7 is very left
}

// ReadByte read and return value at addr
func (l *LCDReg) ReadByte(addr uint16) byte {
	switch addr {
	case 0xFF40: // Control reg
		return l.Control
		
	case 0xFF41: // STAT
		return l.Stat
		
	case 0xFF44: // LY
		return l.LY
		
	case 0xFF45: // LYC
		return l.LYC
		
	case 0xFF47: // BGP
		return l.BGP
		
	case 0xFF48: // OBP0
		return l.OBP0
		
	case 0xFF49: // OBP1
		return l.OBP1
		
	case 0xFF42: // SCY
		return l.SCY
		
	case 0xFF43: // SCX
		return l.SCX

	case 0xFF4A: // WY
		return l.WY
		
	case 0xFF4B: // WX
		return l.WX

	default:
		log.Println("Unkown read in LCD IO at unkown addr:", addr)
	}
	
	return 0xFF
}

// WriteByte write data to given addr
func (l *LCDReg) WriteByte(addr uint16, data byte) {
	switch addr {
	case 0xFF40: // Control reg
		l.Control = data

	case 0xFF41: // STAT
		l.Stat = data

	case 0xFF44: // LY

	case 0xFF45: // LYC
		l.LYC = data

	case 0xFF47: // BGP
		l.BGP = data

	case 0xFF48: // OBP0
		l.OBP0 = data

	case 0xFF49: // OBP1
		l.OBP1 = data

	case 0xFF42: // SCY
		l.SCY = data

	case 0xFF43: // SCX
		l.SCX = data

	case 0xFF4A: // WY
		l.WY = data

	case 0xFF4B: // WX
		l.WX = data
		
	default:
		log.Println("Unkown write to LCD IO at unkown addr:", addr)
	}
}
