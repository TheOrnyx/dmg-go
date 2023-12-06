package cpu

// Instruction struct containing information about an instruction
// kinda plagirazed from github.com/djhworld/gomeboycolor
type Instruction struct {
	OpCode      byte           // the opcode of the instruction, in hex form usually
	Desc        string         // basically just the name of the instruction, used for debug purposes
	OperandAmnt int            // size of the instruction aka the amount of operands it has
	Cycles      int            // Number of cycles - NOTE is in M-cycles so 1 M-cycle = 4 T-cycles
	ExecFun     func(cpu *CPU) // the function for the instruction to execute
}

///////////////////////////
// Instruction Functions //
///////////////////////////
// TODO - actually like properly implement these as they're mostly base stuff atm
// NOTE - at the moment every instance of 'var r byte' is a stand in for a like value gotten from a register which will be added later
// TODO - Possibly change the functions to use like 8Bit instead of 8b
// TODO - probably reordering some of my stuff to be based on specific things might be better - also possibly renaming

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

// AddRegToA add value stored in register r to A
func (cpu *CPU) AddRegToA() {
	var r byte //imagine this as specific register value
	cpu.Reg.A = cpu.addBytes(cpu.Reg.A, r)
}

// Add16bRegToHLReg add 16b data located in reg1,reg2 to data located in HL
func (cpu *CPU) Add16bRegToHLReg(reg1, reg2 *uint8) {
	hl := cpu.Reg.HL()
	regPair := JoinBytes(*reg1, *reg2)

	result := cpu.Add16(hl, regPair)

	high, low := Split16(result)
	cpu.Reg.H = high
	cpu.Reg.L = low
}

// AddSPToHLReg Add the stack pointer value to the HL Register
func (cpu *CPU) AddSPToHLReg()  {
	result := cpu.Add16(cpu.Reg.HL(), cpu.SP)
	high, low := Split16(result)
	cpu.Reg.H, cpu.Reg.L = high, low
}

//// INC FUNCTIONS /////

// Inc8BitReg increment the value stored in specified register
func (cpu *CPU) Inc8BitReg(r *byte) {
	*r = cpu.incByte(*r)
}

// Inc16BitRegPair increment the specified register pair r1,r2
func (cpu *CPU) Inc16BitRegPair(r1, r2 *byte) {
	regPair := JoinBytes(*r1, *r2)
	regPair += 1

	*r1, *r2 = Split16(regPair)
}

// Inc16BitRegister increment 16-bit register
// once again, pretty much only for the stack pointer
func (cpu *CPU) Inc16BitRegister(r *uint16)  {
	*r += 1
}

// Inc16BitRegData increment the data located at address located by r1,r2 pair
func (cpu *CPU) Inc16BitRegData(r1, r2 *byte, useFlags bool)  {
	address := JoinBytes(*r1, *r2)
	oldVal := cpu.ReadByte(address)
	result := oldVal+1
	cpu.WriteByteToAddr(address, result)

	if useFlags {
		cpu.SetFlag(Z, result == 0)
		cpu.SetFlag(N, false)
		cpu.SetFlag(H, halfCarryAdd8b(oldVal, 1))
	}
	// TODO - check this is the correct way
}

// IncHLRegData Increment the data pointed to by the HL reg adress.
// Only making one for the HL register as it's the only one to have this
func (cpu *CPU) IncHLRegData()  {
	high, low := cpu.Reg.HLByte()
	cpu.Inc16BitRegData(high, low, true) // NOTE - meh, not really needed can probs just call this directly
}

//// DEC FUNCTIONS /////

// Dec8BitReg decrement the data stored in register r
func (cpu *CPU) Dec8BitReg(r *byte) {
	*r = cpu.decByte(*r)
}

// Dec16BitReg decrement the combination register r1,r2
// NOTE - not the data stored in it
func (cpu *CPU) Dec16BitReg(r1, r2 *uint8) {
	regPair := JoinBytes(*r1, *r2)
	regPair -= 1
	*r1, *r2 = Split16(regPair)
}

// Dec16BitRegData decrement the data located at address r1,r2
func (cpu *CPU) Dec16BitRegData(r1, r2 *byte, useFlags bool)  {
	address := JoinBytes(*r1, *r2)
	oldVal := cpu.ReadByte(address)
	cpu.WriteByteToAddr(address, oldVal-1)

	if useFlags {		
		cpu.SetFlag(N, true)
		cpu.SetFlag(Z, oldVal-1 == 0)
		cpu.SetFlag(H, halfCarrySub8b(oldVal, 1)) // TODO - maybe replace this so I'm not doing the calculation 2 times
	}
	// TODO - check this duh
}

// DecHLRegData Decrement the data pointed to by the adress in HL reg
// Once again only just for HL cuz that's where it's only ever used for
func (cpu *CPU) DecHLRegData()  {
	high, low := cpu.Reg.HLByte()
	cpu.Dec16BitRegData(high, low, true)
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
func (cpu *CPU) Load8BitDataInto8BitReg(r *byte) {
	*r = cpu.CurrentInstruction.Operands[0]
}


// Load8bRegInto16bRegAddr load the data from 8-bit register r into the adress specified in 16-bit register r1,r2
func (cpu *CPU) Load8bRegInto16bRegAddr(r1, r2, r *byte) {
	regPair := JoinBytes(*r1, *r2)
	cpu.WriteByteToAddr(regPair, *r)
}

// Load8BitDataInto16BitRegAddr load the immediate data into adress specified by r1,r2
func (cpu *CPU) Load8BitDataInto16BitRegAddr(r1, r2 *byte)  {
	data := cpu.CurrentInstruction.Operands[0]
	addr := JoinBytes(*r1, *r2)
	cpu.WriteByteToAddr(addr, data)
	// TODO - check this is right
}

// Load8bRegInto16bRegAddrInc load the data from 8-bit register r
// into the adress specified in 16-bit register r1,r2 and then
// increments the regPair after
func (cpu *CPU) Load8bRegInto16bRegAddrInc(r1, r2, r *byte) {
	regPair := JoinBytes(*r1, *r2)
	cpu.WriteByteToAddr(regPair, *r)
	cpu.Inc16BitRegData(r1, r2, false) // TODO - check this is right
}

// Load8bRegInto16bRegAddrDec load the data from 8-bit register r
// into the adress specified in 16-bit register r1,r2 and then
// Decrement the regPair data after
func (cpu *CPU) Load8bRegInto16bRegAddrDec(r1, r2, r *byte) {
	regPair := JoinBytes(*r1, *r2)
	cpu.WriteByteToAddr(regPair, *r)
	cpu.Dec16BitRegData(r1, r2, false) // TODO - check this is right
}

// Load16BitAddrIncInto8BitReg Load the data stored in the adress pointed to by r1,r2 into reg r
// Increment the data stored in address at r1,r2 after
func (cpu *CPU) Load16BitAddrIncInto8BitReg(r1, r2, r *byte)  {
	*r = cpu.ReadByte(JoinBytes(*r1, *r2))
	cpu.Inc16BitRegData(r1, r2, false)
	// TODO - check
}

// Load16BitAddrDecInto8BitReg Load the data stored in the adress pointed to by r1,r2 into reg r
// Decrement the data stored in address at r1,r2 after
func (cpu *CPU) Load16BitAddrDecInto8BitReg(r1, r2, r *byte)  {
	*r = cpu.ReadByte(JoinBytes(*r1, *r2))
	cpu.Dec16BitRegData(r1, r2, false)
	// TODO - check
}

// Load16BitDataInto16BitRegPair load the 16bit data into the register pair r1,r2
// NOTE - I think you do the like top byte into one and the bottom byte into another?
func (cpu *CPU) Load16BitDataInto16BitRegPair(r1, r2 *byte) {
	low, high := cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	*r1 = high
	*r2 = low
}

// Load16BitDataInto16BitRegister load the 16bit immediate data into the 16-bit register r
// Basically only used for like the stack pointer
func (cpu *CPU) Load16BitDataInto16BitRegister(r *uint16)  {
	low, high := cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	*r = JoinBytes(high, low)
	// TODO - check this is right
}

// Load16bRegDataIntoReg load the data stored in register combination regHigh,regLow into reg
func (cpu *CPU) Load16bRegDataIntoReg(regHigh, regLow, reg *uint8) {
	*reg = cpu.ReadByte(JoinBytes(*regHigh, *regLow))
}

// LoadStackPointerInto16bData load the data from the stack pointer into the adress specified nn
func (cpu *CPU) LoadStackPointerInto16bData() {
	var low, high byte = cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	address := JoinBytes(high, low)

	cpu.WriteByteToAddr(address, byte(cpu.SP&0xFF))
	cpu.WriteByteToAddr((address+1)&0xFFFF, byte(cpu.SP>>8))
	// TODO check if need to increase PC here
}

//// ROTATE FUNCTIONS /////

// RotateLeftCarryRegA Rotate Register A left by one
// Store the value of bit 7(?) in the carry flag
// 1 M-cycle, 1 byte length
func (cpu *CPU) RotateLeftCarryRegA() {
	var bit7 bool = false

	cpu.ResetAllFlags() // reset cuz only the carry flag is set

	if cpu.Reg.A&0x80 == 0x80 {
		bit7 = true
	}

	outcome := ((cpu.Reg.A << 1) & 0xFF) | (boolToBit(bit7) & 0x01)

	if bit7 {
		cpu.SetFlag(C, true)
	}

	cpu.Reg.A = outcome
}

// RotateLeftRegA rotate Register A left by one
func (cpu *CPU) RotateLeftRegA() {
	var bit7 bool = false

	if cpu.Reg.A&0x80 == 0x80 {
		bit7 = true
	}

	outcome := (cpu.Reg.A << 1) | (boolToBit(cpu.Reg.F.carry) & 0x01)

	cpu.ResetAllFlags()
	cpu.SetFlag(C, bit7)
	cpu.Reg.A = outcome
	// TODO - check this is actually right, have no idea if you are supposed to set bit 0 to previous carry flag
}

// RotateRightCarryRegA rorate Register A right by one
// Store value of bit 0 in the carry flag
func (cpu *CPU) RotateRightCarryRegA() {
	var bit0 bool = false

	cpu.ResetAllFlags()

	if cpu.Reg.A&0x01 == 0x01 {
		bit0 = true
	}

	result := (cpu.Reg.A >> 1) | (boolToBit(bit0) << 7 & 0x80)

	cpu.SetFlag(C, bit0)

	cpu.Reg.A = result
}

// RotateRightRegA rotate register A right by one
func (cpu *CPU) RotateRightRegA()  {
	var bit0 bool = false

	if cpu.Reg.A&0x01 == 0x01 {
		bit0 = true
	}

	outcome := (cpu.Reg.A >> 1) | (boolToBit(cpu.Reg.F.carry) << 7 & 0x80)

	cpu.ResetAllFlags()
	cpu.SetFlag(C, bit0)
	cpu.Reg.A = outcome
}

//// JUMP FUNCTIONS /////

// JumpRelative8bit make an unconditional jump relative by a signed 8-bit operand
// TODO - check that the overflow on the signed value works properly
func (cpu *CPU) JumpRelative8bit() {
	jumpVal := int8(cpu.CurrentInstruction.Operands[0]) // convert to signed value, should hopefully overflow to correct negative

	if jumpVal < 0 { // stupid negative values
		cpu.PC -= uint16(jumpVal)
	} else {
		cpu.PC += uint16(jumpVal)
	}

	cpu.hasJumped = true
}

// JumpConditionalRelative8bit jump relative amount when cond is true
// Can be used for both not and regular jumps
func (cpu *CPU) JumpConditionalRelative8bit(cond *bool, jumpWhen bool)  {
	jumpVal := int8(cpu.CurrentInstruction.Operands[0])

	if *cond == jumpWhen {
		cpu.CurrentInstruction.Instruction.Cycles = 3 // instruction has variable cylce amount
		if jumpVal < 0 { // stupid negative values
			cpu.PC -= uint16(jumpVal)
		} else {
			cpu.PC += uint16(jumpVal)
		}
		cpu.hasJumped = true
	} else {
		cpu.CurrentInstruction.Instruction.Cycles = 2 // just for saffety
	}
}

//// MISC FUNCTIONS /////

// Nop No Operation. Doesn't do anything, only increases counter due to nothing being done
func (cpu *CPU) Nop() {
	// Nop-ing so hard rn
}

// Stop switches the system to STOP mode
func (cpu *CPU) Stop() {
	// TODO - figure out how to like actually implement this
}

// ComplementRegA take the complement of Reg A (flip the bits)
func (cpu *CPU) ComplementRegA()  {
	cpu.Reg.A = ^cpu.Reg.A
	cpu.SetFlag(N, true)
	cpu.SetFlag(H, true)
}

// SetCarryFlag set the carry flag to true
func (cpu *CPU) SetCarryFlag()  {
	cpu.SetFlag(C, true)
}

//// INSTRUCTIONS /////

// InstructionsUnprefixed - the slice to represent each of the
// Unprefixed CPU instructions Accessed using just the index in hex
// format as that's easiest to access, then compiling the operands and
// running the associated function
var InstructionsUnprefixed []*Instruction = []*Instruction{
	&Instruction{0x00, "NOP", 0, 1, func(cpu *CPU) { cpu.Nop() }},
	&Instruction{0x01, "LD BC, d16", 2, 3, func(cpu *CPU) { cpu.Load16BitDataInto16BitRegPair(&cpu.Reg.B, &cpu.Reg.C) }},
	&Instruction{0x02, "LD (BC), A", 0, 2, func(cpu *CPU) { cpu.Load8bRegInto16bRegAddr(&cpu.Reg.B, &cpu.Reg.C, &cpu.Reg.A) }},
	&Instruction{0x03, "INC BC", 0, 2, func(cpu *CPU) { cpu.Inc16BitRegPair(&cpu.Reg.B, &cpu.Reg.C) }},
	&Instruction{0x04, "INC B", 0, 1, func(cpu *CPU) { cpu.Inc8BitReg(&cpu.Reg.B) }},
	&Instruction{0x05, "DEC B", 0, 1, func(cpu *CPU) { cpu.Dec8BitReg(&cpu.Reg.B) }},
	&Instruction{0x06, "LD B d8", 2, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.B) }},
	&Instruction{0x07, "RLCA", 0, 1, func(cpu *CPU) { cpu.RotateLeftCarryRegA() }},
	&Instruction{0x08, "LD (a16), SP", 2, 5, func(cpu *CPU) { cpu.LoadStackPointerInto16bData() }},
	&Instruction{0x09, "ADD HL, BC", 0, 2, func(cpu *CPU) { cpu.Add16bRegToHLReg(&cpu.Reg.B, &cpu.Reg.C) }},
	&Instruction{0x0A, "LD A, (BC)", 0, 2, func(cpu *CPU) { cpu.Load16bRegDataIntoReg(&cpu.Reg.B, &cpu.Reg.C, &cpu.Reg.A) }},
	&Instruction{0x0B, "DEC BC", 0, 2, func(cpu *CPU) { cpu.Dec16BitReg(&cpu.Reg.B, &cpu.Reg.C) }},
	&Instruction{0x0C, "INC C", 0, 1, func(cpu *CPU) { cpu.Inc8BitReg(&cpu.Reg.C) }},
	&Instruction{0x0D, "DEC C", 0, 1, func(cpu *CPU) { cpu.Dec8BitReg(&cpu.Reg.C) }},
	&Instruction{0x0E, "LD C, d8", 1, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.C) }},
	&Instruction{0x0F, "RRCA", 0, 1, func(cpu *CPU) { cpu.RotateRightCarryRegA() }},
	&Instruction{0x10, "STOP", 1, 1, func(cpu *CPU) { cpu.Stop() }},
	&Instruction{0x11, "LD DE, d16", 2, 3, func(cpu *CPU) { cpu.Load16BitDataInto16BitRegPair(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0x12, "LD (DE), A", 0, 2, func(cpu *CPU) { cpu.Load8bRegInto16bRegAddr(&cpu.Reg.D, &cpu.Reg.E, &cpu.Reg.A) }},
	&Instruction{0x13, "INC DE", 0, 2, func(cpu *CPU) { cpu.Inc16BitRegPair(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0x14, "INC D", 0, 1, func(cpu *CPU) { cpu.Inc8BitReg(&cpu.Reg.D) }},
	&Instruction{0x15, "DEC D", 0, 1, func(cpu *CPU) { cpu.Dec8BitReg(&cpu.Reg.D) }},
	&Instruction{0x16, "LD D, d8", 1, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.D) }},
	&Instruction{0x17, "RLA", 0, 1, func(cpu *CPU) { cpu.RotateLeftRegA() }},
	&Instruction{0x18, "JR s8", 1, 3, func(cpu *CPU) { cpu.JumpRelative8bit() }},
	&Instruction{0x19, "ADD HL, DE", 0, 2, func(cpu *CPU) { cpu.Add16bRegToHLReg(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0x1A, "LD A, (DE)", 0, 2, func(cpu *CPU) {cpu.Load16bRegDataIntoReg(&cpu.Reg.D, &cpu.Reg.E, &cpu.Reg.A)}},
	&Instruction{0x1B, "DEC DE", 0, 2, func(cpu *CPU) {cpu.Dec16BitReg(&cpu.Reg.D, &cpu.Reg.E)}},
	&Instruction{0x1C, "INC E", 0, 1, func(cpu *CPU) {cpu.Inc8BitReg(&cpu.Reg.E)}},
	&Instruction{0x1D, "DEC E", 0, 1, func(cpu *CPU) {cpu.Dec8BitReg(&cpu.Reg.E)}},
	&Instruction{0x1E, "LD E, d8", 1, 2, func(cpu *CPU) {cpu.Load8BitDataInto8BitReg(&cpu.Reg.E)}},
	&Instruction{0x1F, "RRA", 0, 1, func(cpu *CPU) {cpu.RotateRightRegA()}},
	&Instruction{0x20, "JR NZ, s8", 1, 2, func(cpu *CPU) {cpu.JumpConditionalRelative8bit(&cpu.Reg.F.zero, false)}},
	&Instruction{0x21, "LD HL, d16", 2, 3, func(cpu *CPU) {cpu.Load16BitDataInto16BitRegPair(&cpu.Reg.H, &cpu.Reg.L)}},
	&Instruction{0x22, "LD (HL+), A", 0, 2, func(cpu *CPU) {cpu.Load8bRegInto16bRegAddrInc(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A)}},
	&Instruction{0x23, "INC HL", 0, 2, func(cpu *CPU) {cpu.Inc16BitRegPair(&cpu.Reg.H, &cpu.Reg.L)}},
	&Instruction{0x24, "INC H", 0, 1, func(cpu *CPU) {cpu.Inc8BitReg(&cpu.Reg.H)}},
	&Instruction{0x25, "DEC H", 0, 1, func(cpu *CPU) {cpu.Dec8BitReg(&cpu.Reg.H)}},
	&Instruction{0x26, "LD H, d8", 1, 2, func(cpu *CPU) {cpu.Load8BitDataInto8BitReg(&cpu.Reg.H)}},
	&Instruction{0x27, "DAA", 0, 1, func(cpu *CPU) {}}, // TODO - implement this
	&Instruction{0x28, "JR Z, s8", 1, 2, func(cpu *CPU) {cpu.JumpConditionalRelative8bit(&cpu.Reg.F.zero, true)}},
	&Instruction{0x29, "ADD HL, HL", 0, 2, func(cpu *CPU) {cpu.Add16bRegToHLReg(&cpu.Reg.H, &cpu.Reg.L)}}, // lmao
	&Instruction{0x2A, "LD A, (HL+)", 0, 2, func(cpu *CPU) {cpu.Load16BitAddrIncInto8BitReg(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A)}},
	&Instruction{0x2B, "DEC HL", 0, 2, func(cpu *CPU) {cpu.Dec16BitReg(&cpu.Reg.H, &cpu.Reg.L)}},
	&Instruction{0x2C, "INC L", 0, 1, func(cpu *CPU) {cpu.Inc8BitReg(&cpu.Reg.L)}},
	&Instruction{0x2D, "DEC L", 0, 1, func(cpu *CPU) {cpu.Dec8BitReg(&cpu.Reg.L)}},
	&Instruction{0x2E, "LD L, d8", 1, 2, func(cpu *CPU) {cpu.Load8BitDataInto8BitReg(&cpu.Reg.L)}},
	&Instruction{0x2F, "CPL", 0, 1, func(cpu *CPU) {cpu.ComplementRegA()}},
	&Instruction{0x30, "JR NC, s8", 2, 2, func(cpu *CPU) {cpu.JumpConditionalRelative8bit(&cpu.Reg.F.carry, false)}},
	&Instruction{0x31, "LD SP, d16", 2, 3, func(cpu *CPU) {cpu.Load16BitDataInto16BitRegister(&cpu.SP)}},
	&Instruction{0x32, "lD (HL-), A", 0, 2, func(cpu *CPU) {cpu.Load8bRegInto16bRegAddrDec(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A)}},
	&Instruction{0x33, "INC SP", 0, 2, func(cpu *CPU) {cpu.Inc16BitRegister(&cpu.SP)}},
	&Instruction{0x34, "INC (HL)", 0, 3, func(cpu *CPU) {cpu.IncHLRegData()}},
	&Instruction{0x35, "DEC (HL)", 0, 3, func(cpu *CPU) {cpu.DecHLRegData()}},
	&Instruction{0x36, "LD (HL), d8", 1, 3, func(cpu *CPU) {cpu.Load8BitDataInto16BitRegAddr(cpu.Reg.HLByte())}},
	&Instruction{0x37, "SCF", 0, 1, func(cpu *CPU) {cpu.SetCarryFlag()}},
	&Instruction{0x38, "JR C, s8", 1, 2, func(cpu *CPU) {cpu.JumpConditionalRelative8bit(&cpu.Reg.F.carry, true)}},
	&Instruction{0x39, "ADD HL, SP", 0, 2, func(cpu *CPU) {cpu.AddSPToHLReg()}},
	&Instruction{0x3A, "LD A, (HL-)", 0, 2, func(cpu *CPU) {}},
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
