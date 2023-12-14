package main

import (
	"log"
	"os"
	"github.com/TheOrnyx/gameboy-golor/window"
)

const UsingSDL = true
const WinScalar = 4 //the scalar used to scale up the gbc screen
const WinWidth, WinHeight = 160 * WinScalar, 144 * WinScalar

var InfoLog = log.New(os.Stdout, "[INFO] ", log.Ldate)
var DebugLog = log.New(os.Stdout, "[DEBUG] ", log.Ldate)
var WarnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
var FatalLog = log.New(os.Stdout, "[FaTAL] ", log.LstdFlags)

func main() {
	InfoLog.Println("Starting...")
	window.StartSDLWindowSystem(WinWidth, WinHeight)
	
	InfoLog.Println("Program finished, exiting...")
}
