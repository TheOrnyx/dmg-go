package ppu

type Sprite struct {
	PosY      byte        // the y-position for the sprites vertical position with offset of 16
	PosX      byte        // the x-position for the sprites horizontal position with offset of 8
	Index     byte        // The tile index of the sprite
	FlagsByte byte        // the attributes and flags byte for the sprites
	Flags     spriteFlags // the attributes and flags struct for a sprite object
}

// Struct for storing flags for the sprites
type spriteFlags struct {
	priority   bool // true if background and window colors 1-3 will draw over obj
	yFlip      bool // true if object is flipped vertically
	xFlip      bool // true if object is flipped horizontally
	dmgPalette bool // true = use OBP1, false = use OBP0 (not used in CGB mode)
	bank       bool // true = fetch tile from vram bank 1, false = fetch tile from vram bank 0 (CGB mode only)
	cgbPalette byte // the last 3 bits (the 1, 2 and 4 bits). chooses which palette from OBP0-7 to use (CGB mode only)
}

// getFlagsFromByte make and return a spriteflags struct from the given flags byte
func getFlagsFromByte(data byte) spriteFlags {
	return spriteFlags{
		priority:   data>>7 == 0x01,
		yFlip:      (data>>6)&0x01 == 0x01,
		xFlip:      (data>>5)&0x01 == 0x01,
		dmgPalette: (data>>4)&0x01 == 0x01,
		bank:       (data>>3)&0x01 == 0x01,
		cgbPalette: data & 0x07,
	}
}
