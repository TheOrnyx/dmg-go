package cartridge

import "log"

////////////
// TODO's //
////////////
//
// TODO - check whatever the hell a "mode flag" and a "Zero bank" is (https://hacktix.github.io/GBEDG/mbcs/mbc1/#zero-bank)
// TODO - possibly implement the like large banks which has stuff like the mode flag

const (
	sixteenMBRom8KBRam = iota
	fourMBRom32KBRam
)

type MBC1 struct {
	name string
	mode uint8 // The banking mode (0 for default/ 16mbROM with 8kbRAM - 1 for 4mbROM with 32KBRAM) - TODO CHECK THIS
	hasBattery bool

	romBank0 []byte // contains first 16kib of the cart ROM
	romBanks [][]byte
	romBank  int // the active rom Bank
	ROMSize  int // the rom size

	ramBanks   [][]byte
	ramBank    int  // the active RAM bank
	hasRAM     bool // whether or not the MBC1 has ram or not (cuz like there's some that don't have ram)
	ramEnabled bool // bool for if RAM is enabled or not
	RAMSize    int  // the ram size
}

// NewMBC1 create a new MBC1 from specifications
func NewMBC1(rom []byte, romSize, ramSize int, hasBattery bool) *MBC1 {
	newMBC := new(MBC1)

	newMBC.name = "CART-MBC1"
	newMBC.mode = sixteenMBRom8KBRam
	newMBC.hasBattery = hasBattery
	newMBC.ROMSize, newMBC.RAMSize = romSize, ramSize

	if ramSize > 0 { // enable ram stuff if ram supported
		newMBC.hasRAM = true
		newMBC.ramEnabled = true
		newMBC.ramBank = 0
		newMBC.ramBanks = createRAMBanks(4)
	}

	newMBC.romBank = 0
	newMBC.romBank0 = rom[0x0000:0x4000]
	newMBC.romBanks = createROMBanks(rom, romSize/0x4000)

	return newMBC
}

// ReadByte read Byte from addr
func (m *MBC1) ReadByte(addr uint16) byte {
	if addr <= 0x3FFF {
		return m.romBank0[addr] // TODO - check the zero bank thing (top of page)
	}

	if addr >= 0x4000 && addr <= 0x7FFF {
		return m.romBanks[m.romBank][addr-0x4000]
	}

	if addr >= 0xA000 && addr <= 0xBFFF {
		if m.hasRAM && m.ramEnabled {
			switch m.mode {
			case fourMBRom32KBRam:
				return m.ramBanks[m.ramBank][addr - 0xA000]
			case sixteenMBRom8KBRam:
				return m.ramBanks[0][addr - 0xA000]
			}
		}
	}

	return 0xFF // NOTE - I think it's meant to be 0xFF? either that or 0x00
}

// WriteByte write given data to addr
func (m *MBC1) WriteByte(addr uint16, data byte)  {
	switch {
	case addr <= 0x1FFF && m.hasRAM: // enable or disable RAM
		if data & 0x0F == 0x0A { // TODO - check that I don't need to check the mode
			log.Printf("%s: Enabling RAM...\n", m.name)
			m.ramEnabled = true
		} else {
			log.Printf("%s: DIsabling RAM...\n", m.name)
			m.ramEnabled = false
		}

	case addr >= 0x2000 && addr <= 0x3FFF: // switch ROM bank
		m.switchROMBank(int(data & 0x1F)) // TODO - check if I need to do any other masking for sizes
		

	case addr >= 0x4000 && addr <= 0x5FFF: // switch RAM bank
		m.switchRAMBank(int(data & 0x03))
		

	case addr >= 0x6000 && addr <= 0x7FFF: // mode select
		m.mode = data & 0x1
		// TODO - check

	case addr >= 0xA000 && addr <= 0xBFFF: // External RAM - TODO check this is right
		if m.hasRAM && m.ramEnabled {
			switch m.mode {
			case fourMBRom32KBRam:
				m.ramBanks[m.ramBank][addr-0xA000] = data
			case sixteenMBRom8KBRam:
				m.ramBanks[0][addr-0xA000] = data
			}
		}
	}
}

// switchROMBank switch active rom bank to val
func (m *MBC1) switchROMBank(bank int)  {
	m.romBank = bank
}

// switchRAMBank switch active ram bank to val
func (m *MBC1) switchRAMBank(bank int)  {
	m.ramBank = bank
}
