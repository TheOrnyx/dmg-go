package main

// WorkRAM struct for holding the information about the Work RAM
type WorkRAM struct {
	Bytes [WorkRAMSize]byte
}

// Read read the WorkRAM in a specific adress
func (w *WorkRAM) Read(address uint16) byte {
	return w.Bytes[address]
}

// Write write data to a specific adress in the workRAM
func (w *WorkRAM) Write(address uint16, data byte)  {
	w.Bytes[address] = data
}


// VideoRAM struct for holding the information about the Video RAM
type VideoRAM struct {
	Bytes [VideoRAMSize]byte
}
