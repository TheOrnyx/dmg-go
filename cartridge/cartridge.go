package cartridge

import (
	"fmt"
)

// Memory bank type constants - mapped to their equivalent value in rom[0x0147]
const (
	MBC_0                    = 0x00
	MBC_1                    = 0x01
	MBC_1_RAM                = 0x02
	MBC_1_RAM_BATTERY        = 0x03
	MBC_2                    = 0x05
	MBC_2_BATTERY            = 0x06
	MBC_3_TIMER_BATTERY      = 0x0F
	MBC_3_TIMER_RAM_BATTERY  = 0x10
	MBC_3                    = 0x11
	MBC_3_RAM                = 0x12
	MBC_3_RAM_BATTERY        = 0x13
	MBC_5                    = 0x19
	MBC_5_RAM                = 0x1A
	MBC_5_RAM_BATTERY        = 0x1B
	MBC_5_RUMBLE             = 0x1C
	MBC_5_RUMBLE_RAM         = 0x1D
	MBC_5_RUMBLE_RAM_BATTERY = 0x1E
	MBC_6                    = 0x20

	//cbf doing rest
	//ignoring the rom and ram options as they aren't used?
	// TODO add the MMM01 stuff - excluding for now as they aren't really used
)

var CartTypes map[byte]CartType = map[byte]CartType{
	MBC_0:                    CartType{MBC_0, "ROM ONLY"},
	MBC_1:                    CartType{MBC_1, "MBC1"},
	MBC_1_RAM:                CartType{MBC_1_RAM, "MBC1+RAM"},
	MBC_1_RAM_BATTERY:        CartType{MBC_1_RAM_BATTERY, "MBC1+RAM+BATTERY"},
	MBC_2:                    CartType{MBC_2, "MBC2"},
	MBC_2_BATTERY:            CartType{MBC_2_BATTERY, "MBC2+BATTERY"},
	MBC_3_TIMER_BATTERY:      CartType{MBC_3_TIMER_BATTERY, "MBC3+TIMER+BATTERY"},
	MBC_3_TIMER_RAM_BATTERY:  CartType{MBC_3_TIMER_RAM_BATTERY, "MBC3+TIMER+RAM+BATTERY"},
	MBC_3:                    CartType{MBC_3, "MBC3"},
	MBC_3_RAM:                CartType{MBC_3_RAM, "MBC3+RAM"},
	MBC_3_RAM_BATTERY:        CartType{MBC_3_RAM_BATTERY, "MBC3+RAM+BATTERY"},
	MBC_5:                    CartType{MBC_5, "MBC5"},
	MBC_5_RAM:                CartType{MBC_5_RAM, "MBC5+RAM"},
	MBC_5_RAM_BATTERY:        CartType{MBC_5_RAM_BATTERY, "MBC5+RAM+BATTERY"},
	MBC_5_RUMBLE:             CartType{MBC_5_RUMBLE, "MBC5+RUMBLE"},
	MBC_5_RUMBLE_RAM:         CartType{MBC_5_RUMBLE_RAM, "MBC5+RUMBLE+RAM"},
	MBC_5_RUMBLE_RAM_BATTERY: CartType{MBC_5_RUMBLE_RAM_BATTERY, "MBC5+RUMBLE+RAM+BATTERY"},
	MBC_6:                    CartType{MBC_6, "MBC6"},
}

type CartType struct {
	ID   byte
	Desc string
}

type Cartridge struct {
	Title string
	ROM           []byte
	RAMSize       int
	ROMSize       int
	HasCGBSupport bool // whether the CGB flag is set -> fonud in rom[0x0143]
	RomType       byte
	Type          CartType

	MBC MemoryBankController // the memory bank controller
	
	IsJapanese bool
	OldLicenseeCode byte // the old licensee code, if 33 then use new licensee code
	NewLicenseeCode byte // the new licensee code, only used if OldLicenseeCode is 33
	
}

// LoadROM load and initialize a ROM based on cart path
func LoadROM(rom []byte) (*Cartridge, error) {
	newCart := new(Cartridge)
	newCart.ROM = rom // TODO - check this?
	
	newCart.InitCart(rom)
	
	return newCart, nil
}

// InitCart initialize the cart
func (c *Cartridge) InitCart(rom []byte) error {

	c.HasCGBSupport = rom[0x0143] == 0x80 || rom[0x0143] == 0xC0
	c.Title = string(rom[0x134:0x0143]) //TODO - make this actually good
	romSize := rom[0x0148]
	c.ROMSize = 0x8000 << romSize // TODO maybe error handle here?

	switch rom[0x0149] {
	case 0x00:
		c.RAMSize = 0
	case 0x01:
		c.RAMSize = 2048
	case 0x02:
		c.RAMSize = 8192
	case 0x03:
		c.RAMSize = 32678
	case 0x04:
		c.RAMSize = 131072
	}

	cType, found := CartTypes[rom[0x0147]]
	if !found {
		return fmt.Errorf("Unkown cart type found at %v", rom[0x0147])
	}
	c.Type = cType

	c.IsJapanese = rom[0x014A] == 0x00

	c.OldLicenseeCode = rom[0x014B]
	// Check for new licensee code later

	switch c.Type.ID {
	case MBC_0:
		c.MBC = NewMBC0(rom)
	case MBC_1, MBC_1_RAM:
		c.MBC = NewMBC1(rom, c.ROMSize, c.RAMSize, false)
	case MBC_1_RAM_BATTERY:
		c.MBC = NewMBC1(rom, c.ROMSize, c.RAMSize, true)
	}

	return nil
}
