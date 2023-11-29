package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Context struct {
	Window *sdl.Window
	Renderer *sdl.Renderer
	Screen *Screen
}

// StartSDLWindowSystem initialize and start running the sdl windowsystem
func StartSDLWindowSystem(s *Screen, width, height int32)  {
	c := new(Context)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		FatalLog.Println("Failed to Init SDL: ", err)
	}
	defer sdl.Quit()

	win, err := sdl.CreateWindow("Gameboy-Golor", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		FatalLog.Println("Failed to Create SDL window: ", err)
	}
	c.Window = win
	defer c.Window.Destroy()

	renderer, err := sdl.CreateRenderer(c.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		FatalLog.Println("Failed to creeate SDL renderer: ", err)
	}
	c.Renderer = renderer
	defer c.Renderer.Destroy()
	if err = c.Renderer.SetScale(WinScalar, WinScalar); err != nil {
		WarnLog.Println("Failed to scale up Renderer, continuing anyway but might not function")
	}

	InfoLog.Println("SDL initialized succesfully, beginning main loop")
	c.mainLoop()
}

// mainLoop run the mainloop
func (c *Context) mainLoop()  {
	running := true
	for running {
		c.clearScreen()
		c.renderScreen()
		c.Renderer.Present()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				InfoLog.Println("Exiting SDL window...")
				running = false
				break
			}
		}
	}
}

// clearScreen clears the screen
func (c *Context) clearScreen()  {
	c.Renderer.SetDrawColor(0, 0, 0, 255)
	c.Renderer.Clear()
}

// renderScreen render the gameboy screen to the sdl Window
func (c *Context) renderScreen()  {
	//put stuff here to actuall render the screen later
	c.Renderer.SetDrawColor(212, 202, 25, 255)
	c.Renderer.FillRect(&sdl.Rect{X: 32, Y: 32, W: 32, H: 32})
	c.Renderer.FillRect(&sdl.Rect{X: 96, Y: 32, W: 32, H: 32})
	c.Renderer.FillRect(&sdl.Rect{X: 0, Y: 64, W: 16, H: 64})
	c.Renderer.FillRect(&sdl.Rect{X: 144, Y: 64, W: 16, H: 64})
	c.Renderer.FillRect(&sdl.Rect{X: 16, Y: 112, W: 128, H: 16})
}
