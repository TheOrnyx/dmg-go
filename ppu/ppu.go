// package PPU is the Picture Processing Unit for the gameboy emulator
// handles all the the PPU stuff
package ppu

import (
	"fmt"

	"github.com/TheOrnyx/gameboy-golor/timer"
)

///////////////
// PPU stuff //
///////////////

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

// ReadByte read byte from addr in OAM
func (o *OAM) ReadByte(addr uint16) byte {
	return o.Data[addr]
}

// WriteByte write data to addr in OAM
func (o *OAM) WriteByte(addr uint16, data byte) {
	o.Data[addr] = data
}

// getSprite get sprite information at given 0-based index and return it
func (p *PPU) getSprite(index uint16) Sprite {
	startAddr := 0xFE00 + (4 * index)
	flagByte := p.ReadByte(startAddr + 3)

	return Sprite{
		PosY:      p.ReadByte(startAddr),
		PosX:      p.ReadByte(startAddr + 1),
		Index:     p.ReadByte(startAddr + 2),
		FlagsByte: flagByte,
		Flags:     getFlagsFromByte(flagByte),
	}
}

// PPU Modes Consts
const (
	OAMScan = 2 // OAM scan mode
	Drawing = 3 // Drawing mode
	HBlank  = 0 // H-blank mode
	VBlank  = 1 // V-blank mode
)

type PPU struct {
	VRAM       VideoRam
	OAM        OAM // the Object Attribute Memory
	LCD        LCDReg
	timer      *timer.Timer
	Screen     Screen   // The screen to store the scanlines in
	Fetcher    Fetcher  // the Pixel FIFO fetcher
	bgTileMap  uint16   // the tilemap to use for background tiles
	winTileMap uint16   // the tilemap to use for window tiles
	tileData   uint16   // the tileData address mode to use
	drawX      uint8    // the current X position relative to the LCD to draw to
	cycles     uint16   // the current cycles for the current scanline
	sprites    []Sprite // Sprites from OAM for current scanline
}

// NewPPU create and return a new ppu
func NewPPU(timer *timer.Timer) *PPU {
	ppu := new(PPU)
	ppu.timer = timer
	ppu.Fetcher = Fetcher{ppu: ppu}

	return ppu
}

// ReadByte read byte at location addr in PPU
func (p *PPU) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF: // vram

		newAddr := addr - 0x8000
		return p.VRAM.ReadByte(newAddr)

	case addr >= 0xFE00 && addr <= 0xFE9F: // OAM
		newAddr := addr - 0xFE00
		return p.OAM.ReadByte(newAddr)
	}
	return 0xFF
}

// WriteByte write byte value data to address addr in ppu
func (p *PPU) WriteByte(addr uint16, data byte) {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF: // vram
		newAddr := addr - 0x8000
		p.VRAM.WriteByte(newAddr, data)

	case addr >= 0xFE00 && addr <= 0xFE9F: // OAM
		newAddr := addr - 0xFE00
		p.OAM.WriteByte(newAddr, data)
	}
}

// setPPUMode set the current ppu mode in the STAT register
func (p *PPU) setPPUMode(mode int) {
	switch mode {
	case HBlank:
		p.LCD.Stat = (p.LCD.Stat & 0xFC)
	case VBlank:
		p.LCD.Stat = (p.LCD.Stat & 0xFC) | 0x01
	case OAMScan:
		p.LCD.Stat = (p.LCD.Stat & 0xFC) | 0x02
	case Drawing:
		p.LCD.Stat = (p.LCD.Stat & 0xFC) | 0x03
	}
}

// isLCDEnabled check if the LCD is enabled and return true if it is
func (p *PPU) isLCDEnabled() bool {
	return (p.LCD.Control>>7)&0x01 == 0x01
}

// String get debug string representation of the PPU
func (p *PPU) String() string {
	var mode string
	switch {
	case p.LCD.LY >= 144:
		mode = "VBLANK MODE"
	case p.cycles <= 80:
		mode = "Scan OAM Mode"
	case p.drawX < 160:
		mode = "Draw mode"
	default:
		mode = "Hblank mode"
	}

	return fmt.Sprintf("Currently in %v, drawX:%v, LY:%v, Cycles:%v", mode, p.drawX, p.LCD.LY, p.cycles)
}

// reloadInfo check and set appropriate variable information
// Such as the tilemap and tiledata
func (p *PPU) reloadInfo() {
	if (p.LCD.Control>>3)&0x01 == 0x01 {
		p.bgTileMap = 0x9C00
	} else {
		p.bgTileMap = 0x9800
	}

	if (p.LCD.Control>>6)&0x01 == 0x01 {
		p.winTileMap = 0x9C00
	} else {
		p.winTileMap = 0x9800
	}

	if (p.LCD.Control>>4)&0x01 == 0x01 {
		p.tileData = 0x8000
	} else {
		p.tileData = 0x8800
	}
}

// Step step the PPU and draw appropriate things
func (p *PPU) Step() {
	if p.cycles >= 456 {
		p.nextScanline()
	}
	p.reloadInfo()

	if p.LCD.LY >= 144 {
		p.setPPUMode(VBlank)
		p.vBlank()
	} else if p.cycles <= 80 {
		p.setPPUMode(OAMScan)
		p.scanOAM(p.LCD.LY)
	} else if p.drawX < 160 {
		p.setPPUMode(Drawing)
		p.drawMode(p.LCD.LY)
	} else {
		p.setPPUMode(HBlank)
		p.hBlank()
	}

	if p.LCD.LY == p.LCD.LYC {
		p.LCD.Stat |= 0x04 // TODO - check
	}
}

// nextScanline reset specific fields and set stuff according to scanline
// Also sets stuff to the next Frame if needed
func (p *PPU) nextScanline() {
	p.drawX = 0
	p.setPPUMode(OAMScan)
	p.LCD.LY++
	if p.LCD.LY >= 154 {
		p.LCD.LY = 0
	}
	p.cycles = 0
	p.sprites = []Sprite{}
	p.Fetcher.posX = 0 // TODO - check
}

// tick tick timer by cycles amount and increase field
func (p *PPU) tick(cycles int) {
	p.cycles += uint16(cycles)
	p.timer.TickT(cycles)
}

// scanOAM scan the OAM for a max of 10 sprites for the scanline and add them to the spritelist if applicable
// TODO - maybe do the like sprite validity check in the getSprite function cuz it'd be more optimized idk
func (p *PPU) scanOAM(scanY uint8) {
	var spriteHeight = p.LCD.controlObjSize()

	p.tick(2)
	if len(p.sprites) >= 10 {
		return
	}

	sprite := p.getSprite(uint16(len(p.sprites)))
	if scanY+16 >= sprite.PosY && scanY+16 < sprite.PosY+spriteHeight {
		p.sprites = append(p.sprites, sprite)
	}
}

// hBlank basically just ticks the timer to get the scanlines cycles up to 456
// cycles is the amount of cycles this scanline has used so far
func (p *PPU) hBlank() {

	p.tick(1)
}

// vBlank run the vblank mode for a single scanline (meant to be used in a for loop for the 10 scanlines)
func (p *PPU) vBlank() {
	p.tick(1)
}

// drawMode the drawing mode, return number of cycles the draw mode took
func (p *PPU) drawMode(scanY uint8) {

	p.Fetcher.Step()
	p.tick(2)

	if len(p.Fetcher.FIFO.Pixels) != 0 {
		popPixel := p.Fetcher.FIFO.Pixels[0]
		p.Fetcher.FIFO.Pixels = p.Fetcher.FIFO.Pixels[1:]
		p.Screen.FinalScreen[scanY][p.drawX] = popPixel
		fmt.Println(popPixel)
		fmt.Println(p.Fetcher.FIFO.Pixels)
		p.drawX++
	}

	// TODO - implement this properly
}

// getTileIndex get specific background tile index based on x and y coordinate
func (p *PPU) getBackgTileIndex(x, y uint8) uint8 {
	startAddr := p.bgTileMap + (32 * uint16(y))
	return p.ReadByte(startAddr + uint16(x))
}
