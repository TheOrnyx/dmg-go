package mmu

import (
	"fmt"

	"github.com/TheOrnyx/gameboy-golor/cartridge"
	"github.com/TheOrnyx/gameboy-golor/timer"
)

///////////////
//   TODO's  //
///////////////
//
// TODO - Check whether or not reading from vram and work ram needs to be masked like & 0x0FFF or smth
// TODO - IMPLEMENT THE BOOT ROM IMPORTANT PLEASE
// TODO - finish teh read and write for io

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

type WorkRam struct {
	RAM [0x2000]byte
}

// ReadByte read byte from active bank in vram
func (w *WorkRam) ReadByte(addr uint16) byte {
	// bankAddr := addr & 0x0FFF // TODO - check whether to mask
	return w.RAM[addr]
}

// WriteByte write given data to addr
func (w *WorkRam) WriteByte(addr uint16, data uint8) {
	w.RAM[addr] = data // TODO - check whether to mask
}

// IO the struct for the IO registers in the MMU
type IO struct {
	JoypadInput    byte         // joypad input byte					(0xFF00)
	SerialTransfer [2]byte      // serial transfer						(0xFF01 - 0xFF02)
	TimerControl   *timer.Timer // timer and divider					(0xFF04 - 0xFF07)
	Audio          [23]byte     // Audio								(0xFF10 - 0xFF26)
	Wave           [16]byte     // Wave Pattern							(0xFF30 - 0xFF3F)
	LCD            [12]byte     // LCD control and other stuff			(0xFF40 - 0xFF4B)
	VramBankSel    byte         // CGB byte for swapping vram bank		(0xFF4F)
	BootROMEnabled byte         // Set to non-zero to disable boot rom	(0xFF50)
	VramDMA        [5]byte      // VRAM DMA for CGB						(0xFF51 - 0xFF55)
	Palettes       [4]byte      // Background and OBJ Palettes in CGB	(0xFF68 - 0xFF6B)
	WramBankSel    byte         // CGB work ram bank select				(0xFF70)
}

// ReadByte read and return byte in addr from the IO registers
func (io *IO) ReadByte(addr uint16) byte {
	switch {
	case addr == 0xFF00: // joypad
		return io.JoypadInput

	case addr >= 0xFF01 && addr <= 0xFF02: // serial transfer
		return io.SerialTransfer[addr-0xFF01]

	case addr >= 0xFF04 && addr <= 0xFF07: // Timer and divider
		return io.TimerControl.Read(addr)
		
	case addr >= 0xFF10 && addr <= 0xFF26: // AUDIO
		
	case addr >= 0xFF30 && addr <= 0xFF3F: // Wave pattern

	case addr >= 0xFF40 && addr <= 0xFF4B: // LCD

	case addr == 0xFF4F: // Vram Bank Select (CGB)

	case addr == 0xFF50: // Boot ROM
		return io.BootROMEnabled

	case addr >= 0xFF51 && addr <= 0xFF55: // VRAM DMA (CGB)

	case addr >= 0xFF68 && addr <= 0xFF6B: // Palettes (CGB)

	case addr == 0xFF70: // Wram bank Select (CGB)

	}

	return 0
}

// WriteByte write given byte data to addr in the io registers
// TODO - finish implementing
func (io *IO) WriteByte(addr uint16, data byte) {
	switch {
	case addr == 0xFF00: // joypad
		io.JoypadInput = data // TODO - actually implement the proper input

	case addr >= 0xFF01 && addr <= 0xFF02: // serial transfer
		// TODO - maybe put silly debug stuff here
		if addr == 0xFF01 {
			fmt.Println("writing to serial:", data)
		}
		io.SerialTransfer[addr - 0xFF01] = data

	case addr >= 0xFF04 && addr <= 0xFF07: // Timer and divider
		io.TimerControl.Write(addr, data)
		
	case addr >= 0xFF10 && addr <= 0xFF26: // AUDIO

	case addr >= 0xFF30 && addr <= 0xFF3F: // Wave pattern

	case addr >= 0xFF40 && addr <= 0xFF4B: // LCD

	case addr == 0xFF4F: // Vram Bank Select (CGB)

	case addr == 0xFF50: // Boot ROM
		io.BootROMEnabled = data

	case addr >= 0xFF51 && addr <= 0xFF55: // VRAM DMA (CGB)

	case addr >= 0xFF68 && addr <= 0xFF6B: // Palettes (CGB)

	case addr == 0xFF70: // Wram bank Select (CGB)

	}
}

// MMU Memory Mapped Unit - basically the hub for all the memory stuff
type MMU struct {
	// bootROM [256]byte // the boot/bios ram | 0x0000 -> 0x00FF
	VRAM             VideoRam
	WRAM             WorkRam
	HRAM             [0x7F]byte // High ram
	IO               IO
	OAM              [0xA0]byte // object attribute memory
	interruptEnabled byte
	interruptsFlag   byte
	Cart             *cartridge.Cartridge
}

// NewMMU create and return a new MMU
func NewMMU(cart *cartridge.Cartridge, timer *timer.Timer) *MMU {
	newMMU := new(MMU)
	newMMU.Cart = cart
	newMMU.IO.TimerControl = timer
	return newMMU
}

// ReadByte read and return the byte located at address addr
// TODO - finish and check
func (mmu *MMU) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x7FFF: // Fixed cart bank (don't need to implement switchable as different since mbc handles that)
		return mmu.Cart.MBC.ReadByte(addr)

	case addr >= 0x8000 && addr <= 0x9FFF: // Video RAM
		newAddr := addr - 0x8000
		return mmu.VRAM.ReadByte(newAddr)

	case addr >= 0xA000 && addr <= 0xBFFF: // external ram on cart
		newAddr := addr - 0xA000
		return mmu.Cart.MBC.ReadByte(newAddr)

	case addr >= 0xC000 && addr <= 0xCFFF: // first work ram bank
		newAddr := addr - 0xC000
		return mmu.WRAM.ReadByte(newAddr)

	case addr >= 0xD000 && addr <= 0xDFFF: // switchable Work ram bank (but for regular gameboy it's just the same)
		newAddr := addr - 0xD000
		return mmu.WRAM.ReadByte(newAddr)

		// Ignoring the 0xE000 -> 0xFDFF - nintendo says not allowed >:(

	case addr == 0xFF0F: // Interrupt Flag
		return mmu.interruptsFlag

	case addr == 0xFFFF: // interrupts enabled
		return mmu.interruptEnabled

	case addr >= 0xFE00 && addr <= 0xFE9F: // object attribute memory
		newAddr := addr - 0xFE00
		return mmu.OAM[newAddr]

		// ignore 0xFEA0 -> 0xFEFF - nintendo says not allowed again

	case addr >= 0xFF00 && addr <= 0xFF7F: // I/O registers
		return mmu.IO.ReadByte(addr)

	case addr >= 0xFF80 && addr <= 0xFFFE: // high ram (HRAM)
		newAddr := addr - 0xFF80
		return mmu.HRAM[newAddr]

	default:

	}
	return 0
}

// WriteByte write byte value data to location specified in address addr
func (mmu *MMU) WriteByte(addr uint16, data byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x7FFF: // Read from Cart
		mmu.Cart.MBC.WriteByte(addr, data)
		
	case addr >= 0x8000 && addr <= 0x9FFF: // Video ram
		newAddr := addr - 0x8000
		mmu.VRAM.WriteByte(newAddr, data)

	case addr >= 0xA000 && addr <= 0xBFFF: // External ram on cart
		newAddr := addr - 0xA000
		mmu.Cart.MBC.WriteByte(newAddr, data)

	case addr >= 0xC000 && addr <= 0xDFFF: // work ram
		newAddr := addr - 0xC000
		mmu.WRAM.WriteByte(newAddr, data)

	case addr == 0xFF0F: // interruptsFlag
		mmu.interruptsFlag = data

	case addr == 0xFFFF: // interrupts enabled
		mmu.interruptEnabled = data

	case addr >= 0xFF00 && addr <= 0xFF7F: // I/O registers
		mmu.IO.WriteByte(addr, data)

	case addr >= 0xFF80 && addr <= 0xFFFE: // high ram
		mmu.HRAM[addr - 0xFF80] = data
		

	}
}
