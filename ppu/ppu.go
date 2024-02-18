// package PPU is the Picture Processing Unit for the gameboy emulator
// handles all the the PPU stuff
package ppu

import (
	"fmt"
	"slices"

	"github.com/TheOrnyx/gameboy-golor/timer"
)

///////////////
// PPU stuff //
///////////////

const (
	// Interrupt request codes
	vBlankInt = iota
	lcdInt
)

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

// PPU Modes Consts
const (
	OAMScanMode = 2 // OAM scan mode
	DrawMode    = 3 // Drawing mode
	HBlankMode  = 0 // H-blank mode
	VBlankMode  = 1 // V-blank mode
)

type PPU struct {
	VRAM             VideoRam
	OAM              OAM // the Object Attribute Memory
	LCD              LCDReg
	WLY              byte            // the LY for the window
	RequestInterrupt func(code byte) // function pointer to request interrupts
	timer            *timer.Timer
	Screen           Screen // The screen to store the scanlines in
	cycles           uint16 // the current cycles for the current scanline
}

// NewPPU create and return a new ppu
func NewPPU(timer *timer.Timer, requestInterrupt func(code byte)) *PPU {
	ppu := new(PPU)
	ppu.timer = timer
	ppu.RequestInterrupt = requestInterrupt
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
	case HBlankMode:
		p.LCD.Stat = (p.LCD.Stat & 0xFC)
		if p.LCD.modeZeroInt() {
			p.RequestInterrupt(lcdInt)
		}

	case VBlankMode:
		p.LCD.Stat = (p.LCD.Stat & 0xFC) | 0x01
		p.RequestInterrupt(vBlankInt)
		if p.LCD.modeOneInt() {
			p.RequestInterrupt(lcdInt)
		}

	case OAMScanMode:
		p.LCD.Stat = (p.LCD.Stat & 0xFC) | 0x02
		if p.LCD.modeTwoInt() {
			p.RequestInterrupt(lcdInt)
		}

	case DrawMode:
		p.LCD.Stat = (p.LCD.Stat & 0xFC) | 0x03
	}
}

// Mode get the current PPU Mode
func (p *PPU) Mode() byte {
	return p.LCD.StatMode()
}

// checkLYCInterrupt Check and set LYC value and set LYCInterrupt as well
func (p *PPU) checkLYCInterrupt() {
	var equal bool

	if p.LCD.LY == p.LCD.LYC {
		p.LCD.Stat |= 0x04 // TODO - check
		equal = true
	} else {
		p.LCD.Stat &= 0xFB
	}

	if equal && (p.LCD.Stat>>6)&0x01 == 0x01 {
		p.RequestInterrupt(vBlankInt)
	}
}

// String get debug string representation of the PPU
func (p *PPU) String() string {
	var mode string
	switch p.Mode() {
	case OAMScanMode:
		mode = "OAM Scan Mode"
	case HBlankMode:
		mode = "HBlank Mode"
	case DrawMode:
		mode = "Draw Mode"
	case VBlankMode:
		mode = "VBlank Mode"
	}

	return fmt.Sprintf("Currently in %v, LY:%v, Cycles:%v", mode, p.LCD.LY, p.cycles)
}

// Step step the PPU and draw appropriate things
func (p *PPU) Step(cycles uint16) {
	p.cycles += cycles
	p.checkLYCInterrupt()

	if p.cycles >= 456 { // go to next scanline
		p.cycles -= 456

		p.LCD.LY = (p.LCD.LY + 1) % 154
		p.checkLYCInterrupt()

		// vblank
		if p.LCD.LY >= 144 && p.Mode() != VBlankMode {
			p.setPPUMode(VBlankMode)

			p.WLY = 0

			if p.LCD.modeOneInt() {
				p.RequestInterrupt(lcdInt) // request lcd interrupt
			}
			p.RequestInterrupt(vBlankInt)
		}
	}

	if p.LCD.LY < 144 {
		switch {
		case (p.cycles <= 80) && p.Mode() != OAMScanMode:
			// p.Screen.Reset() // TODO - CHECK THIS
			p.setPPUMode(OAMScanMode)

		case (p.cycles >= 81 && p.cycles <= 252) && p.Mode() != DrawMode: // NOTE - idk why 252, that's what a guy did
			p.setPPUMode(DrawMode)

		case (p.cycles >= 253 && p.cycles <= 455) && p.Mode() != HBlankMode:
			p.setPPUMode(HBlankMode)
			p.renderScanline()
		}
	}
}

// renderScanline create the individual scanlines for the specific layers and combine them onto the main screen
func (p *PPU) renderScanline() {
	if !p.LCD.isLcdOn() {
		return
	}

	p.DrawBGScanline()
	p.DrawWinScanline()
	p.DrawObjectScanline()
	p.Screen.CombineLine(int(p.LCD.LY))
}

// tick tick timer by cycles amount and increase field
func (p *PPU) tick(cycles int) {
	p.cycles += uint16(cycles)
	p.timer.TickT(cycles)
}

// getSprite get sprite information at given 0-based index and return it
func (p *PPU) getSprite(index uint16) Sprite {
	startAddr := 0xFE00 + (4 * index)

	return Sprite{
		PosY:  p.ReadByte(startAddr),
		PosX:  p.ReadByte(startAddr + 1),
		Index: p.ReadByte(startAddr + 2),
		Flags: p.ReadByte(startAddr + 3),
		Position: byte(index),
	}
}

// scanOAM scan the OAM and return a max of 10 sprites for the scanline
func (p *PPU) scanOAM(scanY uint8) []Sprite {
	var spriteHeight = p.LCD.objSize()
	var sprites []Sprite

	for i := range uint16(40) { // scan all objects in OAM
		if len(sprites) >= 10 {
			break
		}
		
		sprite := p.getSprite(i)
		if scanY+16 >= sprite.PosY && scanY+16 < sprite.PosY+spriteHeight {
			sprites = append(sprites, sprite)
		}
	}

	return sprites
}

// DrawBGScanline draw a background scanline and push it to the screen layer
func (p *PPU) DrawBGScanline() {
	if !p.LCD.EnableBGWin() {
		return
	}

	tileMapAddr := p.LCD.BGTileMap()

	for col := range uint8(160) {
		x := col + p.LCD.SCX
		y := p.LCD.LY + p.LCD.SCY
		pixelVal := p.getTilePixel(x, y, tileMapAddr)
		color := (p.LCD.BGP >> (pixelVal * 2)) & 0x03
		p.Screen.Background[p.LCD.LY][col] = Pixel{Color: color, Opaque: true}
	}
}

// DrawWinScanline draw a window scanline and push it to the screen layer
func (p *PPU) DrawWinScanline() {
	if !p.LCD.WindowEnabled() || p.LCD.LY < p.LCD.WY {
		return
	}

	tilemapAddr := p.LCD.WinTileMap()
	var pixelDrawn bool

	for x := range uint8(160) {
		winX := 0 - (int(p.LCD.WX) - 7) + int(x)

		// don't draw if pixel off screen
		if winX < 0 || winX >= 160 {
			continue
		}

		pixel := p.getTilePixel(uint8(winX), p.WLY, tilemapAddr)

		p.Screen.Window[p.LCD.LY][x] = Pixel{Color: pixel, Opaque: true}
		pixelDrawn = true
	}

	if pixelDrawn {
		p.WLY += 1
	}
}

// DrawObjectScanline draw a scanline for the objects and push to object layer in the screen
// TODO - maybe split into seperate functions and simplify
func (p *PPU) DrawObjectScanline() {
	if !p.LCD.objEnabled() {
		return
	}

	sprites := p.scanOAM(p.LCD.LY)
	slices.Reverse(sprites)
	
	for i := range sprites {
		sprite := sprites[i]

		var spriteNum uint16
		if p.LCD.objSize() == 16 {
			spriteNum = uint16(sprite.Index) & 0xFE
		} else {
			spriteNum = uint16(sprite.Index) & 0xFF
		}
		
		var spriteY uint16
		height := uint16(p.LCD.objSize())
		
		if sprite.yFlip() {
			spriteY = (height-1 - (uint16(p.LCD.LY) - (uint16(sprite.PosY) - 16))) // TODO - very gross, check
		} else {
			spriteY = (uint16(p.LCD.LY) - (uint16(sprite.PosY) - 16))
		}

		spriteDataAddr := 0x8000 + (spriteNum * 16) + (spriteY * 2)
		spriteDataLow := p.ReadByte(spriteDataAddr)
		spriteDataHigh := p.ReadByte(spriteDataAddr+1)

		for pixel := range uint8(7) {
			posX := sprite.PosX - 8 // the relative posX
			outOfRangeLeft := (int(sprite.PosX) - 8) + int(pixel) < 0
			outOfRangeRight := (sprite.PosX - 8) + pixel >= 160
			if outOfRangeLeft || outOfRangeRight {
				continue
			}

			coverPixel := !p.Screen.Background[p.LCD.LY][posX].Opaque || !p.Screen.Window[p.LCD.LY][posX].Opaque
			if !p.LCD.EnableBGWin() && sprite.priority() && coverPixel {
				continue
			}

			var pixelNum byte = pixel
			if sprite.xFlip() {
				pixelNum = 7-pixel
			}

			pixelVal := getPixel(spriteDataLow, spriteDataHigh, pixelNum)
			color := (p.getDMGPalette(&sprite) >> (pixelVal * 2)) & 0x03
			p.Screen.Objects[p.LCD.LY][posX+pixel] = Pixel{Color: color, Opaque: color != 0}
		}
	}
}

// getDMGPalette get the actual palette data for a given obj
func (p *PPU) getDMGPalette(sprite *Sprite) byte {
	if sprite.dmgPalette() {
		return p.LCD.OBP1
	}

	return p.LCD.OBP0
}

// getTilePixel get the pixel value for the pixel at x and y based on tileMapAddr
// Heavily based on https://github.com/ablakey/gameboy/blob/dc2abcae57f271cb7873f9fd5cc77cd6e30fb4d1/src/guest/systems/ppu.rs#L5
// As I'm too stupid to understand myself
func (p *PPU) getTilePixel(x, y uint8, tileMapAddr uint16) uint8 {
	tileRowNum := y / 8
	tileColNum := x / 8
	tileNum := uint16(tileRowNum)*32 + uint16(tileColNum)

	tileDataNum := p.ReadByte(tileMapAddr + tileNum)
	tileDataAddr := getTileDataAddress(p.LCD.TileDataAddrMode(), tileDataNum)

	pixelRowNum := y % 8
	pixelColNum := x % 8

	tileRowIndex := tileDataAddr + (uint16(pixelRowNum) * 2)
	tileDataLow := p.ReadByte(tileRowIndex)
	tileDataHigh := p.ReadByte(tileRowIndex + 1)

	return getPixel(tileDataLow, tileDataHigh, pixelColNum)
}

// getPixel get the two bits for the color pixel based on the two bytes and the pixel num
func getPixel(tileLow, tileHigh, pixelNum uint8) uint8 {
	pixel0 := (tileLow >> (7 - pixelNum)) & 0x01
	pixel1 := (tileHigh >> (7 - pixelNum)) & 0x01
	return (pixel1 << 1) + pixel0
}

// getTileDataAddress get the address for the tile data based on the addressing mode
func getTileDataAddress(baseAddress uint16, tileNum uint8) uint16 {
	if baseAddress == 0x8800 {
		return baseAddress + (uint16(int16(int8(tileNum))+128) * 16)
	}

	return baseAddress + (uint16(tileNum) * 16)
}
