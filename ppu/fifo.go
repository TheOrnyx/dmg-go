package ppu

// The fetcher state constants
const (
	ReadTileID byte = iota
	ReadTileDataLow
	ReadTileDataHigh
	PushToFIFO
)

// Fetcher the pixel fetcher for the PPU to use
type Fetcher struct {
	FIFO  fifo // the fifo queue to store the pixels in
	ppu   *PPU // the PPU to allow memory reads
	ticks int  // the clock cycle counter, used for timings
	state byte // the fetchers current state, set using constants
	posX byte // the current x position of the fethcer
	msb byte // the msb for the tile data being grabbed
	lsb byte // the lsb for the tile data being grabbed
	tileIndex byte // the index for the tile being grabbed
}


// Step step the fetcher by one state
// Once again based off the article above
func (f *Fetcher) Step()  {
	// f.ticks += 1
	// if f.ticks < 2 {
	// 	return
	// }
	// f.ticks = 0 // reset the ticks and continue

	if true {
		f.stepBackground()
	}
}

// stepBackground performs a step for background pixels
func (f *Fetcher) stepBackground()  {
	switch f.state {
	case ReadTileID:
		newX := ((f.ppu.LCD.SCX/8) + f.posX) & 0x1F
		newY := (f.ppu.LCD.LY + f.ppu.LCD.SCY) & 0xFF
		f.tileIndex = f.ppu.getBackgTileIndex(newX, newY)
		
		f.state = ReadTileDataLow

	case ReadTileDataLow:
		offset := f.tileIndex + (2 * ((f.ppu.LCD.LY + f.ppu.LCD.SCY) % 8))
		f.lsb = f.ppu.ReadByte(f.ppu.tileData + uint16(offset)) // TODO - replace with signed stuff for other address modes
		
		f.state = ReadTileDataHigh
		
	case ReadTileDataHigh:
		offset := f.tileIndex + (2 * ((f.ppu.LCD.LY + f.ppu.LCD.SCY) % 8))
		f.msb = f.ppu.ReadByte(f.ppu.tileData + uint16(offset) + 1) // TODO - replace with signed stuff for other address modes
		// TODO - apparently there can be a push here or smth?
		f.state = PushToFIFO

	case PushToFIFO:
		if len(f.FIFO.Pixels) != 0 {
			return
		}

		f.FIFO.Pixels = bytesToPixels(f.msb, f.lsb)
		f.posX += 1
		f.state = ReadTileID
	}
}

// bytesToPixels convert two bytes to a row of 8 pixels
func bytesToPixels(msb, lsb byte) []Pixel {
	pixels := make([]Pixel, 8)
	for i := byte(7); i > 0; i-- {
		msbBit := (msb >> i & 0x01) << 1
		lsbBit := lsb >> i & 0x01
		pixels[i] = Pixel{Color: msbBit | lsbBit}
	}
	return pixels
}


type fifo struct {
	Pixels []Pixel // the "queue" for the pixels
}

// clear clears the fifo queue
func (f *fifo) clear()  {
	f.Pixels = []Pixel{}
}
