package cpu

import (
	"log"
	"os"

	"github.com/TheOrnyx/gameboy-golor/mmu"
)

var infoLog = log.New(os.Stdout, "[INFO] ", log.Ldate)
var debugLog = log.New(os.Stdout, "[DEBUG] ", log.Ldate)
var warnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
var fatalLog = log.New(os.Stdout, "[FaTAL] ", log.LstdFlags)

// CurrentInstruction struct to hold information about the current instruction
type CurrentInstruction struct {
	Operands [2]byte
	Instruction *Instruction
}

// CPU struct to contain information about the CPU
type CPU struct {
	Reg Registers
	PC uint16 // program counter
	SP uint16 // Stack pointer
	CurrentInstruction *CurrentInstruction //the current instruction to be run
	MMU mmu.MMU
	hasJumped bool // bool for whether the CPU has just ran a jump instruction, TODO - implement
}

// Reset reset the cpu
func (cpu *CPU) Reset()  {
	infoLog.Println("Resetting CPU...")
	cpu.PC = 0
	cpu.SP = 0
	cpu.Reg.reset()
	cpu.ResetAllFlags()
	cpu.CurrentInstruction = &CurrentInstruction{Instruction: InstructionsUnprefixed[0x00], Operands: [2]byte{}}
}


// ResetFlag reset given flag to it's default
func (cpu *CPU) ResetFlag(flag int)  {
	switch flag {
	case Z:
		cpu.Reg.F.zero = false
	case N:
		cpu.Reg.F.subtract = false
	case H:
		cpu.Reg.F.half_carry = false
	case C:
		cpu.Reg.F.carry = false
	}
}

// ResetAllFlags resets all flags back to false
func (cpu *CPU) ResetAllFlags()  {
	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.ResetFlag(C)
}

// SetFlag set given flag to specified state
func (cpu *CPU) SetFlag(flag int, state bool)  {
	switch flag {
	case Z:
		cpu.Reg.F.zero = state
	case N:
		cpu.Reg.F.subtract = state
	case H:
		cpu.Reg.F.half_carry = state
	case C:
		cpu.Reg.F.carry = state
	}
}


// Step the step function for the cpu
func (cpu *CPU) Step() int {
	opCode := cpu.ReadByte(cpu.PC)
	cpu.PC++
	
	if opCode == 0xCB { // use the prefixed instructions 
		
	} else {
		cpu.CompileInstruction(InstructionsUnprefixed[opCode])
	}

	cpu.CurrentInstruction.Instruction.ExecFun(cpu)

	if true { // replace with some sort of jump check later
		cpu.IncrementPC(len(cpu.CurrentInstruction.Operands)+1)
	}
	
	return 0
}

// IncrementPC increment the PC by amnt
func (cpu *CPU) IncrementPC(amnt int)  {
	cpu.PC += uint16(amnt)
}

// Add16 add 2 16bit numbers together, set flags and return the result
func (cpu *CPU) Add16(a, b uint16) uint16 {
	result := a + b

	cpu.ResetFlag(N)
	cpu.ResetFlag(C)
	cpu.ResetFlag(H)
	
	if result < a {
		cpu.SetFlag(C, true)
	}

	if result & (1 << 11) != 0 {
		cpu.SetFlag(H, true)
	}

	return result
}


// CompileInstruction compile the given instruction, properly reading
// and attaching the correct number of operands
// pretty much stolen tbh
func (cpu *CPU) CompileInstruction(instruction *Instruction)  {
	cpu.CurrentInstruction.Instruction = instruction
	switch instruction.OperandAmnt {
	case 1:
		cpu.CurrentInstruction.Operands[0] = cpu.ReadByte(cpu.PC + 1)
	case 2:
		cpu.CurrentInstruction.Operands[0] = cpu.ReadByte(cpu.PC + 1)
		cpu.CurrentInstruction.Operands[1] = cpu.ReadByte(cpu.PC + 2)
	}
}

// ReadByte reads the byte at address addr and returns it
func (cpu *CPU) ReadByte(addr uint16) byte {
	// NOTE - this should read a byte from the MMU but I'm too lazy to implement that atm
	return 0
}

// byteToInstruction convert a byte value to an instruction
// if converted value does not map to an instruction then return an error
func byteToInstruction(b uint8) (int, error) {
	
	return 0, nil
}

// incByte increment given byte and set appropriate flags based on the outcome
// also return the outcome
// used in instruction functions
func (cpu *CPU) incByte(val byte) byte {
	origByte := val
	// result, overflow := overflowAdd[byte](val, 1) // NOTE - overflow not handled?
	result := val + 1

	// reset the flags - can also just use else statements but eh
	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	if result == 0 {
		cpu.SetFlag(Z, true)
	}
	// TODO - check if the half carry thing can be changed
	if halfCarryAdd8b(origByte, 1) { // maybe don't do calculation twice idk (optimization)
		cpu.SetFlag(H, true)
	}

	return result
}

// decByte decrement the given byte by 1 and set the appropriate flags
func (cpu *CPU) decByte(val byte) byte {
	result := val - 1

	cpu.ResetFlag(Z)
	cpu.SetFlag(N, true)
	cpu.ResetFlag(H)

	if result == 0 {
		cpu.SetFlag(Z, true)
	}
	if halfCarrySub8b(val, 1) {
		cpu.SetFlag(H, true)
	}
	
	
	return result
}

// WriteByteToAddr write data to address located at addr
func (cpu *CPU) WriteByteToAddr(addr uint16, data byte)  {
	// TODO add bit for writing to the MMU here
}

