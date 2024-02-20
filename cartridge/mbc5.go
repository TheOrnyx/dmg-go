package cartridge

type MBC5 struct {
	romSize int // the rom size
	romBanks [][]byte // the rom banks (allows for 512 of them and no individual bank0 field as bank0 is here)
	romBank uint16 // the active rom bank
	hasRAM bool // whether or not cart has RAM
	ramSize int // the RAM size (if any)
	ramBanks [][]byte // the ram banks
	ramBank byte // the active ram bank
	ramEnabled bool // whether or not the ram is enabled
	hasBattery bool // whether or not cart has battery
}

// NewMBC5 create and return a new MBC5
func NewMBC5(rom []byte, hasBattery bool, ramSize, romSize int) *MBC5 {
	mbc := new(MBC5)
	mbc.hasBattery = hasBattery
	mbc.romSize, mbc.ramSize = romSize, ramSize

	if ramSize > 0 {
		mbc.hasRAM = true
		mbc.ramEnabled = true
		mbc.ramBank = 0
		mbc.ramBanks = createRAMBanks(16) // TODO - maybe change to use ramSize
	}

	mbc.romBank = 0
	mbc.romBanks = createROMBanksFromZero(rom, romSize/0x4000)

	return mbc
}

// ReadByte read byte from memory
func (m *MBC5) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // rom bank 0
		return m.romBanks[0][addr]
	case addr >= 0x4000 && addr <= 0x7FFF: // Switchable rom bank
		return m.romBanks[m.romBank][addr - 0x4000]
	case addr >= 0xA000 && addr <= 0xBFFF: // switchable ram bank (if enabled)
		if !m.hasRAM || !m.ramEnabled {
			return 0xFF
		}

		return m.ramBanks[m.ramBank][addr - 0xA000]
	}

	return 0xFF
}

// WriteByte write given data to addr
func (m *MBC5) WriteByte(addr uint16, data byte)  {
	switch {
	case addr >= 0x0000 && addr <= 0x1FFF: // Ram enable
		if !m.hasRAM {
			return
		}

		if data & 0x0F == 0x0A {
			m.ramEnabled = true
		} else {
			m.ramEnabled = false
		}

	case addr >= 0x2000 && addr <= 0x2FFF: // set lower 8 bits of ROM Bank value to data
		m.romBank = (m.romBank & 0xFF00) | uint16(data) // TODO - check if this is right

	case addr >= 0x3000 && addr <= 0x3FFF: // set 9th bit of rom bank value to data's lowest bit (bit 0)
		m.romBank = (m.romBank & 0xFF) | ((uint16(data)&0x01) << 8)

	case addr >= 0x4000 && addr <= 0x5FFF: // set ram bank value
		m.ramBank = data

	case addr >= 0xA000 && addr <= 0xBFFF: // write to external ram
		if !m.hasRAM || !m.ramEnabled {
			return
		}

		m.ramBanks[m.ramBank][addr - 0xA000] = data
	}
}


// SwitchRAMBank switch ram bank
func (m *MBC5) switchRAMBank(bank int)  {
	
}

// SwitchROMBank switch ram bank
func (m *MBC5) switchROMBank(bank int)  {
	
}
