package main

import (
	_"image"
	"image/color"
)

const ScreenWidth, ScreenHeight = 160, 144 //the width and height of the screen

type ColorPalette struct {
	Colors [32 * 32]color.RGBA
}

type Screen struct {
	Pixels [][]uint8
	Palette ColorPalette
}

// MakeNewScreen create a new screen to use
func MakeNewScreen() *Screen {
	newScreen := new(Screen)
	pixels := make([][]uint8, ScreenHeight)
	for i := range pixels {
		pixels[i] = make([]uint8, ScreenWidth)
	}

	newScreen.Pixels = pixels
	newScreen.Palette = initColorPalette()

	InfoLog.Println("Screen Succesfully set up")
	return newScreen
}


// initColorPalette initialize the color palette (TODO implement later)
func initColorPalette() ColorPalette {
	newPalette := new(ColorPalette)
	return *newPalette
}
