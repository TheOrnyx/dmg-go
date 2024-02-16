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

// CombineLine combine the 3 layers from line line together and add them to the line in FinalScreen
func (s *Screen) CombineLine(line int)  {
	var screen [160]Pixel

	for x := range 160 {
		winPixel := s.Window[line][x]
		bgPixel := s.Background[line][x]
		objPixel := s.Objects[line][x]

		screen[x] = bgPixel

		if winPixel.Opaque {
			screen[x] = winPixel
		}
		
		if objPixel.Opaque { // for the sprite pixels
			screen[x] = objPixel
		}
	}
	s.FinalScreen[line] = screen
}

// Pixel - a struct to hold pixel data
type Pixel struct {
	Color   byte // color number for the pixel
	Palette byte // Value for which palette to use
	Opaque bool // whether or not to draw that pixel (true = drawn)
}
