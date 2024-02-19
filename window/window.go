package window

import (
	"log"
	"os"

	"github.com/TheOrnyx/gameboy-golor/joypad"
	"github.com/TheOrnyx/gameboy-golor/ppu"
	"github.com/veandco/go-sdl2/sdl"
)

var FatalLog = log.New(os.Stdout, "[FaTAL] ", log.LstdFlags)

const (
	gbScreenWidth = 160
	gbScreenHeight = 144
)

type Screen interface {
	ClearScreen()
	RenderScreen(screen *ppu.Screen)
	CloseScreen()
	GetInput() [8]bool
}

type Context struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
}

var GreenPalette [4]sdl.Color = [4]sdl.Color{
	sdl.Color{R: 155, G: 188, B: 15, A: 255},
	sdl.Color{R: 139, G: 172, B: 15, A: 255},
	sdl.Color{R: 48, G: 98, B: 48, A: 255},
	sdl.Color{R: 15, G: 56, B: 15, A: 255},
}

var GrayPalette [4]sdl.Color = [4]sdl.Color{
	sdl.Color{R: 255, G: 255, B: 255, A: 255},
	sdl.Color{R: 172, G: 172, B: 172, A: 255},
	sdl.Color{R: 84, G: 84, B: 84, A: 255},
	sdl.Color{R: 0, G: 0, B: 0, A: 255},
}

var mainPalette = GreenPalette


// StartSDLWindowSystem initialize and start running the sdl windowsystem
func InitSDLWindowSystem(width, height, scale int32) Screen {
	c := new(Context)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		FatalLog.Println("Failed to Init SDL: ", err)
	}

	win, err := sdl.CreateWindow("Gameboy-Golor", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		FatalLog.Println("Failed to Create SDL window: ", err)
	}
	c.Window = win

	renderer, err := sdl.CreateRenderer(c.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		FatalLog.Println("Failed to creeate SDL renderer: ", err)
	}
	c.Renderer = renderer
	if err = c.Renderer.SetScale(float32(scale), float32(scale)); err != nil {
		log.Println("Failed to scale up Renderer, continuing anyway but might not function")
	}

	return c
}

// ClearScreen clears the screen
func (c *Context) ClearScreen() {
	c.Renderer.SetDrawColor(0, 0, 0, 255)
	c.Renderer.Clear()
}

// RenderScreen render the gameboy screen to the sdl Window
func (c *Context) RenderScreen(screen *ppu.Screen) {
	for drawY := 0; drawY < gbScreenHeight; drawY++ {
		for drawX := 0; drawX < gbScreenWidth; drawX++ {
			
			c.Renderer.SetDrawColor(colorsFromSDLCol(mainPalette[screen.FinalScreen[drawY][drawX].Color]))
			c.Renderer.DrawPoint(int32(drawX), int32(drawY))
		}
	}
	c.Renderer.Present()
}

// colorsFromSDLCol return individual colors for a given sdl color
func colorsFromSDLCol(col sdl.Color) (r, g, b, a uint8) {
	return col.R, col.G, col.B, col.A
}

// CloseScreen shut down the screen
func (c *Context) CloseScreen()  {
	sdl.Quit()
	c.Window.Destroy()
	c.Renderer.Destroy()
}

// ReceiveInput return the current list of inputs based on their activeness
// Returns bools based on an index using the joypad constants
// TODO - implement mapping
func (c *Context) GetInput() [8]bool {
	var inputs [8]bool
	sdl.PollEvent()
	keys := sdl.GetKeyboardState()
	inputs[joypad.ButtonA] = keys[sdl.SCANCODE_Z] == 1
	inputs[joypad.ButtonB] = keys[sdl.SCANCODE_X] == 1
	inputs[joypad.ButtonSel] = keys[sdl.SCANCODE_S] == 1
	inputs[joypad.ButtonStart] = keys[sdl.SCANCODE_A] == 1

	inputs[joypad.DpadRight] = keys[sdl.SCANCODE_RIGHT] == 1
	inputs[joypad.DpadLeft] = keys[sdl.SCANCODE_LEFT] == 1
	inputs[joypad.DpadUp] = keys[sdl.SCANCODE_UP] == 1
	inputs[joypad.DpadDown] = keys[sdl.SCANCODE_DOWN] == 1

	if keys[sdl.SCANCODE_R] == 1 {
		if mainPalette == GreenPalette {
			mainPalette = GrayPalette
		} else {
			mainPalette = GreenPalette
		}
	}
	
	return inputs
}
