package joypad

// the button and dpad constants
const (
	ButtonA = iota
	ButtonB
	ButtonSel
	ButtonStart
	DpadRight
	DpadLeft
	DpadUp
	DpadDown
)

const (
	DpadMode = 0x20
	ButtonMode = 0x10
)

type Joypad struct {
	keys [2]byte // the two sets of keypressed (0 = dpad, 1 = buttons)
	keyMode byte // which set of keys to use
	
	requestInterrupt func(code byte) // the function pointer to request joypad interrupt
}

// NewJoypad create a new joypad
func NewJoypad(requestInterrupt func(code byte)) *Joypad  {
	return &Joypad{requestInterrupt: requestInterrupt}
}

// ResetInput reset bits 0-3
func (j *Joypad) ResetInput()  {
	j.keys[0], j.keys[1] = 0x0F, 0x0F
	j.keyMode = 0x00
}

// HandleInput handle input given by the emulator (basically set the bits)
func (j *Joypad) HandleInput(buttons [8]bool)  {
	for i, state := range buttons {
		if state {
			j.keyDown(i)
		} else {
			j.keyUp(i)
		}
	}
}

// keyDown handle a keydown for specific key
func (j *Joypad) keyDown(key int)  {
	switch key {
	case ButtonA:
		if j.keys[1] & 0x01 == 0x1 {
			j.requestInterrupt(0x04)
		}
		j.keys[1] &= 0xE
	case ButtonB:
		if (j.keys[1]) & 0x02 == 0x02 {
			j.requestInterrupt(0x04)
		}
		j.keys[1] &= 0xD
	case ButtonSel:
		if j.keys[1] & 0x04 == 0x04 {
			j.requestInterrupt(0x04)
		}
		j.keys[1] &= 0xB
	case ButtonStart:
		if j.keys[1] & 0x08 == 0x08 {
			j.requestInterrupt(0x04)
		}
		j.keys[1] &= 0x7
	case DpadRight:
		if j.keys[0] & 0x01 == 0x1 {
			j.requestInterrupt(0x04)
		}
		j.keys[0] &= 0xE
	case DpadLeft:
		if j.keys[0] & 0x02 == 0x02 {
			j.requestInterrupt(0x04)
		}
		j.keys[0] &= 0xD
	case DpadUp:
		if j.keys[0] & 0x04 == 0x04 {
			j.requestInterrupt(0x04)
		}
		j.keys[0] &= 0xB
	case DpadDown:
		if j.keys[0] & 0x08 == 0x08 {
			j.requestInterrupt(0x04)
		}
		j.keys[0] &= 0x7
	}
}

// keyUp perform keyup
func (j *Joypad) keyUp(key int)  {
	switch key {
	case ButtonA:
		j.keys[1] |= 0x01
	case ButtonB:
		j.keys[1] |= 0x02
	case ButtonSel:
		j.keys[1] |= 0x04
	case ButtonStart:
		j.keys[1] |= 0x08
	case DpadRight:
		j.keys[0] |= 0x01
	case DpadLeft:
		j.keys[0] |= 0x02
	case DpadUp:
		j.keys[0] |= 0x04
	case DpadDown:
		j.keys[0] |= 0x08
	}
}

// Read read the joypad
func (j *Joypad) Read() byte {
	switch j.keyMode {
	case DpadMode:
		return j.keys[0] | DpadMode
	case ButtonMode:
		return j.keys[1] | ButtonMode
	default:
		return 0xF
	}
	
}

// WriteData write to the joypad
// NOTE - I assume this only ever writes to bit 4 and 5? 
func (j *Joypad) WriteData(data byte)  {
	j.keyMode = data & 0x30
}
