package ppu

// Screen - a struct to hold the pixel data for each of the scanlines
type Screen struct {
	FinalScreen [144][160]Pixel // the pixels for the screen
	Background  [144][160]Pixel // the pixels for the Background layer
	Window      [144][160]Pixel // The pixels for the Window layer
	Objects     [144][160]Pixel // the pixels for the Objects layer
}

// Reset reset the screen (usually done at the beginning of each frame)
func (s *Screen) Reset() {
	s.FinalScreen	= [144][160]Pixel{}
	s.Background	= [144][160]Pixel{}
	s.Window		= [144][160]Pixel{}
	s.Objects		= [144][160]Pixel{}
}

// Pixel - a struct to hold pixel data
type Pixel struct {
	Color   byte // color number for the pixel
	Palette byte // Value for which palette to use
}
