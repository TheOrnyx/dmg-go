package ppu

// Screen - a struct to hold the pixel data for each of the scanlines
type Screen struct {
	Pixels [160][144]Pixel // the pixels for the screen
}

// Pixel - a struct to hold pixel data
type Pixel struct {
	Color byte // color number for the pixel
	Palette byte // Value for which palette to use
}
