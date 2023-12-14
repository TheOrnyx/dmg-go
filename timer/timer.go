package timer

import "log"

type Timer struct {
	div                    uint16 // the div counter - only upper 8-bits can be read by cpu
	tima                   byte // the timer counter
	tma                    byte   // value to set timaCounter to when it overflows
	tac                    byte   // timer control - controls behaviour of the TIMA reg
	lastBit                uint16  // the last and result of the chose DIV value and the timer enable bit, used to detect falling edge
	overflowed             bool   // whether or not timer has overflowed
	doubleSpeed            bool   // whether or not the timer is running on double speed - TODO IMPLEMENT
	timaReload             bool   // whether or not you're in the process of reloading the TIMA
	cyclesTilTIMAInterrupt int    // number of cycles until the TIMA IRQ interrupt flag is raised
}

// Tick tick timer by cycles amnt of m-cycles
// NOTE - only increments by one at the moment as I haven't adapted it
// TODO - probably switch to just inc by 1 and call this whenever it's needed
// TODO - add a return for overflow or smth - I like the function pointer idea but will see
// TODO - check out that silly thing about DIV only being like 14 bits (https://discord.com/channels/465585922579103744/465586075830845475/1184659262421618729)
func (t *Timer) Tick(cycles int) {
	// t.div += uint16(cycles*4)
	// timerEnabled := t.tac & 0x04 == 0x04

	t.timaReload = false
	if t.cyclesTilTIMAInterrupt > 0 {
		t.cyclesTilTIMAInterrupt -= 1
		if t.cyclesTilTIMAInterrupt == 0 {
			// TODO - add bit for raising the interrupt flag
			t.overflowed = true
			t.tima = t.tma
			t.timaReload = true
		}
	}
	t.changeDiv(t.div + 1)
}

// changeDiv change the value of the div and also adjust TMA accordingly
// Also check things like falling edges to find out whether to increase TIMA
// Thanks to https://github.com/raddad772/jsmoo/blob/main/system/gb/gb_cpu.js#L4 for providing a good example
func (t *Timer) changeDiv(newVal uint16)  {
	t.div = newVal
	var chosenBit uint16 // the chosen bit to use from the DIV based on the TAC
	// TODO - check this is correct
	switch t.tac & 0x03 {
	case 0: // use bit 9
		chosenBit = (t.div >> 9) & 0x01
	case 1: // use bit 3
		chosenBit = (t.div >> 3) & 0x01
	case 2: // use bit 5
		chosenBit = (t.div >> 5) & 0x01
	case 3: // use bit 7
		chosenBit = (t.div >> 7) & 0x01
	}

	timerEnabled := (t.tac & 0x04) >> 2 // get the third bit of the tac (the timer enabled bit)
	chosenBit &= uint16(timerEnabled) // TODO - check

	t.detectFallingEdge(t.lastBit, chosenBit)
	t.lastBit = chosenBit	
}

// detectFallingEdge detect prescense of a falling edge and handle it if one exists
func (t *Timer) detectFallingEdge(oldVal, newVal uint16)  {
	if oldVal == 1 && newVal == 0 {
		t.tima = (t.tima + 1) & 0xFF // increment TIMA and mask to 0xFF to detect overflow
		if t.tima == 0 { // detect overflow and schedule interruption if it overflowed
			t.cyclesTilTIMAInterrupt = 1
		}
	}
}

// Write write data to the timer - handle the address options
func (t *Timer) Write(addr uint16, data byte)  {
	switch addr {
	case 0xFF04: // the DIV register - clears it
		t.changeDiv(0)
	case 0xFF05: // TIMA register
		if !t.timaReload {
			t.tima = data
		}
		if t.cyclesTilTIMAInterrupt == 1 { // I think this is to prevent interrupt flag from being set and stuff
			t.cyclesTilTIMAInterrupt = 0
		}
	case 0xFF06: // TMA register
		if t.timaReload {
			t.tima = data
		}
		t.tma = data
	case 0xFF07: // TAC register
		lastBit := t.lastBit
		timerEnabled := (t.tac & 0x04) >> 2
		t.lastBit &= uint16(timerEnabled)

		t.detectFallingEdge(lastBit, t.lastBit)
		t.tac = data
	}
}

// Read read and return the data at addr 
func (t *Timer) Read(addr uint16) byte {
	switch addr {
	case 0xFF04: // DIV - top 8 bits of it
		return byte(t.div >> 8 & 0xFF)
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	default:
		log.Println("Unkown Timer read address:", addr)
		return 0
	}
}
