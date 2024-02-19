package mmu

import (
	"fmt"

	"github.com/TheOrnyx/gameboy-golor/cartridge"
	"github.com/TheOrnyx/gameboy-golor/joypad"
	"github.com/TheOrnyx/gameboy-golor/ppu"
	"github.com/TheOrnyx/gameboy-golor/timer"
)

///////////////
//   TODO's  //
///////////////
//
// TODO - Check whether or not reading from vram and work ram needs to be masked like & 0x0FFF or smth
// TODO - IMPLEMENT THE BOOT ROM IMPORTANT PLEASE
// TODO - finish teh read and write for io

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
	Joypad *joypad.Joypad       // joypad input         				(0xFF00)
	SerialTransfer [2]byte      // serial transfer						(0xFF01 - 0xFF02)
	TimerControl   *timer.Timer // timer and divider					(0xFF04 - 0xFF07)
	Audio          [23]byte     // Audio								(0xFF10 - 0xFF26)
	Wave           [16]byte     // Wave Pattern							(0xFF30 - 0xFF3F)
	LCD            *ppu.LCDReg  // LCD control and other stuff			(0xFF40 - 0xFF4B)
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
		return io.Joypad.Read()

	case addr >= 0xFF01 && addr <= 0xFF02: // serial transfer
		return io.SerialTransfer[addr-0xFF01]

	case addr >= 0xFF04 && addr <= 0xFF07: // Timer and divider
		return io.TimerControl.Read(addr)

	case addr >= 0xFF10 && addr <= 0xFF26: // AUDIO

	case addr >= 0xFF30 && addr <= 0xFF3F: // Wave pattern

	case addr >= 0xFF40 && addr <= 0xFF4B: // LCD
		// if addr == 0xFF44 {
		// 	return 0x90
		// }
		return io.LCD.ReadByte(addr)

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
		io.Joypad.WriteData(data)

	case addr >= 0xFF01 && addr <= 0xFF02: // serial transfer
		io.SerialTransfer[addr-0xFF01] = data

	case addr >= 0xFF04 && addr <= 0xFF07: // Timer and divider
		io.TimerControl.Write(addr, data)

	case addr >= 0xFF10 && addr <= 0xFF26: // AUDIO

	case addr >= 0xFF30 && addr <= 0xFF3F: // Wave pattern

	case addr >= 0xFF40 && addr <= 0xFF4B: // LCD
		io.LCD.WriteByte(addr, data)

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
	PPU              *ppu.PPU // the PPU (needed to acces VRAM and OAM)
	WRAM             WorkRam
	HRAM             [0x7F]byte // High ram
	IO               IO
	interruptEnabled byte
	interruptsFlag   byte
	Cart             *cartridge.Cartridge
	DebugMode        bool     // whether or not to record debug information
	DebugRecords     []string // the debug information for read and write operations
}

// NewMMU create and return a new MMU
func NewMMU(cart *cartridge.Cartridge, timer *timer.Timer, ppu *ppu.PPU, joypad *joypad.Joypad) *MMU {
	newMMU := new(MMU)
	newMMU.Cart = cart
	newMMU.PPU = ppu
	newMMU.IO.TimerControl = timer
	newMMU.IO.LCD = &ppu.LCD
	newMMU.IO.BootROMEnabled = 1
	newMMU.IO.Joypad = joypad
	newMMU.DebugMode = false // TODO - change later
	return newMMU
}

// ReadByte read and return the byte located at address addr
// TODO - finish and check
func (mmu *MMU) ReadByte(addr uint16) byte {	
	switch {
	case addr >= 0x0000 && addr <= 0x7FFF: // Fixed cart bank (don't need to implement switchable as different since mbc handles that)
		data := mmu.Cart.MBC.ReadByte(addr)
		mmu.addReadToDebug(addr, data, "Cart Banks")
		return data

	case addr >= 0x8000 && addr <= 0x9FFF: // Video RAM
		
		data := mmu.PPU.ReadByte(addr)
		mmu.addReadToDebug(addr, data, "VRAM")
		return data

	case addr >= 0xA000 && addr <= 0xBFFF: // external ram on cart
		newAddr := addr - 0xA000
		data := mmu.Cart.MBC.ReadByte(newAddr)
		mmu.addReadToDebug(addr, data, "External Cart RAM")
		return data

	case addr >= 0xC000 && addr <= 0xDFFF: // first work ram bank
		newAddr := addr - 0xC000
		data := mmu.WRAM.ReadByte(newAddr)
		mmu.addReadToDebug(addr, data, "WRAM bank 0")
		return data

		// Ignoring the 0xE000 -> 0xFDFF - nintendo says not allowed >:(
		
	case addr == 0xFF0F: // Interrupt Flag
		mmu.addReadToDebug(addr, mmu.interruptsFlag, "IF")
		return mmu.interruptsFlag

	case addr == 0xFFFF: // interrupts enabled
		mmu.addReadToDebug(addr, mmu.interruptEnabled, "IEF")
		return mmu.interruptEnabled

	case addr >= 0xFE00 && addr <= 0xFE9F: // object attribute memory
		data := mmu.PPU.ReadByte(addr)
		mmu.addReadToDebug(addr, data, "OAM")
		return data

		// ignore 0xFEA0 -> 0xFEFF - nintendo says not allowed again

	case addr >= 0xFF00 && addr <= 0xFF7F: // I/O registers
		data := mmu.IO.ReadByte(addr)
		mmu.addReadToDebug(addr, data, "I/O")
		return data

	case addr >= 0xFF80 && addr <= 0xFFFE: // high ram (HRAM)
		newAddr := addr - 0xFF80
		data := mmu.HRAM[newAddr]
		mmu.addReadToDebug(addr, data, "HRAM")
		return data

	default:

	}
	return 0
}

// WriteByte write byte value data to location specified in address addr
func (mmu *MMU) WriteByte(addr uint16, data byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x7FFF: // Write to from Cart
		mmu.Cart.MBC.WriteByte(addr, data)
		mmu.addWriteToDebug(addr, data, "Cart")

	case addr >= 0x8000 && addr <= 0x9FFF: // Video ram
		mmu.PPU.WriteByte(addr, data)
		mmu.addWriteToDebug(addr, data, "VRAM")

	case addr >= 0xA000 && addr <= 0xBFFF: // External ram on cart
		newAddr := addr - 0xA000
		mmu.Cart.MBC.WriteByte(newAddr, data)
		mmu.addWriteToDebug(addr, data, "External Cart RAM")

	case addr >= 0xC000 && addr <= 0xDFFF: // work ram
		newAddr := addr - 0xC000
		mmu.WRAM.WriteByte(newAddr, data)
		mmu.addWriteToDebug(addr, data, "WRAM")

	case addr >= 0xFE00 && addr <= 0xFE9F: // OAM
		mmu.PPU.WriteByte(addr, data)
		mmu.addWriteToDebug(addr, data, "OAM")

	case addr == 0xFF0F: // interruptsFlag
		mmu.interruptsFlag = data
		mmu.addWriteToDebug(addr, data, "IF")

	case addr == 0xFF46: // DMA OAM transfer
		mmu.DMATransfer(data)
		mmu.addWriteToDebug(addr, data, "DMA OAM Transfer")

	case addr == 0xFFFF: // interrupts enabled
		mmu.interruptEnabled = data
		mmu.addWriteToDebug(addr, data, "IEF")

	case addr >= 0xFF00 && addr <= 0xFF7F: // I/O registers
		mmu.IO.WriteByte(addr, data)
		mmu.addWriteToDebug(addr, data, "I/O")
		

	case addr >= 0xFF80 && addr <= 0xFFFE: // high ram
		mmu.HRAM[addr-0xFF80] = data
		mmu.addWriteToDebug(addr, data, "HRAM")
		
	}
}

// DMATransfer perform an OAM DMA transfer
// TODO - maybe implement the timings if needed
func (mmu *MMU) DMATransfer(data byte)  {
	var addr uint16 = uint16(data) << 8
	for i := uint16(0); i < 0xA0; i++ {
		mmu.PPU.WriteByte(0xFE00+i, mmu.ReadByte(addr + i))
	}
	mmu.PPU.LCD.PrevOAM = data
}

// addWriteToDebug add the write attempt to the DebugRecords if debugMode is on
func (mmu *MMU) addWriteToDebug(addr uint16, data uint8, location string)  {
	if !mmu.DebugMode {
		return
	}

	mmu.DebugRecords = append(mmu.DebugRecords, fmt.Sprintf("Writing %v to 0x%04X in %s", data, addr, location))
}

// addReadToDebug add write attempt to debugrecords if debugmode is on
func (mmu *MMU) addReadToDebug(addr uint16, data uint8, location string)  {
	if !mmu.DebugMode {
		return
	}

	mmu.DebugRecords = append(mmu.DebugRecords, fmt.Sprintf("Read %v from 0x%04X in %s", data, addr, location))
}
