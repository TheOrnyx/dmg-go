// package PPU is the Picture Processing Unit for the gameboy emulator
// handles all the the PPU stuff
package ppu

type VideoRam struct {
	RAM [0x2000]byte
}

// ReadByte read byte from active bank in vram
func (v *VideoRam) ReadByte(addr uint16) byte {
	// bankAddr := addr & 0x1FFF
	return v.RAM[addr]
}

// WriteByte write given data to addr
func (v *VideoRam) WriteByte(addr uint16, data uint8) {
	v.RAM[addr] = data
}

type OAM struct {
	Data [0xA0]byte // The OAM data (probably replcae with like sprite objects later)
}

// PPU Modes Consts
const (
	OAMScan = 2 // OAM scan mode
	Drawing = 3 // Drawing mode
	HBlank = 0 // H-blank mode
	VBlank = 1 // V-blank mode
)

type PPU struct {
	VRAM VideoRam
	OAM OAM // the Object Attribute Memory
	LCD LCDReg

	Screen Screen // The screen to store the scanlines in
}

// ReadByte read byte at location addr in PPU
func (p *PPU) ReadByte(addr uint16) byte {

	return 0xFF
}

// WriteByte write byte value data to address addr in ppu
func (p *PPU) WriteByte(addr uint16, data byte) {
	
}
