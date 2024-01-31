package ppu

type Sprite struct {
	PosY byte // the y-position for the sprites vertical position with offset of 16
	PosX byte // the x-position for the sprites horizontal position with offset of 8
	Index byte // The tile index of the sprite
	Flags byte // the attributes and flags for the sprites
	
}
