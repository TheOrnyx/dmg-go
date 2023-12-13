package cartridge

type MemoryBankController interface {
	ReadByte(addr uint16) byte
	WriteByte(addr uint16, data byte)
	switchRAMBank(bank int)
	switchROMBank(bank int)
}

// createROMBanks create and populate banknum amount of ROM banks from romData and then return
// NOTE - not too sure about this one, kinda just ripped it off from the gomeboy project
func createROMBanks(romData []byte, bankNum int) [][]byte {
	romBanks := make([][]byte, bankNum)

	romBanks[0] = romData[0x4000:0x8000] // TODO - check this
	chunk := 0x4000

	for i := 1; i < bankNum; i++ {
		romBanks[i] = romData[chunk : chunk+0x4000]
		chunk += 0x4000
	}
	
	return romBanks
}

// createRAMBanks create bankNum amount of RAM Banks and populate them before returning
func createRAMBanks(bankNum int) [][]byte {
	ramBanks := make([][]byte, bankNum)

	for i := range ramBanks {
		ramBanks[i] = make([]byte, 0x2000)
	}

	return ramBanks
}
