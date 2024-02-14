package cpu

import "fmt"

// GetInstrDebug get information about the current instruction and
// return it in string form without changing cpu state
func (cpu *CPU) GetInstrDebug() string {
	cpu.MMU.DebugMode = false // disable reporting
	pc := cpu.PC
	var currentInstr *Instruction
	var operands [2]byte
	newOpcode := cpu.MMU.ReadByte(pc)
	
	if newOpcode == 0xCB {
		pc++
		newOpcode = cpu.MMU.ReadByte(pc)
		currentInstr = InstructionsPrefixed[newOpcode]
	} else {
		currentInstr = InstructionsUnprefixed[newOpcode]
	}

	switch currentInstr.OperandAmnt {
	case 1:
		operands[0] = cpu.MMU.ReadByte(pc+1)
	case 2:
		operands[0] = cpu.MMU.ReadByte(pc+1)
		operands[1] = cpu.MMU.ReadByte(pc+2)
	}

	cpu.MMU.DebugMode = true
	return fmt.Sprintf("OpCode:0x%02X  Name:%s  Operands:%v", currentInstr.OpCode, currentInstr.Desc, operands)
}
