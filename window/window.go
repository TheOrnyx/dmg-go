package window

import (
	"log"
	"os"

	"github.com/TheOrnyx/gameboy-golor/ppu"
	"github.com/veandco/go-sdl2/sdl"
)

var FatalLog = log.New(os.Stdout, "[FaTAL] ", log.LstdFlags)

const (
	gbScreenWidth = 160
	gbScreenHeight = 144
)

type Context struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
}

var Palette [4]sdl.Color = [4]sdl.Color{
	sdl.Color{R: 155, G: 188, B: 15, A: 255},
	sdl.Color{R: 139, G: 172, B: 15, A: 255},
	sdl.Color{R: 48, G: 98, B: 48, A: 255},
	sdl.Color{R: 15, G: 56, B: 15, A: 255},
}

// StartSDLWindowSystem initialize and start running the sdl windowsystem
func InitSDLWindowSystem(width, height, scale int32) *Context {
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

// // mainLoop run the mainloop
// func (c *Context) mainLoop() {
// 	running := true
// 	for running {
// 		c.ClearScreen()
// 		c.RenderScreen()
// 		c.Renderer.Present()
// 		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
// 			switch event.(type) {
// 			case *sdl.QuitEvent:
// 				InfoLog.Println("Exiting SDL window...")
// 				running = false
// 				break
// 			}
// 		}
// 	}
// }

// ClearScreen clears the screen
func (c *Context) ClearScreen() {
	c.Renderer.SetDrawColor(0, 0, 0, 255)
	c.Renderer.Clear()
}

// RenderScreen render the gameboy screen to the sdl Window
func (c *Context) RenderScreen(screen *ppu.Screen) {
	for drawY := 0; drawY < gbScreenHeight; drawY++ {
		for drawX := 0; drawX < gbScreenWidth; drawX++ {
			
			c.Renderer.SetDrawColor(colorsFromSDLCol(Palette[screen.FinalScreen[drawY][drawX].Color]))
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
