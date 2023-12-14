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
// TODO - make sure that all the functions properly assign the results to things in thigns like ADD

//// ADD FUNCTIONS /////

// addBytes add and then return the addition of two bytes
func (cpu *CPU) addBytes(a, b byte) byte {
	result := a + b

	carry := (uint16(a)&0xFF)+(uint16(b)&0xFF) > 0xFF // TODO - check and maybe change

	cpu.SetFlag(N, false)
	cpu.SetFlag(Z, result == 0)
	cpu.SetFlag(H, halfCarryAdd8b(a, b))
	cpu.SetFlag(C, carry)

	return result
}

// AddRegToA add value stored in register r to A
func (cpu *CPU) AddRegToA(r *byte) {
	cpu.Reg.A = cpu.addBytes(cpu.Reg.A, *r)
	// TODO - check this is actually correct cuz idk
}

// AddRegAToHLRegData add the data located in HL reg address to reg A and set reg A to the result
func (cpu *CPU) AddRegAToHLRegData() {
	data := cpu.ReadByte(cpu.Reg.HL())
	cpu.Reg.A = cpu.addBytes(cpu.Reg.A, data)
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

// Add8BitRegToRegAWithCarry add contents of register r to register A with carry
func (cpu *CPU) Add8BitRegToRegAWithCarry(r *byte) {
	result := cpu.Reg.A + boolToBit(cpu.Reg.F.carry) + *r

	carry := (uint16(cpu.Reg.A)&0xFF)+(uint16(*r)&0xFF)+uint16(boolToBit(cpu.Reg.F.carry)) > 0xFF
	halfCarry := (cpu.Reg.A&0x0F)+(boolToBit(cpu.Reg.F.carry)&0x0F)+(*r&0x0F) > 0x0F

	cpu.SetFlag(Z, result == 0)
	cpu.SetFlag(N, false)
	cpu.SetFlag(C, carry)
	cpu.SetFlag(H, halfCarry)

	cpu.Reg.A = result
	// TODO - test the carry and half carry work
}

// Add8BitDataToRegA add the 8-bit immediate data to reg A and store result in reg A
func (cpu *CPU) Add8BitDataToRegA() {
	cpu.Reg.A = cpu.addBytes(cpu.Reg.A, cpu.CurrentInstruction.Operands[0])
	// TODO - test this is right
}

// AddHLRegDataToRegAWithCarry add the data stored in address at reg HL to reg A with carry
func (cpu *CPU) AddHLRegDataToRegAWithCarry() {
	data := cpu.ReadByte(cpu.Reg.HL())

	cpu.Add8BitRegToRegAWithCarry(&data)

	// TODO - check this
}

// AddSPToHLReg Add the stack pointer value to the HL Register
func (cpu *CPU) AddSPToHLReg() {
	result := cpu.Add16(cpu.Reg.HL(), cpu.SP)
	high, low := Split16(result)
	cpu.Reg.H, cpu.Reg.L = high, low
}

// Add8BitDataToSP add immediate 8-bit data to the stack pointer and store result as stack pointer
func (cpu *CPU) Add8BitDataToSP() {
	var result uint16
	var carry, halfCarry bool
	data := int8(cpu.CurrentInstruction.Operands[0])

	if data < 0 {
		result = cpu.SP - uint16(-data)                            // TODO - check
		carry = int(cpu.SP)-int(-data) < 0                         // TODO - gross
		halfCarry = int(cpu.SP&0x0FF)-int(uint16(-data)&0x0FF) < 0 // TODO - check
	} else {
		result = cpu.SP + uint16(data)
		carry = uint32(cpu.SP)+uint32(data) > 0xFFFF
		halfCarry = (cpu.SP&0x0FF)+(uint16(data)&0x0FF) > 0x0FF
		// TODO - check
	}

	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)
	cpu.SetFlag(C, carry)
	cpu.SetFlag(H, halfCarry)

	cpu.SP = result
}

// Add8BitDataToSPIntoHLReg add the immediate 8-bit data to the stack
// pointer value and store the result in the HL register
func (cpu *CPU) Add8BitDataToSPIntoHLReg() {
	var result uint16
	var carry, halfCarry bool
	data := int8(cpu.CurrentInstruction.Operands[0])

	if data < 0 {
		result = cpu.SP - uint16(-data)                            // TODO - check
		carry = int(cpu.SP)-int(-data) < 0                         // TODO - gross
		halfCarry = int(cpu.SP&0x0FF)-int(uint16(-data)&0x0FF) < 0 // TODO - check
	} else {
		result = cpu.SP + uint16(data)
		carry = uint32(cpu.SP)+uint32(data) > 0xFFFF
		halfCarry = (cpu.SP&0x0FF)+(uint16(data)&0x0FF) > 0x0FF
		// TODO - check
	}

	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)
	cpu.SetFlag(C, carry)
	cpu.SetFlag(H, halfCarry)

	cpu.Reg.H, cpu.Reg.L = Split16(result)
}

//// SUB FUNCTIONS /////

// Sub8BitRegFromRegA subtract reg r from reg A (A - r) and store the result in reg A
func (cpu *CPU) Sub8BitRegFromRegA(r *byte) {
	result := cpu.Reg.A - *r

	carry := int(cpu.Reg.A)-int(*r) < 0
	halfCarry := int(cpu.Reg.A&0x0F)-int(*r&0x0F) < 0

	cpu.SetFlag(N, true)
	cpu.SetFlag(Z, result == 0)
	cpu.SetFlag(C, carry)
	cpu.SetFlag(H, halfCarry)

	cpu.Reg.A = result
}

// SubHLRegDataFromRegA subtract the data stored at address in reg HL
// from reg A and store the result in reg A
func (cpu *CPU) SubHLRegDataFromRegA() {
	data := cpu.ReadByte(cpu.Reg.HL())
	cpu.Sub8BitRegFromRegA(&data)
}

// Sub8BitRegFromRegAWithCarry subtract reg r from reg A with carry (A - C - r) and store result in reg A
func (cpu *CPU) Sub8BitRegFromRegAWithCarry(r *byte) {
	carryFlag := boolToBit(cpu.Reg.F.carry)
	result := cpu.Reg.A - carryFlag - *r

	carry := int(cpu.Reg.A)-int(carryFlag)-int(*r) < 0
	halfCarry := int(cpu.Reg.A&0x0F)-int(carryFlag&0x0F)-int(*r&0x0F) < 0

	cpu.SetFlag(N, true)
	cpu.SetFlag(Z, result == 0)
	cpu.SetFlag(C, carry)
	cpu.SetFlag(H, halfCarry)

	// TODO- check this

	cpu.Reg.A = result
}

// SubHLRegDataFromRegAWithCarry subtract the data stored at the
// address in reg HL from reg A with carry and store the result in reg A
func (cpu *CPU) SubHLRegDataFromRegAWithCarry() {
	data := cpu.ReadByte(cpu.Reg.HL())
	cpu.Sub8BitRegFromRegAWithCarry(&data)
}

// Sub8BitDataFromRegA subtract the immediate 8-bit data from reg A and store the result in reg A
func (cpu *CPU) Sub8BitDataFromRegA() {
	cpu.Sub8BitRegFromRegA(&cpu.CurrentInstruction.Operands[0])
}

//// COMPARISON FUNCTIONS /////

// AndRegAWithReg execute an AND comparison between reg A and r
// Store the result in register A
func (cpu *CPU) AndRegAWithReg(r *byte) {
	result := cpu.Reg.A & *r

	cpu.SetFlag(Z, result == 0)
	cpu.SetFlag(N, false)
	cpu.SetFlag(H, true)
	cpu.SetFlag(C, false)

	cpu.Reg.A = result
}

// AndRegAWithHLData execute an AND comparison between reg A and data pointed to by HL Reg
// Store the result in register A
func (cpu *CPU) AndRegAWithHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())
	cpu.AndRegAWithReg(&data)
}

// XorRegAWithReg execute a XOR comparison between reg A and r
// Store the result in register A
func (cpu *CPU) XorRegAWithReg(r *byte) {
	result := cpu.Reg.A ^ *r

	cpu.SetFlag(Z, result == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.ResetFlag(C)

	cpu.Reg.A = result
}

// XorRegAWithHLData execute XOR comparison between reg A and HL Reg data stored in address
// Store the result in reg A
func (cpu *CPU) XorRegAWithHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())
	cpu.XorRegAWithReg(&data)
}

// OrRegAWithReg execute an OR comparison between reg A and r
// Store result in reg A
func (cpu *CPU) OrRegAWithReg(r *byte) {
	result := cpu.Reg.A | *r

	cpu.SetFlag(Z, result == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.ResetFlag(C)

	cpu.Reg.A = result
}

// OrRegAWithHLData execute an OR comparison between reg A and data stored in HL reg address
// Store result in reg A
func (cpu *CPU) OrRegAWithHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())
	cpu.OrRegAWithReg(&data)
}

// CompareRegAWithReg subtract r from reg A (A - r) and set flags accordingly
// Does not change reg A, only sets the flags
func (cpu *CPU) CompareRegAWithReg(r *byte) {
	result := cpu.Reg.A - *r

	carry := int(cpu.Reg.A)-int(*r) < 0
	halfCarry := int(cpu.Reg.A&0x0F)-int(*r&0x0F) < 0

	cpu.SetFlag(N, true)
	cpu.SetFlag(Z, result == 0)
	cpu.SetFlag(C, carry)
	cpu.SetFlag(H, halfCarry)
}

// CompareRegAWithHLData subtract (compare) data stored in HL address (A - data) and set flags
// Doesn't change reg A, only sets the flags
func (cpu *CPU) CompareRegAWithHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())
	cpu.CompareRegAWithReg(&data)
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
func (cpu *CPU) Inc16BitRegister(r *uint16) {
	*r += 1
}

// Inc16BitRegData increment the data located at address located by r1,r2 pair
func (cpu *CPU) Inc16BitRegData(r1, r2 *byte, useFlags bool) {
	address := JoinBytes(*r1, *r2)
	oldVal := cpu.ReadByte(address)
	result := oldVal + 1
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
func (cpu *CPU) IncHLRegData() {
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
func (cpu *CPU) Dec16BitRegData(r1, r2 *byte, useFlags bool) {
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
func (cpu *CPU) DecHLRegData() {
	high, low := cpu.Reg.HLByte()
	cpu.Dec16BitRegData(high, low, true)
}

// DecSPReg decrement the stack pointer register
func (cpu *CPU) DecSPReg() {
	cpu.SP -= 1
	// TODO - check I don't need to do anything else here
}

//// LD (load) FUNCTIONS /////

// Load8bRegTo8bReg load data from one 8-bit register r1 into another 8-bit register r2
// LD r, r’: Load register (register)
func (cpu *CPU) Load8BitRegInto8BitReg(r1, r2 *byte) {
	*r1 = *r2
}

// Load8bDataTo8bReg load immediate data n into 8-bit register r
// LD r, n: Load register (immediate)
func (cpu *CPU) Load8BitDataInto8BitReg(r *byte) {
	*r = cpu.CurrentInstruction.Operands[0]
}

// Load8BitRegInto16BitRegAddr load the data from 8-bit register r into the adress specified in 16-bit register r1,r2
func (cpu *CPU) Load8BitRegInto16BitRegAddr(r1, r2, r *byte) {
	regPair := JoinBytes(*r1, *r2)
	cpu.WriteByteToAddr(regPair, *r)
}

// LoadHLDataInto8BitReg load the data located in the HL register into the register R
func (cpu *CPU) LoadHLDataInto8BitReg(r *byte) {
	data := cpu.ReadByte(cpu.Reg.HL())
	*r = data
}

// Load8BitRegIntoHLAddr load the data in register r into address pointed to by HL
func (cpu *CPU) Load8BitRegIntoHLAddr(r *byte) {
	cpu.WriteByteToAddr(cpu.Reg.HL(), *r)
	// TODO - chec this is right
}

// Load8BitDataInto16BitRegAddr load the immediate data into adress specified by r1,r2
func (cpu *CPU) Load8BitDataInto16BitRegAddr(r1, r2 *byte) {
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
func (cpu *CPU) Load16BitAddrIncInto8BitReg(r1, r2, r *byte) {
	*r = cpu.ReadByte(JoinBytes(*r1, *r2))
	cpu.Inc16BitRegData(r1, r2, false)
	// TODO - check
}

// Load16BitAddrDecInto8BitReg Load the data stored in the adress pointed to by r1,r2 into reg r
// Decrement the data stored in address at r1,r2 after
func (cpu *CPU) Load16BitAddrDecInto8BitReg(r1, r2, r *byte) {
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
func (cpu *CPU) Load16BitDataInto16BitRegister(r *uint16) {
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

// LoadRegAIntoInternalRam load contents of register A into internal ram pointed to by immediate 8-bit(?) data
func (cpu *CPU) LoadRegAIntoInternalRam() {
	data := cpu.CurrentInstruction.Operands[0]
	addr := 0xFF00 + uint16(data)
	cpu.WriteByteToAddr(addr, cpu.Reg.A)
	// TODO - check this is right
}

// LoadRegAIntoRegCInternalRam load the contents of reg A into the location of internal ram pointed to by reg C
func (cpu *CPU) LoadRegAIntoRegCInternalRam() {
	addr := 0xFF00 + uint16(cpu.Reg.C)
	cpu.WriteByteToAddr(addr, cpu.Reg.A)
	// TODO - check
}

// LoadRegCInteralRamIntoRegA load the data stored in interal ram at the address pointed to by register C into reg A
func (cpu *CPU) LoadRegCInteralRamIntoRegA() {
	cpu.Reg.A = cpu.ReadByte(0xFF00 + uint16(cpu.Reg.C))
	// TODO - check
}

// LoadRegAIntoInternalRamData load the contents of register A into internal ram referenced in immediate 16-bit data
func (cpu *CPU) LoadRegAIntoInternalRamData() {
	data := JoinBytes(cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1])
	cpu.WriteByteToAddr(data, cpu.Reg.A)
	// TODO - check
}

// LoadInternalRamDataIntoRegA load the data in internal ram located at immediate data address into register A
func (cpu *CPU) LoadInternalRamDataIntoRegA() {
	data := cpu.CurrentInstruction.Operands[0]
	cpu.Reg.A = cpu.ReadByte(0xFF00 + uint16(data))
	// TODO - check
}

// LoadHLRegIntoSP load the HL register value into the stack pointer
// basically just SP = HL
func (cpu *CPU) LoadHLRegIntoSP() {
	cpu.SP = cpu.Reg.HL()
	// TODO check
}

// Load16BitDataIntoRegA load data pointed to by immediate 16-bit address into register A
func (cpu *CPU) Load16BitDataIntoRegA() {
	lsb, msb := cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	cpu.Reg.A = cpu.ReadByte(JoinBytes(msb, lsb)) // TODO - check this is right order and stuff
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

// RotateLeftReg Rotate reg r left
func (cpu *CPU) RotateLeftReg(r *byte)  {
	var bit7 bool = false

	if *r & 0x80 == 0x80 {
		bit7 = true
	}

	outcome := (*r << 1) | (boolToBit(cpu.Reg.F.carry) & 0x01)

	cpu.SetFlag(C, bit7)
	cpu.SetFlag(Z, outcome == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	*r = outcome
}

// RotateLeftHLData rotate the HL reg data left
func (cpu *CPU) RotateLeftHLData()  {
	data := cpu.ReadByte(cpu.Reg.HL())
	
	var bit7 bool = false

	if data & 0x80 == 0x80 {
		bit7 = true
	}

	outcome := (data << 1) | (boolToBit(cpu.Reg.F.carry) & 0x01)

	cpu.SetFlag(C, bit7)
	cpu.SetFlag(Z, outcome == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.WriteByteToAddr(cpu.Reg.HL(), outcome)
}

// RotateLeftCarryReg rotate given register r left through carry flag
func (cpu *CPU) RotateLeftCarryReg(r *byte)  {
	var bit7 bool = false

	if *r & 0x80 == 0x80 {
		bit7 = true
	}

	outcome := ((*r << 1) & 0xFF) | (boolToBit(bit7) & 0x01)

	
	cpu.SetFlag(C, bit7)
	cpu.SetFlag(Z, outcome == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	*r = outcome
}

// RotateLeftCarryHLData rotate the data in the HL addr left through the carry flag
func (cpu *CPU) RotateLeftCarryHLData()  {
	data := cpu.ReadByte(cpu.Reg.HL())

	var bit7 bool = false

	if data & 0x80 == 0x80 {
		bit7 = true
	}

	outcome := ((data << 1) & 0xFF) | (boolToBit(bit7) & 0x01)

	
	cpu.SetFlag(C, bit7)
	cpu.SetFlag(Z, outcome == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.WriteByteToAddr(cpu.Reg.HL(), outcome)
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
func (cpu *CPU) RotateRightRegA() {
	var bit0 bool = false

	if cpu.Reg.A&0x01 == 0x01 {
		bit0 = true
	}

	outcome := (cpu.Reg.A >> 1) | (boolToBit(cpu.Reg.F.carry) << 7 & 0x80)

	cpu.ResetAllFlags()
	cpu.SetFlag(C, bit0)
	cpu.Reg.A = outcome
}

// RotateRightReg rotate specific reg r right by one
func (cpu *CPU) RotateRightReg(r *byte)  {
	var bit0 bool = false

	if *r & 0x01 == 0x01 {
		bit0 = true
	}

	outcome := (*r >> 1) | (boolToBit(cpu.Reg.F.carry) << 7 & 0x80)

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit0)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)
	
	*r = outcome
}

// RotateRightHLData rotate the Hl reg data right
func (cpu *CPU) RotateRightHLData()  {
	data := cpu.ReadByte(cpu.Reg.HL())
	
	var bit0 bool = false

	if data & 0x01 == 0x01 {
		bit0 = true
	}

	outcome := (data >> 1) | (boolToBit(cpu.Reg.F.carry) << 7 & 0x80)

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit0)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)
	
	cpu.WriteByteToAddr(cpu.Reg.HL(), outcome)
}

// RotateRightCarryReg rotate reg r right through carry flag
func (cpu *CPU) RotateRightCarryReg(r *byte)  {
	var bit0 bool = false

	if *r & 0x01 == 0x01 {
		bit0 = true
	}

	result := (*r >> 1) | (boolToBit(bit0) << 7 & 0x80)

	cpu.SetFlag(C, bit0)
	cpu.SetFlag(Z, result == 0)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	*r = result
}

// RotateRightCarryHLData rotate the data stored in HL address right through the carry flag
func (cpu *CPU) RotateRightCarryHLData()  {
	data := cpu.ReadByte(cpu.Reg.HL())

	var bit0 bool = false

	if data & 0x01 == 0x01 {
		bit0 = true
	}

	result := (data >> 1) | (boolToBit(bit0) << 7 & 0x80)

	cpu.SetFlag(C, bit0)
	cpu.SetFlag(Z, result == 0)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	cpu.WriteByteToAddr(cpu.Reg.HL(), result)
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
func (cpu *CPU) JumpConditionalRelative8bit(cond *bool, jumpWhen bool) {
	jumpVal := int8(cpu.CurrentInstruction.Operands[0])

	if *cond == jumpWhen {
		cpu.CurrentInstruction.Instruction.Cycles = 3 // instruction has variable cylce amount
		if jumpVal < 0 {                              // stupid negative values
			cpu.PC -= uint16(jumpVal)
		} else {
			cpu.PC += uint16(jumpVal)
		}
		cpu.hasJumped = true
	} else {
		cpu.CurrentInstruction.Instruction.Cycles = 2 // just for saffety
	}
}

// JumpConditional16Bit jump to immediate 16-bit value when cond = jumpWhen
// basically PC = data when cond = jumpWhen
func (cpu *CPU) JumpConditional16Bit(cond *bool, jumpWhen bool) {
	lsb, msb := cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	data := JoinBytes(msb, lsb)

	if *cond == jumpWhen {
		cpu.PC = data
		cpu.hasJumped = true
		cpu.CurrentInstruction.Instruction.Cycles = 4
	} else {
		cpu.CurrentInstruction.Instruction.Cycles = 3
	}
}

// Jump16Bit jump to the value represented by the immediate 16-bit value
func (cpu *CPU) Jump16Bit() {
	lsb, msb := cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	cpu.PC = JoinBytes(msb, lsb)
	cpu.hasJumped = true
}

// CallFunctionConditional conditional functional call to absolute address pointed to in immediate 16-bit data
func (cpu *CPU) CallFunctionConditional(flag *bool, callWhen bool) {
	lsb, msb := cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	cpu.PC += 2 // TODO - check whether should be 2 or 3 (one for each byte read, opcode, data etc)

	if *flag == callWhen {
		cpu.pushSP(cpu.PC)
		cpu.PC = JoinBytes(msb, lsb)
		cpu.hasJumped = true
		cpu.CurrentInstruction.Instruction.Cycles = 6
	} else {
		cpu.CurrentInstruction.Instruction.Cycles = 3
	}
}

// CallFunctionUnconditional make an unconditional function call to address specified by immediate 16-bit data
// Also put the current PC on the top of the Stack pointer
func (cpu *CPU) CallFunctionUnconditional() {
	newVal := JoinBytes(cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1])
	cpu.PC += 2

	cpu.pushSP(cpu.PC)
	cpu.PC = newVal
	// TODO - check this is right
}

// Restart Unconditional function call to the absolute fixed address defined by the opcode.
// Basically just like moves you to the specific point based on the opcode - passed in as loc
func (cpu *CPU) Restart(loc byte) {
	cpu.pushSP(cpu.PC)
	cpu.PC = uint16(loc)
	cpu.hasJumped = true
}

// JumpToHLReg jump to the value stored in the HL register
// NOTE - not the data stored in the address, just the data in the register pair
func (cpu *CPU) JumpToHLReg() {
	cpu.PC = cpu.Reg.HL()
	cpu.hasJumped = true // TODO - maybeeee
}

//// SHIFT FUNCTIONS /////

// ShiftLeftReg shift reg r left
func (cpu *CPU) ShiftLeftReg(r *byte)  {
	var bit7 bool

	if *r & 0x80 == 0x80 {
		bit7 = true
	}

	outcome := *r << 1

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit7)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	*r = outcome
}

// ShiftLeftHLData shift data stored in HL reg addr left
func (cpu *CPU) ShiftLeftHLData()  {
	data := cpu.ReadByte(cpu.Reg.HL())
	var bit7 bool

	if data & 0x80 == 0x80 {
		bit7 = true
	}

	outcome := data << 1

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit7)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.WriteByteToAddr(cpu.Reg.HL(), outcome)
}

// ShiftRightReg shift reg r right
// Also keep the bit in bit7 as unchanged
func (cpu *CPU) ShiftRightReg(r *byte)  {
	var bit0 bool

	if *r & 0x01 == 0x01 {
		bit0 = true
	}

	bit7 := *r & 0x80 == 0x80

	outcome := *r >> 1 | (boolToBit(bit7) << 7) // TODO - check

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	
	*r = outcome
}

// ShiftRightHLData shift the data in HL reg addr right
// Also keep the bit in bit7 as unchanged
func (cpu *CPU) ShiftRightHLData()  {
	data := cpu.ReadByte(cpu.Reg.HL())
	var bit0 bool

	if data & 0x01 == 0x01 {
		bit0 = true
	}

	bit7 := data & 0x80 == 0x80

	outcome := data >> 1 | (boolToBit(bit7) << 7) // TODO - check

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	
	cpu.WriteByteToAddr(cpu.Reg.HL(), outcome)
}


//// MISC FUNCTIONS /////

// Nop No Operation. Doesn't do anything, only causes an increase in
// counter due to you know, the opcode being read
func (cpu *CPU) Nop() {
	// Nop-ing so hard rn
}

// Stop switches the system to STOP mode
func (cpu *CPU) Stop() {
	// TODO - figure out how to like actually implement this
}

// Halt stop system clock and enter halt mode
func (cpu *CPU) Halt() {
	cpu.Halted = true
	// TODO - implement this
}

// ComplementRegA take the complement of Reg A (flip the bits)
func (cpu *CPU) ComplementRegA() {
	cpu.Reg.A = ^cpu.Reg.A
	cpu.SetFlag(N, true)
	cpu.SetFlag(H, true)
}

// SetCarryFlag set the carry flag to true
func (cpu *CPU) SetCarryFlag() {
	cpu.SetFlag(C, true)
}

// ComplementCarryFlag make the carry flag be it's complement and resets the N and H flag
func (cpu *CPU) ComplementCarryFlag() {
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.Reg.F.carry = !cpu.Reg.F.carry
}

// ReturnConditional return from function when flag is same as
func (cpu *CPU) ReturnConditional(flag *bool, returnWhen bool) {
	if *flag == returnWhen {
		cpu.PC = cpu.popSP()
		cpu.hasJumped = true
		cpu.CurrentInstruction.Instruction.Cycles = 5
	} else {
		cpu.CurrentInstruction.Instruction.Cycles = 2
	}
}

// ReturnFromFunc unconditional return from a function
// Basically just set the program counter to be the popped off value from the Stack pointer
func (cpu *CPU) ReturnFromFunc() {
	cpu.PC = cpu.popSP()
}

// ReturnFromFuncInterrupt return from a function unconditionally and enable interrupts
func (cpu *CPU) ReturnFromFuncInterrupt() {
	cpu.PC = cpu.popSP()
	cpu.InterruptsEnabled = true
}

// PopSPIntoRegPair pop the stack pointer and put the address into reg pair r1,r2
func (cpu *CPU) PopSPIntoRegPair(r1, r2 *byte) {
	addr := cpu.popSP()
	*r2, *r1 = Split16(addr) // TODO - check the order is right here
}

// PopSPIntoAFRegPair pop the stack pointer and put the value into regpair AF
// Exists as it means I need to set the flag register properly
func (cpu *CPU) PopSPIntoAFRegPair() {
	addr := cpu.popSP()
	high, low := Split16(addr)
	cpu.Reg.A = high
	cpu.Reg.F = byteToFlagsRegister(low)
	// TODO - check this is like actually right
}

// PushRegPairOntoSP push the contents of regpair r1,r2 onto the top of the stack pointer
func (cpu *CPU) PushRegPairOntoSP(r1, r2 *byte) {
	cpu.pushSP(JoinBytes(*r1, *r2)) // TODO - check this like actually works
}

// PushAFRegOntoSP push the AF register pair onto the stack pointer
func (cpu *CPU) PushAFRegOntoSP() {
	flagReg := cpu.Reg.F.toByte()
	cpu.pushSP(JoinBytes(cpu.Reg.A, flagReg))
	// TODO - check this
}

// DisableInterrupts disable interrupts by setting interrupt flag to 0 (false)
func (cpu *CPU) DisableInterrupts() {
	cpu.InterruptsEnabled = false
	// TODO - check if don't need ot cancel anything
}

// EnableInterrupts enable interrupts
func (cpu *CPU) EnableInterrupts() {
	cpu.InterruptsEnabled = true // TODO - check this function is just this
}

// Daa decimal adjust accumulator
// TODO - I have no idea what the FUCK this does
func (cpu *CPU) Daa() {
	var a uint16 = uint16(cpu.Reg.A)

	if cpu.Reg.F.subtract {
		if cpu.Reg.F.half_carry {
			a = (a - 6) & 0xFF
		}

		if cpu.Reg.F.carry {
			a -= 0x60
		}
	} else {
		if cpu.Reg.F.half_carry || (a&0x0F) > 9 {
			a += 0x06
		}

		if cpu.Reg.F.carry || a > 0x9F {
			a += 0x60
		}
	}

	cpu.ResetFlag(H)

	if (a & 0x100) == 0x100 {
		cpu.SetFlag(C, true)
	}

	a &= 0xFF
	cpu.SetFlag(Z, a == 0)
	cpu.Reg.A = byte(a)
}

// SwapReg swap the upper four bits and lower four bits around in reg r
func (cpu *CPU) SwapReg(r *byte)  {
	btm := *r & 0x0F // the bottom 4 bits
	top := *r & 0xF0 // the top 4 bis

	result := (btm << 4) | (top >> 4)

	cpu.SetFlag(Z, result == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.ResetFlag(C)

	*r = result
}

//// INSTRUCTIONS /////

var unkownInstruction = Instruction{
	Desc:        "Unkown Instruction",
	OperandAmnt: 0,
	Cycles:      1,
	ExecFun:     func(cpu *CPU) { warnLog.Println("Unkown instruction executed.... continuing anyway") },
}

// InstructionsUnprefixed - the slice to represent each of the
// Unprefixed CPU instructions Accessed using just the index in hex
// format as that's easiest to access, then compiling the operands and
// running the associated function
var InstructionsUnprefixed []*Instruction = []*Instruction{
	&Instruction{0x00, "NOP", 0, 1, func(cpu *CPU) { cpu.Nop() }},
	&Instruction{0x01, "LD BC, d16", 2, 3, func(cpu *CPU) { cpu.Load16BitDataInto16BitRegPair(&cpu.Reg.B, &cpu.Reg.C) }},
	&Instruction{0x02, "LD (BC), A", 0, 2, func(cpu *CPU) { cpu.Load8BitRegInto16BitRegAddr(&cpu.Reg.B, &cpu.Reg.C, &cpu.Reg.A) }},
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
	&Instruction{0x12, "LD (DE), A", 0, 2, func(cpu *CPU) { cpu.Load8BitRegInto16BitRegAddr(&cpu.Reg.D, &cpu.Reg.E, &cpu.Reg.A) }},
	&Instruction{0x13, "INC DE", 0, 2, func(cpu *CPU) { cpu.Inc16BitRegPair(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0x14, "INC D", 0, 1, func(cpu *CPU) { cpu.Inc8BitReg(&cpu.Reg.D) }},
	&Instruction{0x15, "DEC D", 0, 1, func(cpu *CPU) { cpu.Dec8BitReg(&cpu.Reg.D) }},
	&Instruction{0x16, "LD D, d8", 1, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.D) }},
	&Instruction{0x17, "RLA", 0, 1, func(cpu *CPU) { cpu.RotateLeftRegA() }},
	&Instruction{0x18, "JR s8", 1, 3, func(cpu *CPU) { cpu.JumpRelative8bit() }},
	&Instruction{0x19, "ADD HL, DE", 0, 2, func(cpu *CPU) { cpu.Add16bRegToHLReg(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0x1A, "LD A, (DE)", 0, 2, func(cpu *CPU) { cpu.Load16bRegDataIntoReg(&cpu.Reg.D, &cpu.Reg.E, &cpu.Reg.A) }},
	&Instruction{0x1B, "DEC DE", 0, 2, func(cpu *CPU) { cpu.Dec16BitReg(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0x1C, "INC E", 0, 1, func(cpu *CPU) { cpu.Inc8BitReg(&cpu.Reg.E) }},
	&Instruction{0x1D, "DEC E", 0, 1, func(cpu *CPU) { cpu.Dec8BitReg(&cpu.Reg.E) }},
	&Instruction{0x1E, "LD E, d8", 1, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.E) }},
	&Instruction{0x1F, "RRA", 0, 1, func(cpu *CPU) { cpu.RotateRightRegA() }},
	&Instruction{0x20, "JR NZ, s8", 1, 2, func(cpu *CPU) { cpu.JumpConditionalRelative8bit(&cpu.Reg.F.zero, false) }},
	&Instruction{0x21, "LD HL, d16", 2, 3, func(cpu *CPU) { cpu.Load16BitDataInto16BitRegPair(&cpu.Reg.H, &cpu.Reg.L) }},
	&Instruction{0x22, "LD (HL+), A", 0, 2, func(cpu *CPU) { cpu.Load8bRegInto16bRegAddrInc(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A) }},
	&Instruction{0x23, "INC HL", 0, 2, func(cpu *CPU) { cpu.Inc16BitRegPair(&cpu.Reg.H, &cpu.Reg.L) }},
	&Instruction{0x24, "INC H", 0, 1, func(cpu *CPU) { cpu.Inc8BitReg(&cpu.Reg.H) }},
	&Instruction{0x25, "DEC H", 0, 1, func(cpu *CPU) { cpu.Dec8BitReg(&cpu.Reg.H) }},
	&Instruction{0x26, "LD H, d8", 1, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.H) }},
	&Instruction{0x27, "DAA", 0, 1, func(cpu *CPU) { cpu.Daa() }},
	&Instruction{0x28, "JR Z, s8", 1, 2, func(cpu *CPU) { cpu.JumpConditionalRelative8bit(&cpu.Reg.F.zero, true) }},
	&Instruction{0x29, "ADD HL, HL", 0, 2, func(cpu *CPU) { cpu.Add16bRegToHLReg(&cpu.Reg.H, &cpu.Reg.L) }}, // lmao
	&Instruction{0x2A, "LD A, (HL+)", 0, 2, func(cpu *CPU) { cpu.Load16BitAddrIncInto8BitReg(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A) }},
	&Instruction{0x2B, "DEC HL", 0, 2, func(cpu *CPU) { cpu.Dec16BitReg(&cpu.Reg.H, &cpu.Reg.L) }},
	&Instruction{0x2C, "INC L", 0, 1, func(cpu *CPU) { cpu.Inc8BitReg(&cpu.Reg.L) }},
	&Instruction{0x2D, "DEC L", 0, 1, func(cpu *CPU) { cpu.Dec8BitReg(&cpu.Reg.L) }},
	&Instruction{0x2E, "LD L, d8", 1, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.L) }},
	&Instruction{0x2F, "CPL", 0, 1, func(cpu *CPU) { cpu.ComplementRegA() }},
	&Instruction{0x30, "JR NC, s8", 2, 2, func(cpu *CPU) { cpu.JumpConditionalRelative8bit(&cpu.Reg.F.carry, false) }},
	&Instruction{0x31, "LD SP, d16", 2, 3, func(cpu *CPU) { cpu.Load16BitDataInto16BitRegister(&cpu.SP) }},
	&Instruction{0x32, "lD (HL-), A", 0, 2, func(cpu *CPU) { cpu.Load8bRegInto16bRegAddrDec(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A) }},
	&Instruction{0x33, "INC SP", 0, 2, func(cpu *CPU) { cpu.Inc16BitRegister(&cpu.SP) }},
	&Instruction{0x34, "INC (HL)", 0, 3, func(cpu *CPU) { cpu.IncHLRegData() }},
	&Instruction{0x35, "DEC (HL)", 0, 3, func(cpu *CPU) { cpu.DecHLRegData() }},
	&Instruction{0x36, "LD (HL), d8", 1, 3, func(cpu *CPU) { cpu.Load8BitDataInto16BitRegAddr(cpu.Reg.HLByte()) }},
	&Instruction{0x37, "SCF", 0, 1, func(cpu *CPU) { cpu.SetCarryFlag() }},
	&Instruction{0x38, "JR C, s8", 1, 2, func(cpu *CPU) { cpu.JumpConditionalRelative8bit(&cpu.Reg.F.carry, true) }},
	&Instruction{0x39, "ADD HL, SP", 0, 2, func(cpu *CPU) { cpu.AddSPToHLReg() }},
	&Instruction{0x3A, "LD A, (HL-)", 0, 2, func(cpu *CPU) { cpu.Load8bRegInto16bRegAddrDec(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A) }},
	&Instruction{0x3B, "DEC SP", 0, 2, func(cpu *CPU) { cpu.DecSPReg() }},
	&Instruction{0x3C, "INC A", 0, 1, func(cpu *CPU) { cpu.Inc8BitReg(&cpu.Reg.A) }},
	&Instruction{0x3D, "DEC A", 0, 1, func(cpu *CPU) { cpu.Dec8BitReg(&cpu.Reg.A) }},
	&Instruction{0x3E, "LD A, d8", 1, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.A) }},
	&Instruction{0x3F, "CCF", 0, 1, func(cpu *CPU) { cpu.ComplementCarryFlag() }},

	&Instruction{0x40, "LD B, B", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.B, &cpu.Reg.B) }},
	&Instruction{0x41, "LD B, C", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.B, &cpu.Reg.C) }},
	&Instruction{0x42, "LD B, D", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.B, &cpu.Reg.D) }},
	&Instruction{0x43, "LD B, E", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.B, &cpu.Reg.E) }},
	&Instruction{0x44, "LD B, H", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.B, &cpu.Reg.H) }},
	&Instruction{0x45, "LD B, L", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.B, &cpu.Reg.L) }},
	&Instruction{0x46, "LD B, (HL)", 0, 2, func(cpu *CPU) { cpu.LoadHLDataInto8BitReg(&cpu.Reg.B) }},
	&Instruction{0x47, "LD B, A", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.B, &cpu.Reg.A) }},

	&Instruction{0x48, "LD C, B", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.C, &cpu.Reg.B) }},
	&Instruction{0x49, "LD C, C", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.C, &cpu.Reg.C) }},
	&Instruction{0x4A, "LD C, D", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.C, &cpu.Reg.D) }},
	&Instruction{0x4B, "LD C, E", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.C, &cpu.Reg.E) }},
	&Instruction{0x4C, "LD C, H", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.C, &cpu.Reg.H) }},
	&Instruction{0x4D, "LD C, L", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.C, &cpu.Reg.L) }},
	&Instruction{0x4E, "LD C, (HL)", 0, 2, func(cpu *CPU) { cpu.LoadHLDataInto8BitReg(&cpu.Reg.C) }},
	&Instruction{0x4F, "LD C, A", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.C, &cpu.Reg.A) }},

	&Instruction{0x50, "LD D, B", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.D, &cpu.Reg.B) }},
	&Instruction{0x51, "LD D, C", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.D, &cpu.Reg.C) }},
	&Instruction{0x52, "LD D, D", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.D, &cpu.Reg.D) }},
	&Instruction{0x53, "LD D, E", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0x54, "LD D, H", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.D, &cpu.Reg.H) }},
	&Instruction{0x55, "LD D, L", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.D, &cpu.Reg.L) }},
	&Instruction{0x56, "LD D, (HL)", 0, 2, func(cpu *CPU) { cpu.LoadHLDataInto8BitReg(&cpu.Reg.D) }},
	&Instruction{0x57, "LD D, A", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.D, &cpu.Reg.A) }},

	&Instruction{0x58, "LD E, B", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.E, &cpu.Reg.B) }},
	&Instruction{0x59, "LD E, C", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.E, &cpu.Reg.C) }},
	&Instruction{0x5A, "LD E, D", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.E, &cpu.Reg.D) }},
	&Instruction{0x5B, "LD E, E", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.E, &cpu.Reg.E) }},
	&Instruction{0x5C, "LD E, H", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.E, &cpu.Reg.H) }},
	&Instruction{0x5D, "LD E, L", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.E, &cpu.Reg.L) }},
	&Instruction{0x5E, "LD E, (HL)", 0, 2, func(cpu *CPU) { cpu.LoadHLDataInto8BitReg(&cpu.Reg.E) }},
	&Instruction{0x5F, "LD E, A", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.E, &cpu.Reg.A) }},

	&Instruction{0x60, "LD H, B", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.H, &cpu.Reg.B) }},
	&Instruction{0x61, "LD H, C", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.H, &cpu.Reg.C) }},
	&Instruction{0x62, "LD H, D", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.H, &cpu.Reg.D) }},
	&Instruction{0x63, "LD H, E", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.H, &cpu.Reg.E) }},
	&Instruction{0x64, "LD H, H", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.H, &cpu.Reg.H) }},
	&Instruction{0x65, "LD H, L", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.H, &cpu.Reg.L) }},
	&Instruction{0x66, "LD H, (HL)", 0, 2, func(cpu *CPU) { cpu.LoadHLDataInto8BitReg(&cpu.Reg.H) }},
	&Instruction{0x67, "LD H, A", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.H, &cpu.Reg.A) }},

	&Instruction{0x68, "LD L, B", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.L, &cpu.Reg.B) }},
	&Instruction{0x69, "LD L, C", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.L, &cpu.Reg.C) }},
	&Instruction{0x6A, "LD L, D", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.L, &cpu.Reg.D) }},
	&Instruction{0x6B, "LD L, E", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.L, &cpu.Reg.E) }},
	&Instruction{0x6C, "LD L, H", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.L, &cpu.Reg.H) }},
	&Instruction{0x6D, "LD L, L", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.L, &cpu.Reg.L) }},
	&Instruction{0x6E, "LD L, (HL)", 0, 2, func(cpu *CPU) { cpu.LoadHLDataInto8BitReg(&cpu.Reg.L) }},
	&Instruction{0x6F, "LD L, A", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.L, &cpu.Reg.A) }},

	&Instruction{0x70, "LD (HL), B", 0, 1, func(cpu *CPU) { cpu.Load8BitRegIntoHLAddr(&cpu.Reg.B) }},
	&Instruction{0x71, "LD (HL), C", 0, 1, func(cpu *CPU) { cpu.Load8BitRegIntoHLAddr(&cpu.Reg.C) }},
	&Instruction{0x72, "LD (HL), D", 0, 1, func(cpu *CPU) { cpu.Load8BitRegIntoHLAddr(&cpu.Reg.D) }},
	&Instruction{0x73, "LD (HL), E", 0, 1, func(cpu *CPU) { cpu.Load8BitRegIntoHLAddr(&cpu.Reg.E) }},
	&Instruction{0x74, "LD (HL), H", 0, 1, func(cpu *CPU) { cpu.Load8BitRegIntoHLAddr(&cpu.Reg.H) }},
	&Instruction{0x75, "LD (HL), L", 0, 1, func(cpu *CPU) { cpu.Load8BitRegIntoHLAddr(&cpu.Reg.L) }},
	&Instruction{0x76, "HALT", 0, 1, func(cpu *CPU) { cpu.Halt() }},
	&Instruction{0x77, "LD (HL), A", 0, 1, func(cpu *CPU) { cpu.Load8BitRegIntoHLAddr(&cpu.Reg.A) }},

	&Instruction{0x78, "LD A, B", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.A, &cpu.Reg.B) }},
	&Instruction{0x79, "LD A, C", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.A, &cpu.Reg.C) }},
	&Instruction{0x7A, "LD A, D", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.A, &cpu.Reg.D) }},
	&Instruction{0x7B, "LD A, E", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.A, &cpu.Reg.E) }},
	&Instruction{0x7C, "LD A, H", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.A, &cpu.Reg.H) }},
	&Instruction{0x7D, "LD A, L", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.A, &cpu.Reg.L) }},
	&Instruction{0x7E, "LD A, (HL)", 0, 2, func(cpu *CPU) { cpu.LoadHLDataInto8BitReg(&cpu.Reg.A) }},
	&Instruction{0x7F, "LD A, A", 0, 1, func(cpu *CPU) { cpu.Load8BitRegInto8BitReg(&cpu.Reg.A, &cpu.Reg.A) }},

	&Instruction{0x80, "ADD A, B", 0, 1, func(cpu *CPU) { cpu.AddRegToA(&cpu.Reg.B) }},
	&Instruction{0x81, "ADD A, C", 0, 1, func(cpu *CPU) { cpu.AddRegToA(&cpu.Reg.C) }},
	&Instruction{0x82, "ADD A, D", 0, 1, func(cpu *CPU) { cpu.AddRegToA(&cpu.Reg.D) }},
	&Instruction{0x83, "ADD A, E", 0, 1, func(cpu *CPU) { cpu.AddRegToA(&cpu.Reg.E) }},
	&Instruction{0x84, "ADD A, H", 0, 1, func(cpu *CPU) { cpu.AddRegToA(&cpu.Reg.H) }},
	&Instruction{0x85, "ADD A, L", 0, 1, func(cpu *CPU) { cpu.AddRegToA(&cpu.Reg.L) }},
	&Instruction{0x86, "ADD A, (HL)", 0, 2, func(cpu *CPU) { cpu.AddRegAToHLRegData() }},
	&Instruction{0x87, "ADD A, A", 0, 1, func(cpu *CPU) { cpu.AddRegToA(&cpu.Reg.A) }},

	&Instruction{0x88, "ADC A, B", 0, 1, func(cpu *CPU) { cpu.Add8BitRegToRegAWithCarry(&cpu.Reg.B) }},
	&Instruction{0x89, "ADC A, C", 0, 1, func(cpu *CPU) { cpu.Add8BitRegToRegAWithCarry(&cpu.Reg.C) }},
	&Instruction{0x8A, "ADC A, D", 0, 1, func(cpu *CPU) { cpu.Add8BitRegToRegAWithCarry(&cpu.Reg.D) }},
	&Instruction{0x8B, "ADC A, E", 0, 1, func(cpu *CPU) { cpu.Add8BitRegToRegAWithCarry(&cpu.Reg.E) }},
	&Instruction{0x8C, "ADC A, H", 0, 1, func(cpu *CPU) { cpu.Add8BitRegToRegAWithCarry(&cpu.Reg.H) }},
	&Instruction{0x8D, "ADC A, L", 0, 1, func(cpu *CPU) { cpu.Add8BitRegToRegAWithCarry(&cpu.Reg.L) }},
	&Instruction{0x8E, "ADC A, (HL)", 0, 2, func(cpu *CPU) { cpu.AddHLRegDataToRegAWithCarry() }},
	&Instruction{0x8F, "ADC A, A", 0, 1, func(cpu *CPU) { cpu.Add8BitRegToRegAWithCarry(&cpu.Reg.A) }},

	&Instruction{0x90, "SUB A, B", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegA(&cpu.Reg.B) }},
	&Instruction{0x91, "SUB A, C", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegA(&cpu.Reg.C) }},
	&Instruction{0x92, "SUB A, D", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegA(&cpu.Reg.D) }},
	&Instruction{0x93, "SUB A, E", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegA(&cpu.Reg.E) }},
	&Instruction{0x94, "SUB A, H", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegA(&cpu.Reg.H) }},
	&Instruction{0x95, "SUB A, L", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegA(&cpu.Reg.L) }},
	&Instruction{0x96, "SUB A, (HL)", 0, 2, func(cpu *CPU) { cpu.SubHLRegDataFromRegA() }},
	&Instruction{0x97, "SUB A, A", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegA(&cpu.Reg.A) }},

	&Instruction{0x98, "SBC A, B", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegAWithCarry(&cpu.Reg.B) }},
	&Instruction{0x99, "SBC A, C", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegAWithCarry(&cpu.Reg.C) }},
	&Instruction{0x9A, "SBC A, D", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegAWithCarry(&cpu.Reg.D) }},
	&Instruction{0x9B, "SBC A, E", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegAWithCarry(&cpu.Reg.E) }},
	&Instruction{0x9C, "SBC A, H", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegAWithCarry(&cpu.Reg.H) }},
	&Instruction{0x9D, "SBC A, L", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegAWithCarry(&cpu.Reg.L) }},
	&Instruction{0x9E, "SBC A, (HL)", 0, 2, func(cpu *CPU) { cpu.SubHLRegDataFromRegAWithCarry() }},
	&Instruction{0x9F, "SBC A, A", 0, 1, func(cpu *CPU) { cpu.Sub8BitRegFromRegAWithCarry(&cpu.Reg.A) }},

	&Instruction{0xA0, "AND B", 0, 1, func(cpu *CPU) { cpu.AndRegAWithReg(&cpu.Reg.B) }},
	&Instruction{0xA1, "AND C", 0, 1, func(cpu *CPU) { cpu.AndRegAWithReg(&cpu.Reg.C) }},
	&Instruction{0xA2, "AND D", 0, 1, func(cpu *CPU) { cpu.AndRegAWithReg(&cpu.Reg.D) }},
	&Instruction{0xA3, "AND E", 0, 1, func(cpu *CPU) { cpu.AndRegAWithReg(&cpu.Reg.E) }},
	&Instruction{0xA4, "AND H", 0, 1, func(cpu *CPU) { cpu.AndRegAWithReg(&cpu.Reg.H) }},
	&Instruction{0xA5, "AND L", 0, 1, func(cpu *CPU) { cpu.AndRegAWithReg(&cpu.Reg.L) }},
	&Instruction{0xA6, "AND (HL)", 0, 2, func(cpu *CPU) { cpu.AndRegAWithHLData() }},
	&Instruction{0xA7, "AND A", 0, 1, func(cpu *CPU) { cpu.AndRegAWithReg(&cpu.Reg.A) }},

	&Instruction{0xA8, "XOR B", 0, 1, func(cpu *CPU) { cpu.XorRegAWithReg(&cpu.Reg.B) }},
	&Instruction{0xA9, "XOR C", 0, 1, func(cpu *CPU) { cpu.XorRegAWithReg(&cpu.Reg.C) }},
	&Instruction{0xAA, "XOR D", 0, 1, func(cpu *CPU) { cpu.XorRegAWithReg(&cpu.Reg.D) }},
	&Instruction{0xAB, "XOR E", 0, 1, func(cpu *CPU) { cpu.XorRegAWithReg(&cpu.Reg.E) }},
	&Instruction{0xAC, "XOR H", 0, 1, func(cpu *CPU) { cpu.XorRegAWithReg(&cpu.Reg.H) }},
	&Instruction{0xAD, "XOR L", 0, 1, func(cpu *CPU) { cpu.XorRegAWithReg(&cpu.Reg.L) }},
	&Instruction{0xAE, "XOR (HL)", 0, 2, func(cpu *CPU) { cpu.XorRegAWithHLData() }},
	&Instruction{0xAF, "XOR A", 0, 1, func(cpu *CPU) { cpu.XorRegAWithReg(&cpu.Reg.A) }},

	&Instruction{0xB0, "OR B", 0, 1, func(cpu *CPU) { cpu.OrRegAWithReg(&cpu.Reg.B) }},
	&Instruction{0xB1, "OR C", 0, 1, func(cpu *CPU) { cpu.OrRegAWithReg(&cpu.Reg.C) }},
	&Instruction{0xB2, "OR D", 0, 1, func(cpu *CPU) { cpu.OrRegAWithReg(&cpu.Reg.D) }},
	&Instruction{0xB3, "OR E", 0, 1, func(cpu *CPU) { cpu.OrRegAWithReg(&cpu.Reg.E) }},
	&Instruction{0xB4, "OR H", 0, 1, func(cpu *CPU) { cpu.OrRegAWithReg(&cpu.Reg.H) }},
	&Instruction{0xB5, "OR L", 0, 1, func(cpu *CPU) { cpu.OrRegAWithReg(&cpu.Reg.L) }},
	&Instruction{0xB6, "OR (HL)", 0, 2, func(cpu *CPU) { cpu.OrRegAWithHLData() }},
	&Instruction{0xB7, "OR A", 0, 1, func(cpu *CPU) { cpu.OrRegAWithReg(&cpu.Reg.A) }},

	&Instruction{0xB8, "CP B", 0, 1, func(cpu *CPU) { cpu.CompareRegAWithReg(&cpu.Reg.B) }},
	&Instruction{0xB9, "CP C", 0, 1, func(cpu *CPU) { cpu.CompareRegAWithReg(&cpu.Reg.C) }},
	&Instruction{0xBA, "CP D", 0, 1, func(cpu *CPU) { cpu.CompareRegAWithReg(&cpu.Reg.D) }},
	&Instruction{0xBB, "CP E", 0, 1, func(cpu *CPU) { cpu.CompareRegAWithReg(&cpu.Reg.E) }},
	&Instruction{0xBC, "CP H", 0, 1, func(cpu *CPU) { cpu.CompareRegAWithReg(&cpu.Reg.H) }},
	&Instruction{0xBD, "CP L", 0, 1, func(cpu *CPU) { cpu.CompareRegAWithReg(&cpu.Reg.L) }},
	&Instruction{0xBE, "CP (HL)", 0, 2, func(cpu *CPU) { cpu.CompareRegAWithHLData() }},
	&Instruction{0xBF, "CP A", 0, 1, func(cpu *CPU) { cpu.CompareRegAWithReg(&cpu.Reg.A) }},

	&Instruction{0xC0, "RET NZ", 0, 2, func(cpu *CPU) { cpu.ReturnConditional(&cpu.Reg.F.zero, false) }},
	&Instruction{0xC1, "POP BC", 0, 3, func(cpu *CPU) { cpu.PopSPIntoRegPair(&cpu.Reg.B, &cpu.Reg.C) }},
	&Instruction{0xC2, "JP NZ, a16", 2, 3, func(cpu *CPU) { cpu.JumpConditional16Bit(&cpu.Reg.F.zero, false) }},
	&Instruction{0xC3, "JP a16", 2, 4, func(cpu *CPU) { cpu.Jump16Bit() }},
	&Instruction{0xC4, "CALL NZ, a16", 2, 3, func(cpu *CPU) { cpu.CallFunctionConditional(&cpu.Reg.F.zero, false) }},
	&Instruction{0xC5, "PUSH BC", 0, 4, func(cpu *CPU) { cpu.PushRegPairOntoSP(&cpu.Reg.B, &cpu.Reg.C) }},
	&Instruction{0xC6, "ADD A, d8", 1, 2, func(cpu *CPU) { cpu.Add8BitDataToRegA() }},
	&Instruction{0xC7, "RST 0", 0, 4, func(cpu *CPU) { cpu.Restart(0x00) }},
	&Instruction{0xC8, "RET Z", 0, 2, func(cpu *CPU) { cpu.ReturnConditional(&cpu.Reg.F.zero, true) }},
	&Instruction{0xC9, "RET", 0, 4, func(cpu *CPU) { cpu.ReturnFromFunc() }},
	&Instruction{0xCA, "JP Z, a16", 2, 3, func(cpu *CPU) { cpu.JumpConditional16Bit(&cpu.Reg.F.zero, true) }},
	&unkownInstruction,
	&Instruction{0xCC, "CALL Z, a16", 2, 3, func(cpu *CPU) { cpu.CallFunctionConditional(&cpu.Reg.F.zero, true) }},
	&Instruction{0xCD, "CALL a16", 2, 6, func(cpu *CPU) { cpu.CallFunctionUnconditional() }},
	&Instruction{0xCE, "ADC A, d8", 1, 2, func(cpu *CPU) { cpu.Add8BitRegToRegAWithCarry(&cpu.CurrentInstruction.Operands[0]) }}, // TODO - check
	&Instruction{0xCF, "RST 1", 0, 4, func(cpu *CPU) { cpu.Restart(0x08) }},
	&Instruction{0xD0, "RET NC", 0, 2, func(cpu *CPU) { cpu.ReturnConditional(&cpu.Reg.F.zero, false) }},
	&Instruction{0xD1, "POP DE", 0, 3, func(cpu *CPU) { cpu.PopSPIntoRegPair(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0xD2, "JP NC, a16", 2, 3, func(cpu *CPU) { cpu.JumpConditional16Bit(&cpu.Reg.F.carry, false) }},
	&unkownInstruction,
	&Instruction{0xD4, "CALL NC, a16", 2, 3, func(cpu *CPU) { cpu.CallFunctionConditional(&cpu.Reg.F.carry, false) }},
	&Instruction{0xD5, "PUSH DE", 0, 4, func(cpu *CPU) { cpu.PushRegPairOntoSP(&cpu.Reg.D, &cpu.Reg.E) }},
	&Instruction{0xD6, "SUB d8", 1, 2, func(cpu *CPU) { cpu.Sub8BitDataFromRegA() }},
	&Instruction{0xD7, "RST 2", 0, 4, func(cpu *CPU) { cpu.Restart(0x10) }},
	&Instruction{0xD8, "RET C", 0, 2, func(cpu *CPU) { cpu.ReturnConditional(&cpu.Reg.F.carry, true) }},
	&Instruction{0xD9, "RETI", 0, 4, func(cpu *CPU) { cpu.ReturnFromFuncInterrupt() }},
	&Instruction{0xDA, "JP C, a16", 2, 3, func(cpu *CPU) { cpu.JumpConditional16Bit(&cpu.Reg.F.carry, true) }},
	&unkownInstruction,
	&Instruction{0xDC, "CALL C, a16", 2, 3, func(cpu *CPU) { cpu.CallFunctionConditional(&cpu.Reg.F.carry, true) }},
	&unkownInstruction,
	&Instruction{0xDE, "SBC A, d8", 1, 2, func(cpu *CPU) { cpu.Sub8BitRegFromRegAWithCarry(&cpu.CurrentInstruction.Operands[0]) }}, // TODO - check
	&Instruction{0xDF, "RST 3", 0, 4, func(cpu *CPU) { cpu.Restart(0x18) }},
	&Instruction{0xE0, "LD (a8), A", 1, 3, func(cpu *CPU) { cpu.LoadRegAIntoInternalRam() }},
	&Instruction{0xE1, "POP HL", 0, 3, func(cpu *CPU) { cpu.PopSPIntoRegPair(cpu.Reg.HLByte()) }}, // TODO - check
	&Instruction{0xE2, "LD (C), A", 0, 2, func(cpu *CPU) { cpu.LoadRegAIntoRegCInternalRam() }},
	&unkownInstruction,
	&unkownInstruction,
	&Instruction{0xE5, "PUSH HL", 0, 4, func(cpu *CPU) { cpu.PushRegPairOntoSP(&cpu.Reg.H, &cpu.Reg.L) }},
	&Instruction{0xE6, "AND d8", 1, 2, func(cpu *CPU) { cpu.AndRegAWithReg(&cpu.CurrentInstruction.Operands[0]) }}, //TODO - check
	&Instruction{0xE7, "RST 4", 0, 4, func(cpu *CPU) { cpu.Restart(0x20) }},
	&Instruction{0xE8, "ADD SP, s8", 1, 4, func(cpu *CPU) { cpu.Add8BitDataToSP() }},
	&Instruction{0xE9, "JP HL", 0, 1, func(cpu *CPU) { cpu.JumpToHLReg() }},
	&Instruction{0xEA, "LD (a16), A", 2, 4, func(cpu *CPU) { cpu.LoadRegAIntoInternalRamData() }},
	&unkownInstruction,
	&unkownInstruction,
	&unkownInstruction,
	&Instruction{0xEE, "XOR d8", 1, 2, func(cpu *CPU) { cpu.XorRegAWithReg(&cpu.CurrentInstruction.Operands[0]) }}, // TODO - check
	&Instruction{0xEF, "RST 5", 0, 4, func(cpu *CPU) { cpu.Restart(0x28) }},
	&Instruction{0xF0, "LD A, (a8)", 1, 3, func(cpu *CPU) { cpu.LoadInternalRamDataIntoRegA() }},
	&Instruction{0xF1, "POP AF", 0, 3, func(cpu *CPU) { cpu.PopSPIntoAFRegPair() }},
	&Instruction{0xF2, "LD A, (C)", 0, 2, func(cpu *CPU) { cpu.LoadRegCInteralRamIntoRegA() }},
	&Instruction{0xF3, "DI", 0, 1, func(cpu *CPU) { cpu.DisableInterrupts() }},
	&unkownInstruction,
	&Instruction{0xF5, "PUSH AF", 0, 4, func(cpu *CPU) { cpu.PushAFRegOntoSP() }},
	&Instruction{0xF6, "OR d8", 1, 2, func(cpu *CPU) { cpu.OrRegAWithReg(&cpu.CurrentInstruction.Operands[0]) }},
	&Instruction{0xF7, "RST 6", 0, 4, func(cpu *CPU) { cpu.Restart(0x30) }},
	&Instruction{0xF8, "LD HL, SP+s8", 1, 3, func(cpu *CPU) { cpu.Add8BitDataToSPIntoHLReg() }},
	&Instruction{0xF9, "LD SP, HL", 0, 2, func(cpu *CPU) { cpu.LoadHLRegIntoSP() }},
	&Instruction{0xFA, "LD A, (a16)", 2, 4, func(cpu *CPU) { cpu.Load16BitDataIntoRegA() }},
	&Instruction{0xFB, "EI", 0, 1, func(cpu *CPU) { cpu.EnableInterrupts() }},
	&unkownInstruction,
	&unkownInstruction,
	&Instruction{0xFE, "CP d8", 1, 2, func(cpu *CPU) { cpu.CompareRegAWithReg(&cpu.CurrentInstruction.Operands[0]) }}, // TODO - check
	&Instruction{0xFF, "RST 7", 0, 4, func(cpu *CPU) { cpu.Restart(0x38) }},
}

var InstructionsPrefixed []*Instruction = []*Instruction{
	&Instruction{0x00, "RLC B", 0, 2, func(cpu *CPU){cpu.RotateLeftCarryReg(&cpu.Reg.B)}},
	&Instruction{0x01, "RLC C", 0, 2, func(cpu *CPU){cpu.RotateLeftCarryReg(&cpu.Reg.C)}},
	&Instruction{0x02, "RLC D", 0, 2, func(cpu *CPU){cpu.RotateLeftCarryReg(&cpu.Reg.D)}},
	&Instruction{0x03, "RLC E", 0, 2, func(cpu *CPU){cpu.RotateLeftCarryReg(&cpu.Reg.E)}},
	&Instruction{0x04, "RLC H", 0, 2, func(cpu *CPU){cpu.RotateLeftCarryReg(&cpu.Reg.H)}},
	&Instruction{0x05, "RLC L", 0, 2, func(cpu *CPU){cpu.RotateLeftCarryReg(&cpu.Reg.L)}},
	&Instruction{0x06, "RLC (HL)", 0, 4, func(cpu *CPU){cpu.RotateLeftCarryHLData()}},
	&Instruction{0x07, "RLC A", 0, 2, func(cpu *CPU){cpu.RotateLeftCarryReg(&cpu.Reg.A)}},
	
	&Instruction{0x08, "RRC B", 0, 2, func(cpu *CPU){cpu.RotateRightCarryReg(&cpu.Reg.B)}},
	&Instruction{0x09, "RRC C", 0, 2, func(cpu *CPU){cpu.RotateRightCarryReg(&cpu.Reg.C)}},
	&Instruction{0x0A, "RRC D", 0, 2, func(cpu *CPU){cpu.RotateRightCarryReg(&cpu.Reg.D)}},
	&Instruction{0x0B, "RRC E", 0, 2, func(cpu *CPU){cpu.RotateRightCarryReg(&cpu.Reg.E)}},
	&Instruction{0x0C, "RRC H", 0, 2, func(cpu *CPU){cpu.RotateRightCarryReg(&cpu.Reg.H)}},
	&Instruction{0x0D, "RRC L", 0, 2, func(cpu *CPU){cpu.RotateRightCarryReg(&cpu.Reg.L)}},
	&Instruction{0x0E, "RRC (HL)", 0, 4, func(cpu *CPU){cpu.RotateRightCarryHLData()}},    
	&Instruction{0x0F, "RRC A", 0, 2, func(cpu *CPU){cpu.RotateRightCarryReg(&cpu.Reg.A)}},
	
	&Instruction{0x10, "RL B", 0, 2, func(cpu *CPU){cpu.RotateLeftReg(&cpu.Reg.B)}}, 
	&Instruction{0x11, "RL C", 0, 2, func(cpu *CPU){cpu.RotateLeftReg(&cpu.Reg.C)}}, 
	&Instruction{0x12, "RL D", 0, 2, func(cpu *CPU){cpu.RotateLeftReg(&cpu.Reg.D)}}, 
	&Instruction{0x13, "RL E", 0, 2, func(cpu *CPU){cpu.RotateLeftReg(&cpu.Reg.E)}}, 
	&Instruction{0x14, "RL H", 0, 2, func(cpu *CPU){cpu.RotateLeftReg(&cpu.Reg.H)}}, 
	&Instruction{0x15, "RL L", 0, 2, func(cpu *CPU){cpu.RotateLeftReg(&cpu.Reg.L)}}, 
	&Instruction{0x16, "RL (HL)", 0, 4, func(cpu *CPU){cpu.RotateLeftHLData()}},     
	&Instruction{0x17, "RL A", 0, 2, func(cpu *CPU){cpu.RotateLeftReg(&cpu.Reg.A)}}, 
					                                                                      
	&Instruction{0x18, "RR B", 0, 2, func(cpu *CPU){cpu.RotateRightReg(&cpu.Reg.B)}},
	&Instruction{0x19, "RR C", 0, 2, func(cpu *CPU){cpu.RotateRightReg(&cpu.Reg.C)}},
	&Instruction{0x1A, "RR D", 0, 2, func(cpu *CPU){cpu.RotateRightReg(&cpu.Reg.D)}},
	&Instruction{0x1B, "RR E", 0, 2, func(cpu *CPU){cpu.RotateRightReg(&cpu.Reg.E)}},
	&Instruction{0x1C, "RR H", 0, 2, func(cpu *CPU){cpu.RotateRightReg(&cpu.Reg.H)}},
	&Instruction{0x1D, "RR L", 0, 2, func(cpu *CPU){cpu.RotateRightReg(&cpu.Reg.L)}},
	&Instruction{0x1E, "RR (HL)", 0, 4, func(cpu *CPU){cpu.RotateRightHLData()}},    
	&Instruction{0x1F, "RR A", 0, 2, func(cpu *CPU){cpu.RotateRightReg(&cpu.Reg.A)}},
	
	&Instruction{0x20, "SLA B", 0, 2, func(cpu *CPU){cpu.ShiftLeftReg(&cpu.Reg.B)}}, 
	&Instruction{0x21, "SLA C", 0, 2, func(cpu *CPU){cpu.ShiftLeftReg(&cpu.Reg.C)}}, 
	&Instruction{0x22, "SLA D", 0, 2, func(cpu *CPU){cpu.ShiftLeftReg(&cpu.Reg.D)}}, 
	&Instruction{0x23, "SLA E", 0, 2, func(cpu *CPU){cpu.ShiftLeftReg(&cpu.Reg.E)}}, 
	&Instruction{0x24, "SLA H", 0, 2, func(cpu *CPU){cpu.ShiftLeftReg(&cpu.Reg.H)}}, 
	&Instruction{0x25, "SLA L", 0, 2, func(cpu *CPU){cpu.ShiftLeftReg(&cpu.Reg.L)}}, 
	&Instruction{0x26, "SLA (HL)", 0, 4, func(cpu *CPU){cpu.ShiftLeftHLData()}},     
	&Instruction{0x27, "SLA A", 0, 2, func(cpu *CPU){cpu.ShiftLeftReg(&cpu.Reg.A)}}, 
					                                                                 
	&Instruction{0x28, "SRA B", 0, 2, func(cpu *CPU){cpu.ShiftRightReg(&cpu.Reg.B)}},
	&Instruction{0x29, "SRA C", 0, 2, func(cpu *CPU){cpu.ShiftRightReg(&cpu.Reg.C)}},
	&Instruction{0x2A, "SRA D", 0, 2, func(cpu *CPU){cpu.ShiftRightReg(&cpu.Reg.D)}},
	&Instruction{0x2B, "SRA E", 0, 2, func(cpu *CPU){cpu.ShiftRightReg(&cpu.Reg.E)}},
	&Instruction{0x2C, "SRA H", 0, 2, func(cpu *CPU){cpu.ShiftRightReg(&cpu.Reg.H)}},
	&Instruction{0x2D, "SRA L", 0, 2, func(cpu *CPU){cpu.ShiftRightReg(&cpu.Reg.L)}},
	&Instruction{0x2E, "SRA (HL)", 0, 4, func(cpu *CPU){cpu.ShiftRightHLData()}},    
	&Instruction{0x2F, "SRA A", 0, 2, func(cpu *CPU){cpu.ShiftRightReg(&cpu.Reg.A)}},
	
	&Instruction{0x30, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x31, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x32, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x33, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x34, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x35, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x36, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x37, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x38, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x39, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x3A, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x3B, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x3C, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x3D, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x3E, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x3F, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x40, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x41, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x42, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x43, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x44, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x45, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x46, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x47, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x48, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x49, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x4A, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x4B, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x4C, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x4D, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x4E, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x4F, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x50, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x51, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x52, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x53, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x54, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x55, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x56, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x57, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x58, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x59, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x5A, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x5B, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x5C, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x5D, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x5E, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x5F, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x60, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x61, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x62, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x63, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x64, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x65, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x66, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x67, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x68, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x69, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x6A, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x6B, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x6C, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x6D, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x6E, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x6F, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x70, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x71, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x72, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x73, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x74, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x75, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x76, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x77, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x78, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x79, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x7A, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x7B, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x7C, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x7D, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x7E, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x7F, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x80, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x81, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x82, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x83, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x84, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x85, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x86, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x87, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x88, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x89, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x8A, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x8B, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x8C, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x8D, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x8E, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x8F, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x90, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x91, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x92, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x93, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x94, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x95, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x96, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x97, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x98, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x99, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x9A, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x9B, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x9C, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x9D, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x9E, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0x9F, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA0, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA1, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA2, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA3, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA4, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA5, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA6, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA7, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA8, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xA9, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xAA, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xAB, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xAC, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xAD, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xAE, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xAF, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB0, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB1, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB2, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB3, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB4, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB5, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB6, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB7, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB8, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xB9, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xBA, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xBB, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xBC, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xBD, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xBE, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xBF, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC0, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC1, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC2, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC3, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC4, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC5, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC6, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC7, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC8, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xC9, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xCA, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xCB, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xCC, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xCD, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xCE, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xCF, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD0, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD1, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD2, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD3, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD4, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD5, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD6, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD7, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD8, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xD9, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xDA, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xDB, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xDC, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xDD, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xDE, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xDF, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE0, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE1, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE2, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE3, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE4, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE5, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE6, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE7, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE8, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xE9, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xEA, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xEB, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xEC, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xED, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xEE, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xEF, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF0, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF1, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF2, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF3, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF4, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF5, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF6, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF7, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF8, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xF9, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xFA, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xFB, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xFC, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xFD, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xFE, "", 0, 0, func(cpu *CPU){}},
	&Instruction{0xFF, "", 0, 0, func(cpu *CPU){}},
}
