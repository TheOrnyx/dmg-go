package cartridge

import (
	"encoding/binary"
	"io"
)

type MemoryBankController interface {
	ReadByte(addr uint16) byte
	WriteByte(addr uint16, data byte)
	switchRAMBank(bank int)
	switchROMBank(bank int)
	HasBattery() bool // return whether or not MBC has battery support
	SaveFile(file io.Writer) error
	LoadFile(file io.Reader) error
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

// createROMBanksFromZero same as createROMBanks but writes the
// 0x0000-0x4000 to the first bank
func createROMBanksFromZero(romData []byte, bankNum int) [][]byte {
	romBanks := make([][]byte, bankNum)

	romBanks[0] = romData[0x0000:0x4000] // TODO - check this
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

// writeRAMToFile write ram array to file
func writeRAMToFile(ram [][]byte, file io.Writer) error {
	for bank := range ram {
		if err := binary.Write(file, binary.LittleEndian, ram[bank]); err != nil {
			return err
		}
	}
	return nil
}

// readRamFromFile read the ram data from given file and return the byte arrays
func readRamFromFile(file io.Reader, bankNum, bankSize int) ([][]byte, error) {
	// reader := bufio.NewReader(file)
	newBanks := make([][]byte, bankNum)
	for i := 0; i < bankNum; i++ {
		bank := make([]byte, bankSize)
		if err := binary.Read(file, binary.LittleEndian, &bank); err != nil {
			return nil, err
		}
		newBanks[i] = bank
	}
	return newBanks, nil
}
