package cpu

// Instruction struct containing information about an instruction
// kinda plagirazed from github.com/djhworld/gomeboycolor
type Instruction struct {
	OpCode  byte           // the opcode of the instruction, in hex form usually
	Desc    string         // basically just the name of the instruction, used for debug purposes
	Size    int            // size of the instruction aka the amount of operands it has
	Cycles  int            // Number of cycles - NOTE is in M-cycles so 1 M-cycle = 4 T-cycles
	ExecFun func(cpu *CPU) // the function for the instruction to execute
}

///////////////////////////
// Instruction Functions //
///////////////////////////
// TODO - actually like properly implement these as they're mostly base stuff atm
// NOTE - at the moment every instance of 'var r byte' is a stand in for a like value gotten from a register which will be added later

//// ADD FUNCTIONS /////

// addBytes add and then return the addition of two bytes
func (cpu *CPU) addBytes(a, b byte) byte {
	result, overflow := overflowAdd[byte](a, b)

	cpu.ResetAllFlags() //reset flags so I can avoid else statements

	if result == 0 {
		cpu.SetFlag(Z, true)
	}
	if halfCarryAdd8b(a, b) {
		cpu.SetFlag(H, true)
	}
	if overflow {
		cpu.SetFlag(C, true)
	}
	
	return result
}

// addByteWithCarry adds two bytes together with the carry
func (cpu *CPU) addByteWithCarry(a, b byte) byte {
	result := a + boolToBit(cpu.Reg.F.carry) + b
	overflow := result < a // TODO check this is actually right?

	cpu.ResetAllFlags()

	cpu.Reg.F.setZeroIf(result)
	cpu.SetFlag(C, overflow)
	// TODO add carry and half carry stuff

	return result
}

// AddRegToA add value stored in register r to A
func (cpu *CPU) AddRegToA() {
	var r byte //imagine this as specific register value
	cpu.Reg.A = cpu.addBytes(cpu.Reg.A, r)
}

//// INC FUNCTIONS /////

// Inc8BitReg increment the value stored in specified register
func (cpu *CPU) Inc8BitReg(r *byte) {
	*r = cpu.incByte(*r)
}

// Inc16BitReg increment the data stored at the specified register pair r1,r2
func (cpu *CPU) Inc16BitReg(r1, r2 *byte)  {
	
}

//// DEC FUNCTIONS /////

// Dec8bReg decrement the data stored in register r
func (cpu *CPU) Dec8bReg(r *byte)  {
	*r = cpu.decByte(*r)
}


//// LD (load) FUNCTIONS /////

// Load8bRegTo8bReg load data from one 8-bit register r1 into another 8-bit register r2
// LD r, r’: Load register (register)
func (cpu *CPU) Load8bRegInto8bReg() {
	var r1, r2 *byte
	*r1 = *r2
}

// Load8bDataTo8bReg load immediate data n into 8-bit register r
// LD r, n: Load register (immediate)
func (cpu *CPU) Load8bDataInto8bReg(r *byte) {
	*r = cpu.ReadByte(cpu.PC)
	// TODO finish these after I figure out how the fuck n8 exists
}

// Load16bRegTo8bReg loads the data from the adress specified by conjoined reg2,reg3 into reg1
// LD r, (HL): Load register (indirect HL)
func (cpu *CPU) Load16bRegInto8bReg(reg1, reg2, reg3 *uint8) {
	*reg1 = ReadByteFrom16bReg(reg2, reg3)
}

// Load8bRegInto16bReg load the data from 8-bit register r1 into the adress specified in 16-bit register r2,r3
func (cpu *CPU) Load8bRegInto16bReg(r1, r2, r3 *byte) {

}

// Load8bDataInto16bRegPair loads 8bit immediate data into register pair r1,r2
func (cpu *CPU) Load8bDataInto16bRegPair(r1, r2 *byte) {

}

// Load16bDataInto16bRegPair load the 16bit data into the register pair r1,r2
// NOTE - I think you do the like top byte into one and the bottom byte into another?
func (cpu *CPU) Load16bDataInto16bRegPair(r1, r2 *byte)  {
	
}

//// MISC FUNCTIONS /////

// Nop No Operation. Doesn't do anything, only increases counter due to nothing being done
func (cpu *CPU) Nop()  {
	// Nop-ing so hard rn
}

// Stop switches the system to STOP mode
func (cpu *CPU) Stop()  {
	// TODO - figure out how to like actually implement this
}



//// INSTRUCTIONS /////

var InstructionsUnprefixed []*Instruction = []*Instruction{
	&Instruction{0x00, "NOP", 1, 1, func(cpu *CPU) {cpu.Nop()}},
	&Instruction{0x01, "LD BC, d16", 3, 3, func(cpu *CPU) {cpu.Load16bDataInto16bRegPair(&cpu.Reg.B, &cpu.Reg.C)}},
	&Instruction{0x02, "LD BC, A", 1, 2, func(cpu *CPU) {cpu.Load8bRegInto16bReg(&cpu.Reg.A, &cpu.Reg.B, &cpu.Reg.C)}},
	&Instruction{0x03, "INC BC", 1, 2, func(cpu *CPU) {cpu.Inc16BitReg(&cpu.Reg.B, &cpu.Reg.C)}},
	&Instruction{0x04, "INC B", 1, 1, func(cpu *CPU) {cpu.Inc8BitReg(&cpu.Reg.B)}},
	&Instruction{0x05, "DEC B", 1, 1, func(cpu *CPU) {cpu.Dec8bReg(&cpu.Reg.B)}},
	&Instruction{0x06, "LD B d8", 2, 2, func(cpu *CPU) {cpu.Load8bDataInto8bReg(&cpu.Reg.B)}}, 
	&Instruction{0x07, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x08, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x09, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x0A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x0B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x0C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x0D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x0E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x0F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x10, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x11, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x12, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x13, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x14, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x15, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x16, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x17, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x18, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x19, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x1A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x1B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x1C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x1D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x1E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x1F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x20, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x21, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x22, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x23, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x24, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x25, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x26, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x27, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x28, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x29, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x2A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x2B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x2C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x2D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x2E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x2F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x30, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x31, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x32, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x33, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x34, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x35, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x36, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x37, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x38, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x39, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x3A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x3B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x3C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x3D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x3E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x3F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x40, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x41, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x42, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x43, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x44, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x45, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x46, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x47, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x48, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x49, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x4A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x4B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x4C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x4D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x4E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x4F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x50, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x51, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x52, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x53, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x54, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x55, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x56, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x57, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x58, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x59, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x5A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x5B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x5C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x5D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x5E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x5F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x60, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x61, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x62, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x63, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x64, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x65, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x66, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x67, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x68, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x69, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x6A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x6B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x6C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x6D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x6E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x6F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x70, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x71, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x72, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x73, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x74, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x75, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x76, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x77, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x78, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x79, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x7A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x7B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x7C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x7D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x7E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x7F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x80, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x81, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x82, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x83, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x84, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x85, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x86, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x87, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x88, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x89, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x8A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x8B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x8C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x8D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x8E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x8F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x90, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x91, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x92, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x93, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x94, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x95, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x96, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x97, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x98, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x99, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x9A, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x9B, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x9C, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x9D, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x9E, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0x9F, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA0, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA1, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA2, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA3, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA4, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA5, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA6, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA7, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA8, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xA9, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xAA, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xAB, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xAC, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xAD, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xAE, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xAF, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB0, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB1, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB2, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB3, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB4, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB5, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB6, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB7, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB8, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xB9, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xBA, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xBB, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xBC, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xBD, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xBE, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xBF, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC0, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC1, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC2, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC3, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC4, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC5, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC6, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC7, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC8, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xC9, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xCA, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xCB, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xCC, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xCD, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xCE, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xCF, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD0, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD1, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD2, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD3, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD4, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD5, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD6, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD7, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD8, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xD9, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xDA, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xDB, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xDC, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xDD, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xDE, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xDF, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE0, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE1, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE2, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE3, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE4, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE5, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE6, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE7, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE8, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xE9, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xEA, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xEB, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xEC, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xED, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xEE, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xEF, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF0, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF1, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF2, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF3, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF4, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF5, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF6, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF7, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF8, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xF9, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xFA, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xFB, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xFC, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xFD, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xFE, "", 0, 0, func(cpu *CPU) {}},
	&Instruction{0xFF, "", 0, 0, func(cpu *CPU) {}},
}
