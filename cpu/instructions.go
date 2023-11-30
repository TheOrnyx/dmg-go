package cpu

// Instruction struct containing information about an instruction
// kinda plagirazed from github.com/djhworld/gomeboycolor
type Instruction struct {
	OpCode byte
	Desc string
	OperandsSize int // size of the instructions operands?
	Cycles int // cycles to complete
	ExecFun func(cpu *CPU)
}



///////////////////////////
// Instruction Functions //
///////////////////////////
// TODO - actually like properly implement these as they're mostly base stuff atm
// NOTE - at the moment every instance of 'var r byte' is a stand in for a like value gotten from a register which will be added later

//// ADD FUNCTIONS /////

// addBytes add and then return the addition of two bytes
func (cpu *CPU) addBytes(a, b byte) byte {
	return a + b
	// TODO - add flags and shit here
}

// AddRegToA add value stored in register r to A
func (cpu *CPU) AddRegToA()  {
	var r byte //imagine this as specific register value
	cpu.Registers.A = cpu.addBytes(cpu.Registers.A, r)
}


//// INC FUNCTIONS /////

// Inc8BitReg increment the value stored in specified register
func (cpu *CPU) Inc8BitReg()  {
	var r *byte
	*r = *r + 1
}


//// LD (load) FUNCTIONS /////

// Load8bRegTo8bReg load data from one 8-bit register r1 into another 8-bit register r2
func (cpu *CPU) Load8bRegTo8bReg()  {
	var r1, r2 *byte
	*r1 = *r2
}

// Load8bDataTo8bReg load immediate data n into 8-bit register r
func (cpu *CPU) Load8bDataTo8bReg()  {
	var r *byte
	// TODO finish these after I figure out how the fuck n8 exists
}


//// INSTRUCTIONS /////


