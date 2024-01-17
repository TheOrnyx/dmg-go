package mmu

import (
	"fmt"

	"github.com/TheOrnyx/gameboy-golor/cartridge"
)

///////////////
//   TODO's  //
///////////////
// 
// TODO - Check whether or not reading from vram and work ram needs to be masked like & 0x0FFF or smth

type VideoRam struct {
	RAM  [2][0x2000]byte
	bank int
}

// SwitchBank switch the videoRam ActiveBank to newIndex bank
func (v *VideoRam) SwitchBank(newIndex int) error {
	if newIndex < 0 || newIndex >= len(v.RAM) {
		return fmt.Errorf("Out of range")
	}

	v.bank = newIndex
	return nil
}

// ReadByte read byte from active bank in vram
func (v *VideoRam) ReadByte(addr uint16) byte {
	bankAddr := addr &0x1FFF
	return v.RAM[v.bank][bankAddr]
}

// WriteByte write given data to addr
func (v *VideoRam) WriteByte(addr uint16, data uint8)  {
	v.RAM[v.bank][addr] = data // TODO - check
}

type WorkRam struct {
	RAM  [8][0x1000]byte
	bank int // index of the active bank
}

// SwitchBank switch the workram ActiveBank to newIndex bank
func (w *WorkRam) SwitchBank(newIndex int) error {
	if newIndex < 0 || newIndex >= len(w.RAM) {
		return fmt.Errorf("Out of range")
	}

	w.bank = newIndex
	return nil
}

// ReadByte read byte from active bank in vram
func (w *WorkRam) ReadByte(addr uint16) byte {
	bankAddr := addr & 0x0FFF
	if addr >= 0xC000 && addr <= 0xCFFF {
		return w.RAM[0][bankAddr]
	} else {
		return w.RAM[w.bank][bankAddr]
	}
}

// WriteByte write given data to addr
func (w *WorkRam) WriteByte(addr uint16, data uint8)  {
	if addr >= 0xC000 && addr <= 0xCFFF {
		w.RAM[0][addr] = data
	}
	
	w.RAM[w.bank][addr] = data
}

// MMU Memory Mapped Unit - basically the hub for all the memory stuff
type MMU struct {
	// boot [256]byte // the boot/bios ram | 0x0000 -> 0x00FF
	VRAM VideoRam
	WRAM WorkRam
	HRAM [0x7F]byte // High ram
	IO   [0x80]byte // I/O registers
	OAM  [0xA0]byte // object attribute memory
	Cart *cartridge.Cartridge
}


// NewMMU create and return a new MMU
func NewMMU(cart *cartridge.Cartridge) *MMU {
	newMMU := new(MMU)
	newMMU.Cart = cart
	return newMMU
}

// ReadByte read and return the byte located at address addr
func (mmu *MMU) ReadByte(addr uint16) byte {
	// TODO finish this
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // Fixed cart bank
		return mmu.Cart.MBC.ReadByte(addr)

	case addr >= 0x4000 && addr <= 0x7FFF: // cart rom bank 01 (switchable)
		return mmu.Cart.MBC.ReadByte(addr)

	case addr >= 0x8000 && addr <= 0x9FFF: // Video RAM
		return mmu.VRAM.ReadByte(addr)

	case addr >= 0xA000 && addr <= 0xBFFF: // external ram on cart
		return mmu.Cart.MBC.ReadByte(addr)

	case addr >= 0xC000 && addr <= 0xCFFF: // first work ram bank (TODO - check this is correct)
		return mmu.WRAM.ReadByte(addr)

	case addr >= 0xD000 && addr <= 0xDFFF: // switchable Work ram bank
		return mmu.WRAM.ReadByte(addr)

		// Ignoring the 0xE000 -> 0xFDFF - nintendo says not allowed >:(

	case addr >= 0xFE00 && addr <= 0xFE9F: // object attribute memory
		return mmu.OAM[addr&0x009F]

		// ignore 0xFEA0 -> 0xFEFF - nintendo says not allowed again

	case addr >= 0xFF00 && addr <= 0xFF7F: // I/O registers
		return mmu.IO[addr-0xFF00] // TODO - check
		// TODO - Implement later

	case addr >= 0xFF80 && addr <= 0xFFFE: // high ram (HRAM)
		return mmu.HRAM[addr&(0xFFFF-0xFF80)]

	default:
		

	}
	return 0
}

// WriteByte write byte value data to location specified in address addr
func (mmu *MMU) WriteByte(addr uint16, data byte) {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF: // Video ram
		mmu.VRAM.WriteByte(addr, data)
	case addr >= 0xA000 && addr <= 0xBFFF: // External ram on cart
		mmu.Cart.MBC.WriteByte(addr, data)
	case addr >= 0xC000 && addr <= 0xDFFF: // work ram
		mmu.WRAM.WriteByte(addr, data)
	case addr >= 0xFF00 && addr <= 0xFF7F: // I/O stuff
		mmu.IO[addr - 0xFF00] = data
	default:
		
	// TODO - finish later
	}
}
