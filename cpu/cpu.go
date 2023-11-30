package cpu


// Instructions
const (
	addInstruction = iota
)
const lastInstruction = addInstruction //used to check that the converting byte to instruction works right (make sure to change)

// Arithmetic Target "enum"
const (
	a = iota
	b
	c
	d
	e
	h
	l
)

// CPU struct to contain information about the CPU
type CPU struct {
	Registers Registers
	PC uint16 // program counter
	Bus MemoryBus
}

// // Execute execute an instruction based on the instruction and target
// // target refers to the adress target
// func (cpu *CPU) Execute(instruction, target int)  {
// 	switch instruction {
// 	case addInstruction: // add
// 		switch target {
// 		case c:
// 			cpu.Registers.a = cpu.add(cpu.Registers.c)
// 			// implement more options
// 		}
// 		// implement more instructions, possibly convert to make it easier
// 	}
// }



// Step the step function for the cpu
func (cpu *CPU) Step() int {
	
	return 0
}

// ReadByte return the next byte in the PC (program counter)
// TODO figure out how to convert the program counter to a byte easily
func (cpu *CPU) ReadByte() byte {
	return byte(cpu.PC)
}

// add adds and returns the result of adding value onto register a
func (c *CPU) add(value uint8) uint8 {
	newVal, overflowed := overflowAdd[uint8](c.Registers.A, value)
	c.Registers.F.zero = newVal == 0
	c.Registers.F.subtract = false
	c.Registers.F.carry = overflowed
	// Half Carry is set if adding the lower nibbles of the value and
	// register A together results in a value larger than 0xF
	c.Registers.F.half_carry = (c.Registers.A & 0xF) + (newVal & 0xF) > 0xF
	return newVal
}

type MemoryBus struct {
	Memory [0xFFFF]byte
}

// readByte read and return byte from the memory bus at specific adress
func (m *MemoryBus) readByte(address uint16) uint8 {
	return m.Memory[address]
}

// byteToInstruction convert a byte value to an instruction
// if converted value does not map to an instruction then return an error
func byteToInstruction(b uint8) (int, error) {

	return addInstruction, nil
}

