package cartridge

import (
	"io"
	"log"
)

type MBC0 struct {
	romBank []byte // MBC0 is only ROM so pretty simple
}

// LoadFile implements MemoryBankController.
func (*MBC0) LoadFile(file io.Reader) error {
	return nil
}

// SaveFile implements MemoryBankController.
func (*MBC0) SaveFile(file io.Writer) error {
	return nil
}

// NewMBC0 create and return a new MBC0
func NewMBC0(rom []byte) *MBC0 {
	newMBC0 := new(MBC0)
	newMBC0.romBank = rom[0x0000:0x8000]

	return newMBC0
}

// ReadByte Read byte at given address and return it
func (m *MBC0) ReadByte(addr uint16) byte {
	if addr < 0 || addr >= 0x8000 {
		log.Fatalf("Bank attempted to read from address %v out of range", addr) // TODO - check if this needs to be fatal
	}

	return m.romBank[addr]
}

// WriteByte write given data to addr
func (m *MBC0) WriteByte(addr uint16, data byte) {
	// log.Println("Unable to write to MBC0, no RAM to write to")
}

// switchROMBank does nothing for MBC0
func (m *MBC0) switchROMBank(bank int) {

}

// switchRAMBank does nothing for MBC0 - here to satisfy interface conditions
func (m *MBC0) switchRAMBank(bank int) {

}

// HasBattery return whether or not MBC supports battery
func (m *MBC0) HasBattery() bool {
	return false
}
