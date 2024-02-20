package cartridge

///////////////
// MBC2 Cart //
///////////////
//
// 256KB ROM
// 512x4 RAM

type MBC2 struct {
	name string

	romBank0 []byte // contains first 16kib of ROM
	romBanks[][]byte // the other rom banks
	romBank byte // Currently selected ROM bank number
	externalRam [512]byte // the internal RAM on the cart (is represented in half-bytes, so only bottom 4 bytes are used)
	RamEnabled bool // whether or not RAM is enabled
}

// NewMBC2 create, map and return a new MBC2
func NewMBC2(rom []byte, romSize int) *MBC2 {
	mbc := new(MBC2)
	mbc.name = "CART-MBC2"
	mbc.romBank0 = rom[0x0000:0x4000]
	mbc.romBanks = createROMBanks(rom, romSize/0x4000)
	mbc.romBank = 1
	return mbc
}

// switchRAMBank UNUSED, done in regular functions
func (m *MBC2) switchRAMBank(bank int)  {
	
}

// switchROMBank UNUSED, done in regular functions
func (m *MBC2) switchROMBank(bank int)  {
	
}

// WriteByte write given byte to addr location
func (m *MBC2) WriteByte(addr uint16, data byte)  {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // enable RAM / ROM bank number
		if (addr >> 8)&0x01 == 0x01 { // switch ROM bank num (can't be 0)
			m.romBank = data&0x0F
		} else {
			if data & 0x0F == 0xA {
				m.RamEnabled = true
			} else {
				m.RamEnabled = false
			}
		}

	case addr >= 0xA000 && addr <= 0xBFFF: // eternal ram
		if !m.RamEnabled {
			return
		}

		m.externalRam[addr & 0x1FF] = data & 0x0F
	}
}

// ReadByte read the byte from the MBC2 memory map
func (m *MBC2) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // ROM Bank 0
		return m.romBank0[addr]
	case addr >= 0x4000 && addr <= 0x7FFF: // selectable rom bank (mapped using 0x4000 * RomBankNum + (addr - 0x4000))
		return m.romBanks[m.romBank][(addr - 0x4000)]
	case addr >= 0xA000 && addr <= 0xBFFF: // external RAM
		if !m.RamEnabled {
			return 0xFF
		}
		return m.externalRam[addr & 0x1FF] | 0xF0
	}
	return 0xFF
}
