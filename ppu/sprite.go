package ppu

type Sprite struct {
	PosY  byte // the y-position for the sprites vertical position with offset of 16
	PosX  byte // the x-position for the sprites horizontal position with offset of 8
	Index byte // The tile index of the sprite
	Flags byte // the attributes and flags byte for the sprites
	Position byte // the position of the sprite in OAM
}

// priority if true background and window colors 1-3 will draw over the obj
func (s *Sprite) priority() bool {
	return s.Flags>>7 == 0x01
}

// yFlip if true then the object is flipped vertically
func (s *Sprite) yFlip() bool {
	return (s.Flags>>6)&0x01 == 0x01
}

// xFlip if true then the object is flipped horizontally
func (s *Sprite) xFlip() bool {
	return (s.Flags>>5)&0x01 == 0x01
}

// dmgPalette true = use OBP1, false = use OBP0 (not used in CGB mode)
// NOTE - use the ppu method instead as that returns the actual palette value
func (s *Sprite) dmgPalette() bool {
	return (s.Flags>>4)&0x01 == 0x01
}

// bank true = fetch tile from vram bank 1, false = fetch tile from vram bank 0 (CGB mode only)
func (s *Sprite) bank() bool {
	return (s.Flags>>3)&0x01 == 0x01
}

// cgbPalette the first 3 bits of flags. chooses which palette from OBP0-7 to use (CGB mode only)
func (s *Sprite) cgbPalette() byte {
	return s.Flags & 0x07
}
