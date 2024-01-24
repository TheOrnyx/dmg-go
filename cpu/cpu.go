package cpu

import (
	"fmt"
	"log"
	"os"

	"github.com/TheOrnyx/gameboy-golor/mmu"
	"github.com/TheOrnyx/gameboy-golor/timer"
)

/////////////
// TODO's  //
/////////////
//
// TODO - add like ticks and stuff
// TODO - probably kill the operands system and replace it with just loading the values

var infoLog = log.New(os.Stdout, "[INFO] ", log.Ldate)
var debugLog = log.New(os.Stdout, "[DEBUG] ", log.Ldate)
var warnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
var fatalLog = log.New(os.Stdout, "[FaTAL] ", log.LstdFlags)

const ( // the interrupt addresses
	granularIEAddr = 0xFFFF
	interruptFlagAddr = 0xFF0F
)

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
	MMU *mmu.MMU // the memory mapped unit 
	hasJumped bool // bool for whether the CPU has just ran a jump instruction, TODO - implement
	InterruptsEnabled bool // bool for whether or not the interrupt flag has been enabled
	Halted bool // whether or not the CPU is halted
	Timer *timer.Timer // the cpu timer

	instrCycles int // the current amount of M-cycles for the instruction
}

// String print cpu
func (cpu *CPU) String() string {
	pc := cpu.MMU.ReadByte(cpu.PC)
	pcOne := cpu.MMU.ReadByte(cpu.PC+1)
	pcTwo := cpu.MMU.ReadByte(cpu.PC+2)
	pcThree := cpu.MMU.ReadByte(cpu.PC+3)
	instructionInfo := fmt.Sprintf("Current Instruction: OpCode '0x%X': Operands '%v'/'%v': Cycles '%v': Name '%s'", cpu.CurrentInstruction.Instruction.OpCode, cpu.CurrentInstruction.Instruction.OperandAmnt, cpu.CurrentInstruction.Operands, cpu.CurrentInstruction.Instruction.Cycles, cpu.CurrentInstruction.Instruction.Desc)
	extraInfo := fmt.Sprintf("F Register Bools: %v\nItems at addresses: HL Addr:0x%04X HL Data:%v", cpu.Reg.F.String(), cpu.Reg.HL(), cpu.MMU.ReadByte(cpu.Reg.HL()))
	return fmt.Sprintf("%s\nCPU Values: %v SP:0x%04X PC:0x%04X PCMEM:%v,%v,%v,%v\n%v", instructionInfo, cpu.Reg.String(), cpu.SP, cpu.PC, pc, pcOne, pcTwo, pcThree, extraInfo)
}

// StringDoctor string for cpu in gameboy doctor form
func (cpu *CPU) StringDoctor() string {
	pc := cpu.MMU.ReadByte(cpu.PC)
	pcOne := cpu.MMU.ReadByte(cpu.PC+1)
	pcTwo := cpu.MMU.ReadByte(cpu.PC+2)
	pcThree := cpu.MMU.ReadByte(cpu.PC+3)
	return fmt.Sprintf("%v SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X", cpu.Reg.StringDoctor(), cpu.SP, cpu.PC, pc, pcOne, pcTwo, pcThree)
}

// NewCPU create and return a new cpu
func NewCPU(mmu *mmu.MMU, timer *timer.Timer) (*CPU, error) {
	newCPU := new(CPU)
	newCPU.MMU = mmu
	newCPU.Timer = timer
	newCPU.Reset()

	return newCPU, nil
}

// Reset reset the cpu
func (cpu *CPU) Reset()  {
	// infoLog.Println("Resetting CPU...")
	cpu.PC = 0
	cpu.SP = 0
	cpu.Reg.reset()
	cpu.ResetAllFlags()
	cpu.Halted = false
	cpu.InterruptsEnabled = false
	cpu.hasJumped = false
	cpu.CurrentInstruction = &CurrentInstruction{Instruction: InstructionsUnprefixed[0x00], Operands: [2]byte{}}
}

// ResetDebug reset the cpu to debug position (basically skip the boot rom)
func (cpu *CPU) ResetDebug()  {
	cpu.Reg.A = 0x01

	cpu.SetFlag(C, true) // TODO - these depend on the checksum in the header
	cpu.SetFlag(H, true) // TODO - these depend on the checksum in the header
	cpu.SetFlag(N, false)
	cpu.SetFlag(Z, true)

	cpu.Reg.B = 0x00
	cpu.Reg.C = 0x13
	cpu.Reg.D = 0x00
	cpu.Reg.E = 0xD8
	cpu.Reg.H = 0x01
	cpu.Reg.L = 0x4D

	cpu.SP = 0xFFFE
	cpu.PC = 0x0100
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

// Tick tick the cpu by cycles amount of M-cycles
func (cpu *CPU) Tick(cycles int)  {
	cpu.Timer.Tick(cycles)
}

// RequestInterrupt request an interrupt with code
func (cpu *CPU) RequestInterrupt(code byte)  {
	
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
	cpu.instrCycles = 1
	if !cpu.Halted {
		// cpu.checkInterrupts()
		opCode := cpu.readPC()
		
		if opCode == 0xCB { // use the prefixed instructions
			newOpCode := cpu.readPC()
			cpu.CompileInstruction(InstructionsPrefixed[newOpCode])
		} else {
			cpu.CompileInstruction(InstructionsUnprefixed[opCode])
		}

		cpu.CurrentInstruction.Instruction.ExecFun(cpu)

		cpu.hasJumped = false
		
	} else {
		// interruptFlagNow := cpu.ReadByte(0xFF0F)
		// cpu.Halted = false
		// cpu.Tick(1)
	}

	// cpu.Tick(cpu.instrCycles)
	return 0
}

// checkInterrupts check for interrupts and return true if should interrupt
func (cpu *CPU) checkInterrupts() bool {
	if !cpu.InterruptsEnabled {
		return false
	}

	ie := cpu.MMU.ReadByte(granularIEAddr)
	iFlag := cpu.MMU.ReadByte(interruptFlagAddr)
	interrupt := ie & iFlag // and the two together to find interrupts that are both enabled and pending
	
	switch {
	case interrupt & 0x01 == 0x01: // VBlank interrupt
		cpu.WriteByteToAddr(interruptFlagAddr, iFlag & 0xFE) // turn off the interrupt request bit and write
		cpu.InterruptsEnabled = false
		cpu.pushSP(cpu.PC)
		cpu.PC = 0x0040
		return true
	case interrupt & 0x02 == 0x02: // LCD interrupt
		cpu.WriteByteToAddr(interruptFlagAddr, iFlag & 0xFD) // turn off interrupt req and write
		cpu.InterruptsEnabled = false
		cpu.pushSP(cpu.PC)
		cpu.PC = 0x0048
		return true
	case interrupt & 0x03 == 0x03: // Timer overflow interrupt
		cpu.WriteByteToAddr(interruptFlagAddr, iFlag & 0xFB)
		cpu.InterruptsEnabled = false
		cpu.pushSP(cpu.PC)
		cpu.PC = 0x0050
		return true
	case interrupt & 0x04 == 0x04: // serial link interrupt
		cpu.WriteByteToAddr(interruptFlagAddr, iFlag & 0xF7) // NOTE - check that 0xF7 is the right hex value to and
		cpu.InterruptsEnabled = false
		cpu.pushSP(cpu.PC)
		cpu.PC = 0x0058
		return true
	case interrupt & 0x05 == 0x05: // joypad interrupt
		cpu.WriteByteToAddr(interruptFlagAddr, iFlag & 0xEF)
		cpu.InterruptsEnabled = false
		cpu.pushSP(cpu.PC)
		cpu.PC = 0x0060
		return true
	default:
		// infoLog.Printf("Unkown interrupt with interrupt %v\n", interrupt)
	}
	
	return false
}

// IncrementPC increment the PC by amnt
func (cpu *CPU) IncrementPC(amnt int)  {
	cpu.PC += uint16(amnt)
}

// readPC read the value at the current PC and increment it
func (cpu *CPU) readPC() byte {
	data := cpu.ReadByte(cpu.PC)
	cpu.IncrementPC(1)
	return data
}

// Add16 add 2 16bit numbers together, set flags and return the result
func (cpu *CPU) Add16(a, b uint16) uint16 {
	result := a + b

	cpu.ResetFlag(N)
	cpu.ResetFlag(C)
	cpu.ResetFlag(H)
	
	if result < a { // TODO - convert
		cpu.SetFlag(C, true)
	}

	if (a & 0xFFF) + (b & 0xFFF) > 0xFFF {
		cpu.SetFlag(H, true)
	}

	return result
}


// CompileInstruction compile the given instruction, properly reading
// and attaching the correct number of operands
// pretty much stolen tbh
func (cpu *CPU) CompileInstruction(instruction *Instruction)  {
	cpu.CurrentInstruction.Instruction = instruction
	cpu.CurrentInstruction.Operands[0] = 0
	cpu.CurrentInstruction.Operands[1] = 0 // TODO - check this doesn't break anything
	switch instruction.OperandAmnt {
	case 1:
		cpu.CurrentInstruction.Operands[0] = cpu.readPC()
	case 2:
		cpu.CurrentInstruction.Operands[0] = cpu.readPC()
		cpu.CurrentInstruction.Operands[1] = cpu.readPC()
	}
}

// ReadByte reads the byte at address addr and returns it
func (cpu *CPU) ReadByte(addr uint16) byte {
	return cpu.MMU.ReadByte(addr)
}

// WriteByteToAddr write data to address located at addr
func (cpu *CPU) WriteByteToAddr(addr uint16, data byte)  {
	cpu.MMU.WriteByte(addr, data)
}

/////////////////////////////
// Stack Pointer Functions //
/////////////////////////////

// pushSP push data onto the stack pointer
// TODO - Check I have the byte order correct
func (cpu *CPU) pushSP(data uint16)  {
	msb, lsb := Split16(data) // split to the MSB and LSB
	
	cpu.SP -= 1
	cpu.WriteByteToAddr(cpu.SP, msb)

	cpu.SP -= 1
	cpu.WriteByteToAddr(cpu.SP, lsb)

}

// popSP pop the top off the stack pointer and return the address
// TODO - check byte order is correct
func (cpu *CPU) popSP() uint16 {
	lsb := cpu.ReadByte(cpu.SP)
	cpu.SP += 1

	msb := cpu.ReadByte(cpu.SP)
	cpu.SP += 1

	return JoinBytes(msb, lsb)
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

