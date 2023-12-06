package mmu

import (
	"github.com/TheOrnyx/gameboy-golor/cartridge"
)


// MMU Memory Mapped Unit - basically the hub for all the memory stuff
type MMU struct {
	boot [256]byte // the boot/bios ram | 0x0000 -> 0x00FF
	cart *cartridge.Cartridge
}

// ReadByte read and return the byte located at address addr
func (mmu *MMU) ReadByte(addr uint16) byte {
	// TODO finish this
	return 0
}

// WriteByte write byte value data to location specified in address addr
func (mmu *MMU) WriteByte(addr uint16, data byte)  {
	
}
