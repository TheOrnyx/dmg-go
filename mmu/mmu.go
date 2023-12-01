package mmu

import (
	"github.com/TheOrnyx/gameboy-golor/cartridge"
)


// MMU Memory Mapped Unit - basically the hub for all the memory stuff
type MMU struct {
	boot [256]byte // the boot/bios ram | 0x0000 -> 0x00FF
	cart *cartridge.Cartridge
}
