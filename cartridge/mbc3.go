package cartridge

import "io"

type MBC3 struct {
	romSize    int      // the rom size
	romBanks   [][]byte // the rom banks (max 128)
	romBank0   []byte   // rom bank 0 (first 16KIB of rom)
	romBank    byte     // the current rom bank (7 bits)
	hasRam     bool     // whether or not cart supports ram
	ramSize    int      // the ram size
	ramBanks   [][]byte // the RAM banks (max 4)
	ramBank    byte     // the current RAM bank (2 bits)
	ramEnabled bool     // whether or not ram and RTC are enabled
	hasTimer   bool     // whether or not MBC3 has a timer
	hasBattery bool     // whether or not has battery
	rtcMapped  bool     // if true then RTC is mapped to 0xA000 - 0xBFFF otherwise ram is
	rtc        *RTC     // the RTC (real time clock) for the cart
}

// LoadFile implements MemoryBankController.
func (m *MBC3) LoadFile(file io.Reader) error {
	if !m.hasBattery {
		return nil
	}
	var err error = nil

	m.ramBanks, err = readRamFromFile(file, len(m.ramBanks), 0x2000)
	return err
}

// SaveFile implements MemoryBankController.
func (m *MBC3) SaveFile(file io.Writer) error {
	if !m.hasBattery {
		return nil
	}

	return writeRAMToFile(m.ramBanks, file)
}

// NewMBC3 create and return a new MBC3, populating the banks
func NewMBC3(rom []byte, hasBattery, hasTimer bool, ramSize, romSize int) *MBC3 {
	mbc := new(MBC3)
	mbc.hasBattery = hasBattery
	mbc.hasTimer = hasTimer
	mbc.romSize, mbc.ramSize = romSize, ramSize

	if ramSize > 0 { // enable ram if supported
		mbc.hasRam = true
		mbc.ramEnabled = true
		mbc.ramBank = 0
		mbc.ramBanks = createRAMBanks(4) // TODO - check it's always 4		
	}

	mbc.romBank = 0
	mbc.romBank0 = rom[0x0000:0x4000]
	mbc.romBanks = createROMBanks(rom, romSize/0x4000)

	return mbc
}

// ReadByte read byte at point in mbc3
func (m *MBC3) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // Rom Bank 0
		return m.romBank0[addr]

	case addr >= 0x4000 && addr <= 0x7FFF: // switchable rom banks
		return m.romBanks[m.romBank][addr-0x4000]

	case addr >= 0xA000 && addr <= 0xC000: // external Ram/ RTC
		if m.rtcMapped && m.hasTimer {
			return m.rtc.readByte(addr)
		}

		if m.ramEnabled && m.hasRam {
			return m.ramBanks[m.ramBank][addr-0xA000]
		}
		return 0xFF
	}

	return 0xFF
}

// WriteByte write byte to addr in MBC3
func (m *MBC3) WriteByte(addr uint16, data byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x1FFF: // enable ram/ timer registers
		if data&0x0F == 0xA {
			m.ramEnabled = true
		} else {
			m.ramEnabled = false
		}

	case addr >= 0x2000 && addr <= 0x3FFF: // ROM bank low
		m.romBank = data & 0x7F

	case addr >= 0x4000 && addr <= 0x5FFF: // Ram Bank/ RTC Reg select
		switch {
		case addr >= 0x00 && addr <= 0x03: // set RAM Banks
			if m.hasRam {
				m.ramBank = data & 0x03
				m.rtcMapped = false
			}

		case addr >= 0x08 && addr <= 0x0C: // RTC Map
			// TODO - write about mapping the given RTC register
			m.rtcMapped = true
		}

	case addr >= 0x6000 && addr <= 0x7FFF: // Latch Clock Data
		// TODO - write this

	case addr >= 0xA000 && addr <= 0xBFFF: // EXternal RAM / RTC reg write
		if !m.ramEnabled {
			return
		}

		if m.rtcMapped {
			m.rtc.writeByte(addr, data)
		} else if m.hasRam {
			m.ramBanks[m.ramBank][addr-0xA000] = data
		}
	}
}

// SwitchRAMBank switch ram bank
func (m *MBC3) switchRAMBank(bank int) {

}

// SwitchROMBank switch ram bank
func (m *MBC3) switchROMBank(bank int) {

}

// HasBattery return whether or not MBC supports battery
func (m *MBC3) HasBattery() bool {
	return m.hasBattery
}

type RTC struct {
}

// readByte read byte from RTC registers
func (r *RTC) readByte(addr uint16) byte {
	return 0xFF
}

// writeByte write byte to RTC
func (r *RTC) writeByte(addr uint16, data byte) {

}
