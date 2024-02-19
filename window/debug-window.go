package window

import (
	"log"

	"github.com/TheOrnyx/gameboy-golor/joypad"
	"github.com/TheOrnyx/gameboy-golor/ppu"
	"github.com/veandco/go-sdl2/sdl"
)

type DebugWindow struct {
	Window *sdl.Window
	Renderer *sdl.Renderer
	midX int32
	midY int32
	width int32
	height int32
}

// StartSDLWindowSystem initialize and start running the sdl windowsystem
func CreateDebugWindow(width, height, scale int32) Screen {
	d := new(DebugWindow)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		FatalLog.Println("Failed to Init SDL: ", err)
	}

	win, err := sdl.CreateWindow("Gameboy-Golor debug", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		FatalLog.Println("Failed to Create SDL window: ", err)
	}
	d.Window = win

	renderer, err := sdl.CreateRenderer(d.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		FatalLog.Println("Failed to creeate SDL renderer: ", err)
	}
	d.Renderer = renderer
	if err = d.Renderer.SetScale(float32(scale)/2, float32(scale)/2); err != nil {
		log.Println("Failed to scale up Renderer, continuing anyway but might not function")
	}

	d.width, d.height = width/2, height/2
	d.midX, d.midY = d.width/2, d.height/2


	return d
}

// ClearScreen clears the screen
func (d *DebugWindow) ClearScreen() {
	d.Renderer.SetDrawColor(0, 0, 0, 255)
	d.Renderer.Clear()
}

// DebugRender render the debug screen with each layer view
func (d *DebugWindow) RenderScreen(screen *ppu.Screen)  {
	
	// Background layer
	for drawY := 0; drawY < int(d.midY); drawY++ { 
		for drawX := 0; drawX < int(d.midX); drawX++ {
			d.Renderer.SetDrawColor(colorsFromSDLCol(mainPalette[screen.Background[drawY][drawX].Color]))
			d.Renderer.DrawPoint(int32(drawX), int32(drawY))
		}
	}

	// Window layer
	for drawY := 0; drawY < int(d.midY); drawY++ {
		for drawX := d.midX; drawX < d.width; drawX++ {
			d.Renderer.SetDrawColor(colorsFromSDLCol(mainPalette[screen.Window[drawY][drawX - d.midX].Color]))
			d.Renderer.DrawPoint(int32(drawX), int32(drawY))
		}
	}

	for drawY := d.midY; drawY < d.height; drawY++ {
		for drawX := 0; drawX < int(d.midX); drawX++ {
			d.Renderer.SetDrawColor(colorsFromSDLCol(mainPalette[screen.Objects[drawY - d.midY][drawX].Color]))
			d.Renderer.DrawPoint(int32(drawX), int32(drawY))
		}
	}

	for drawY := d.midY; drawY < d.height; drawY++ {
		for drawX := d.midX; drawX < d.width; drawX++ {
			d.Renderer.SetDrawColor(colorsFromSDLCol(mainPalette[screen.FinalScreen[drawY - d.midY][drawX - d.midX].Color]))
			d.Renderer.DrawPoint(int32(drawX), int32(drawY))
		}
	}

	d.Renderer.Present()
}

// CloseScreen shut down the screen
func (d *DebugWindow) CloseScreen()  {
	sdl.Quit()
	d.Window.Destroy()
	d.Renderer.Destroy()
}

// ReceiveInput return the current list of inputs based on their activeness
// Returns bools based on an index using the joypad constants
func (d *DebugWindow) GetInput() [8]bool {
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

	return inputs
}
