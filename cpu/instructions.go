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
// TODO - make sure that all the functions properly assign the results to things in thigns like ADD3
// TODO - check that I have the order of my High and Low splits correct
// TODO - maybe switch off using the operands system and just load the stuff when needed

//// ADD FUNCTIONS /////

// addBytes add and then return the addition of two bytes
func (cpu *CPU) addBytes(a, b byte) byte {
	result := a + b

	carry := (uint16(a))+(uint16(b)) > 0xFF 

	cpu.SetFlag(Z, result == 0)
	cpu.SetFlag(N, false)
	cpu.SetFlag(H, halfCarryAdd8b(a, b))
	cpu.SetFlag(C, carry)

	return result
}

// AddRegToA add value stored in register r to A
func (cpu *CPU) AddRegToA(r *byte) {
	cpu.Reg.A = cpu.addBytes(cpu.Reg.A, *r)
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

	carry := (uint16(cpu.Reg.A))+(uint16(*r))+uint16(boolToBit(cpu.Reg.F.carry)) > 0xFF
	halfCarry := (cpu.Reg.A&0x0F)+(boolToBit(cpu.Reg.F.carry)&0x0F)+(*r&0x0F) > 0x0F

	cpu.SetFlag(Z, result == 0)
	cpu.SetFlag(N, false)
	cpu.SetFlag(C, carry)
	cpu.SetFlag(H, halfCarry)

	cpu.Reg.A = result
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
}

// AddSPToHLReg Add the stack pointer value to the HL Register
func (cpu *CPU) AddSPToHLReg() {
	oldZero := cpu.Reg.F.zero
	result := cpu.Add16(cpu.Reg.HL(), cpu.SP)
	high, low := Split16(result)
	cpu.Reg.H, cpu.Reg.L = high, low
	cpu.SetFlag(Z, oldZero)
}

// Add8BitDataToSP add immediate 8-bit data to the stack pointer and store result as stack pointer
// NOTE - for some dumbass reason the stupid thing only needs you to test the bottom 8 bits for the flags
// https://discord.com/channels/465585922579103744/465586075830845475/1199883769759404053
func (cpu *CPU) Add8BitDataToSP() {
	var result uint16
	var carry, halfCarry bool
	data := int8(cpu.CurrentInstruction.Operands[0])
	dataUn := uint16(data)

	result = cpu.SP + uint16(data)
	carry = carryAdd8b(byte(cpu.SP), byte(dataUn))
	halfCarry = halfCarryAdd8b(byte(cpu.SP), byte(dataUn))

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
	data16 := uint16(data)
	
	result = cpu.SP + uint16(data)
	carry = carryAdd8b(byte(cpu.SP), byte(data16))
	halfCarry = halfCarryAdd8b(byte(cpu.SP), byte(data16))

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
}

// Load8BitDataInto16BitRegAddr load the immediate data into adress specified by r1,r2
func (cpu *CPU) Load8BitDataInto16BitRegAddr(r1, r2 *byte) {
	data := cpu.CurrentInstruction.Operands[0]
	addr := JoinBytes(*r1, *r2)
	cpu.WriteByteToAddr(addr, data)
}

// Load8bRegInto16bRegAddrInc load the data from 8-bit register r
// into the adress specified in 16-bit register r1,r2 and then
// increments the regPair after
func (cpu *CPU) Load8bRegInto16bRegAddrInc(r1, r2, r *byte) {
	regPair := JoinBytes(*r1, *r2)
	cpu.WriteByteToAddr(regPair, *r)
	*r1, *r2 = Split16(regPair+1)
}

// Load8bRegInto16bRegAddrDec load the data from 8-bit register r into
// the adress specified in 16-bit register r1,r2 and then Decrement
// the regPair after
func (cpu *CPU) Load8bRegInto16bRegAddrDec(r1, r2, r *byte) {
	regPair := JoinBytes(*r1, *r2)
	cpu.WriteByteToAddr(regPair, *r)
	*r1, *r2 = Split16(regPair-1)
}

// Load16BitAddrIncInto8BitReg Load the data stored in the adress pointed to by r1,r2 into reg r
// Increment the HL register address after, NOT the data
func (cpu *CPU) Load16BitAddrIncInto8BitReg(r1, r2, r *byte) {
	addr := JoinBytes(*r1, *r2)
	data := cpu.ReadByte(addr)
	*r = data
	*r1, *r2 = Split16(addr+1)
}

// Load16BitAddrDecInto8BitReg Load the data stored in the adress pointed to by r1,r2 into reg r
// Decrement the data stored in address at r1,r2 after
func (cpu *CPU) Load16BitAddrDecInto8BitReg(r1, r2, r *byte) {
	addr := JoinBytes(*r1, *r2)
	*r = cpu.ReadByte(addr)
	*r1, *r2 = Split16(addr-1)
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

// LoadStackPointerInto16bData load the data from the stack pointer into the adress specified in immediate 16bit data
func (cpu *CPU) LoadStackPointerInto16bData() {
	var low, high byte = cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	address := JoinBytes(high, low)

	cpu.WriteByteToAddr(address, byte(cpu.SP&0xFF))
	cpu.WriteByteToAddr((address+1)&0xFFFF, byte(cpu.SP>>8))
}

// LoadRegAIntoInternalRam load contents of register A into internal ram pointed to by immediate 8-bit(?) data
func (cpu *CPU) LoadRegAIntoInternalRam() {
	data := cpu.CurrentInstruction.Operands[0]
	addr := 0xFF00 + uint16(data)
	cpu.WriteByteToAddr(addr, cpu.Reg.A)
	// TODO - check this is right
}

// LoadRegAIntoRegCInternalRam load the contents of reg A into the location of internal ram pointed to by reg C
func (cpu *CPU) LoadRegAIntoRegCInternalRam() { // TODO - check IMPORTANT
	addr := 0xFF00 + uint16(cpu.Reg.C)
	cpu.WriteByteToAddr(addr, cpu.Reg.A)
	// TODO - check
}

// LoadRegCInteralRamIntoRegA load the data stored in interal ram at the address pointed to by register C into reg A
func (cpu *CPU) LoadRegCInteralRamIntoRegA() { // TODO - check IMPORTANT
	cpu.Reg.A = cpu.ReadByte(0xFF00 + uint16(cpu.Reg.C))
	// TODO - check
}

// LoadRegAIntoInternalRamData load the contents of register A into internal ram referenced in immediate 16-bit data
func (cpu *CPU) LoadRegAIntoInternalRamData() {
	lsb, msb := cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	addr := JoinBytes(msb, lsb)
	cpu.WriteByteToAddr(addr, cpu.Reg.A)
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
	addr := JoinBytes(msb, lsb)
	cpu.Reg.A = cpu.ReadByte(addr)
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
func (cpu *CPU) RotateLeftReg(r *byte) {
	var bit7 bool = false

	if *r&0x80 == 0x80 {
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
func (cpu *CPU) RotateLeftHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())

	var bit7 bool = false

	if data&0x80 == 0x80 {
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
func (cpu *CPU) RotateLeftCarryReg(r *byte) {
	var bit7 bool = false

	if *r&0x80 == 0x80 {
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
func (cpu *CPU) RotateLeftCarryHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())

	var bit7 bool = false

	if data&0x80 == 0x80 {
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
func (cpu *CPU) RotateRightReg(r *byte) {
	var bit0 bool = false

	if *r&0x01 == 0x01 {
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
func (cpu *CPU) RotateRightHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())

	var bit0 bool = false

	if data&0x01 == 0x01 {
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
func (cpu *CPU) RotateRightCarryReg(r *byte) {
	var bit0 bool = false

	if *r&0x01 == 0x01 {
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
func (cpu *CPU) RotateRightCarryHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())

	var bit0 bool = false

	if data&0x01 == 0x01 {
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
	jumpVal := cpu.CurrentInstruction.Operands[0] // convert to signed value, should hopefully overflow to correct negative

	if jumpVal > 127 { // stupid negative values
		cpu.PC -= uint16(-jumpVal)
	} else {
		cpu.PC += uint16(jumpVal)
	}

	cpu.hasJumped = true
}

// JumpConditionalRelative8bit jump relative amount when cond is true
// Can be used for both not and regular jumps
func (cpu *CPU) JumpConditionalRelative8bit(cond *bool, jumpWhen bool) {
	jumpVal := cpu.CurrentInstruction.Operands[0]

	if *cond == jumpWhen {
		cpu.CurrentInstruction.Instruction.Cycles = 3 // instruction has variable cylce amount
		if jumpVal > 127 {                              // stupid negative values
			cpu.PC -= uint16(-jumpVal)
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
	lsb, msb := cpu.CurrentInstruction.Operands[0], cpu.CurrentInstruction.Operands[1]
	newVal := JoinBytes(msb, lsb)
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
func (cpu *CPU) ShiftLeftReg(r *byte) {
	var bit7 bool

	if *r&0x80 == 0x80 {
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
func (cpu *CPU) ShiftLeftHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())
	var bit7 bool

	if data&0x80 == 0x80 {
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
func (cpu *CPU) ShiftRightReg(r *byte) {
	var bit0 bool

	if *r&0x01 == 0x01 {
		bit0 = true
	}

	bit7 := *r&0x80 == 0x80

	outcome := *r>>1 | (boolToBit(bit7) << 7) // TODO - check

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	*r = outcome
}

// ShiftRightHLData shift the data in HL reg addr right
// Also keep the bit in bit7 as unchanged
func (cpu *CPU) ShiftRightHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())
	var bit0 bool

	if data&0x01 == 0x01 {
		bit0 = true
	}

	bit7 := data&0x80 == 0x80

	outcome := data>>1 | (boolToBit(bit7) << 7) // TODO - check

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.WriteByteToAddr(cpu.Reg.HL(), outcome)
}

// ShiftRightRegLogical Shift reg r right with a logical shift
// bit7 gets reset to 0
func (cpu *CPU) ShiftRightRegLogical(r *byte) {
	var bit0 bool

	if *r&0x01 == 0x01 {
		bit0 = true
	}

	outcome := *r >> 1 // TODO - check

	cpu.SetFlag(Z, outcome == 0)
	cpu.SetFlag(C, bit0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	*r = outcome
}

// ShiftRightHLDataLogical shift the data pointed to by the HL reg right with a logical shift
func (cpu *CPU) ShiftRightHLDataLogical() {
	data := cpu.ReadByte(cpu.Reg.HL())
	var bit0 bool

	if data&0x01 == 0x01 {
		bit0 = true
	}

	outcome := data >> 1

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
// for some reason also sets N and H to false
func (cpu *CPU) SetCarryFlag() {
	cpu.SetFlag(N, false)
	cpu.SetFlag(H, false)
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
	*r1, *r2 = Split16(addr) // TODO - check the order is right here
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
func (cpu *CPU) SwapReg(r *byte) {
	btm := *r & 0x0F // the bottom 4 bits
	top := *r & 0xF0 // the top 4 bis

	result := (btm << 4) | (top >> 4)

	cpu.SetFlag(Z, result == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.ResetFlag(C)

	*r = result
}

// SwapHLData swap the nibbles around in the data stored in the HL Reg and write
func (cpu *CPU) SwapHLData() {
	data := cpu.ReadByte(cpu.Reg.HL())

	btm := data & 0x0F // the bottom 4 bits
	top := data & 0xF0 // the top 4 bis

	result := (btm << 4) | (top >> 4)

	cpu.SetFlag(Z, result == 0)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.ResetFlag(C)

	cpu.WriteByteToAddr(cpu.Reg.HL(), result)
}

// BitTestReg test bit in specified position in reg r
func (cpu *CPU) BitTestReg(bit byte, r *byte) {
	cpu.SetFlag(Z, *r>>bit&1 == 0)
	cpu.SetFlag(H, true)
	cpu.ResetFlag(N)
}

// BitTestHLData test bit in specified position in HL data
func (cpu *CPU) BitTestHLData(bit byte) {
	data := cpu.ReadByte(cpu.Reg.HL())
	cpu.BitTestReg(bit, &data) /// NOTE - maybe not best to do this and should just rewrite
}

// BitResetReg reset specified bit in reg r to 0
func (cpu *CPU) BitResetReg(bit byte, r *byte) {
	result := *r & ^(1 << bit)
	*r = result
}

// BitResetHLData Reset specified bit in HL Data to 0 and write
func (cpu *CPU) BitResetHLData(bit byte) {
	data := cpu.ReadByte(cpu.Reg.HL())
	result := data & ^(1 << bit)
	cpu.WriteByteToAddr(cpu.Reg.HL(), result)
}

// BitSetReg set specified bit in reg r to 1
func (cpu *CPU) BitSetReg(bit byte, r *byte) {
	result := *r | (1 << bit)
	*r = result
}

// BitSetHLData set specified bit in HL Data to 1
func (cpu *CPU) BitSetHLData(bit byte) {
	data := cpu.ReadByte(cpu.Reg.HL())
	result := data | (1 << bit)
	cpu.WriteByteToAddr(cpu.Reg.HL(), result)
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
	&Instruction{0x06, "LD B d8", 1, 2, func(cpu *CPU) { cpu.Load8BitDataInto8BitReg(&cpu.Reg.B) }},
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
	&Instruction{0x30, "JR NC, s8", 1, 2, func(cpu *CPU) { cpu.JumpConditionalRelative8bit(&cpu.Reg.F.carry, false) }},
	&Instruction{0x31, "LD SP, d16", 2, 3, func(cpu *CPU) { cpu.Load16BitDataInto16BitRegister(&cpu.SP) }},
	&Instruction{0x32, "lD (HL-), A", 0, 2, func(cpu *CPU) { cpu.Load8bRegInto16bRegAddrDec(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A) }},
	&Instruction{0x33, "INC SP", 0, 2, func(cpu *CPU) { cpu.Inc16BitRegister(&cpu.SP) }},
	&Instruction{0x34, "INC (HL)", 0, 3, func(cpu *CPU) { cpu.IncHLRegData() }},
	&Instruction{0x35, "DEC (HL)", 0, 3, func(cpu *CPU) { cpu.DecHLRegData() }},
	&Instruction{0x36, "LD (HL), d8", 1, 3, func(cpu *CPU) { cpu.Load8BitDataInto16BitRegAddr(cpu.Reg.HLByte()) }},
	&Instruction{0x37, "SCF", 0, 1, func(cpu *CPU) { cpu.SetCarryFlag() }},
	&Instruction{0x38, "JR C, s8", 1, 2, func(cpu *CPU) { cpu.JumpConditionalRelative8bit(&cpu.Reg.F.carry, true) }},
	&Instruction{0x39, "ADD HL, SP", 0, 2, func(cpu *CPU) { cpu.AddSPToHLReg() }},
	&Instruction{0x3A, "LD A, (HL-)", 0, 2, func(cpu *CPU) { cpu.Load16BitAddrDecInto8BitReg(&cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A) }},
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
	&Instruction{0xD0, "RET NC", 0, 2, func(cpu *CPU) { cpu.ReturnConditional(&cpu.Reg.F.carry, false) }},
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
	&Instruction{0xE1, "POP HL", 0, 3, func(cpu *CPU) { cpu.PopSPIntoRegPair(&cpu.Reg.H, &cpu.Reg.L) }}, // TODO - check
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
	&Instruction{0x00, "RLC B", 0, 2, func(cpu *CPU) { cpu.RotateLeftCarryReg(&cpu.Reg.B) }},
	&Instruction{0x01, "RLC C", 0, 2, func(cpu *CPU) { cpu.RotateLeftCarryReg(&cpu.Reg.C) }},
	&Instruction{0x02, "RLC D", 0, 2, func(cpu *CPU) { cpu.RotateLeftCarryReg(&cpu.Reg.D) }},
	&Instruction{0x03, "RLC E", 0, 2, func(cpu *CPU) { cpu.RotateLeftCarryReg(&cpu.Reg.E) }},
	&Instruction{0x04, "RLC H", 0, 2, func(cpu *CPU) { cpu.RotateLeftCarryReg(&cpu.Reg.H) }},
	&Instruction{0x05, "RLC L", 0, 2, func(cpu *CPU) { cpu.RotateLeftCarryReg(&cpu.Reg.L) }},
	&Instruction{0x06, "RLC (HL)", 0, 4, func(cpu *CPU) { cpu.RotateLeftCarryHLData() }},
	&Instruction{0x07, "RLC A", 0, 2, func(cpu *CPU) { cpu.RotateLeftCarryReg(&cpu.Reg.A) }},

	&Instruction{0x08, "RRC B", 0, 2, func(cpu *CPU) { cpu.RotateRightCarryReg(&cpu.Reg.B) }},
	&Instruction{0x09, "RRC C", 0, 2, func(cpu *CPU) { cpu.RotateRightCarryReg(&cpu.Reg.C) }},
	&Instruction{0x0A, "RRC D", 0, 2, func(cpu *CPU) { cpu.RotateRightCarryReg(&cpu.Reg.D) }},
	&Instruction{0x0B, "RRC E", 0, 2, func(cpu *CPU) { cpu.RotateRightCarryReg(&cpu.Reg.E) }},
	&Instruction{0x0C, "RRC H", 0, 2, func(cpu *CPU) { cpu.RotateRightCarryReg(&cpu.Reg.H) }},
	&Instruction{0x0D, "RRC L", 0, 2, func(cpu *CPU) { cpu.RotateRightCarryReg(&cpu.Reg.L) }},
	&Instruction{0x0E, "RRC (HL)", 0, 4, func(cpu *CPU) { cpu.RotateRightCarryHLData() }},
	&Instruction{0x0F, "RRC A", 0, 2, func(cpu *CPU) { cpu.RotateRightCarryReg(&cpu.Reg.A) }},

	&Instruction{0x10, "RL B", 0, 2, func(cpu *CPU) { cpu.RotateLeftReg(&cpu.Reg.B) }},
	&Instruction{0x11, "RL C", 0, 2, func(cpu *CPU) { cpu.RotateLeftReg(&cpu.Reg.C) }},
	&Instruction{0x12, "RL D", 0, 2, func(cpu *CPU) { cpu.RotateLeftReg(&cpu.Reg.D) }},
	&Instruction{0x13, "RL E", 0, 2, func(cpu *CPU) { cpu.RotateLeftReg(&cpu.Reg.E) }},
	&Instruction{0x14, "RL H", 0, 2, func(cpu *CPU) { cpu.RotateLeftReg(&cpu.Reg.H) }},
	&Instruction{0x15, "RL L", 0, 2, func(cpu *CPU) { cpu.RotateLeftReg(&cpu.Reg.L) }},
	&Instruction{0x16, "RL (HL)", 0, 4, func(cpu *CPU) { cpu.RotateLeftHLData() }},
	&Instruction{0x17, "RL A", 0, 2, func(cpu *CPU) { cpu.RotateLeftReg(&cpu.Reg.A) }},

	&Instruction{0x18, "RR B", 0, 2, func(cpu *CPU) { cpu.RotateRightReg(&cpu.Reg.B) }},
	&Instruction{0x19, "RR C", 0, 2, func(cpu *CPU) { cpu.RotateRightReg(&cpu.Reg.C) }},
	&Instruction{0x1A, "RR D", 0, 2, func(cpu *CPU) { cpu.RotateRightReg(&cpu.Reg.D) }},
	&Instruction{0x1B, "RR E", 0, 2, func(cpu *CPU) { cpu.RotateRightReg(&cpu.Reg.E) }},
	&Instruction{0x1C, "RR H", 0, 2, func(cpu *CPU) { cpu.RotateRightReg(&cpu.Reg.H) }},
	&Instruction{0x1D, "RR L", 0, 2, func(cpu *CPU) { cpu.RotateRightReg(&cpu.Reg.L) }},
	&Instruction{0x1E, "RR (HL)", 0, 4, func(cpu *CPU) { cpu.RotateRightHLData() }},
	&Instruction{0x1F, "RR A", 0, 2, func(cpu *CPU) { cpu.RotateRightReg(&cpu.Reg.A) }},

	&Instruction{0x20, "SLA B", 0, 2, func(cpu *CPU) { cpu.ShiftLeftReg(&cpu.Reg.B) }},
	&Instruction{0x21, "SLA C", 0, 2, func(cpu *CPU) { cpu.ShiftLeftReg(&cpu.Reg.C) }},
	&Instruction{0x22, "SLA D", 0, 2, func(cpu *CPU) { cpu.ShiftLeftReg(&cpu.Reg.D) }},
	&Instruction{0x23, "SLA E", 0, 2, func(cpu *CPU) { cpu.ShiftLeftReg(&cpu.Reg.E) }},
	&Instruction{0x24, "SLA H", 0, 2, func(cpu *CPU) { cpu.ShiftLeftReg(&cpu.Reg.H) }},
	&Instruction{0x25, "SLA L", 0, 2, func(cpu *CPU) { cpu.ShiftLeftReg(&cpu.Reg.L) }},
	&Instruction{0x26, "SLA (HL)", 0, 4, func(cpu *CPU) { cpu.ShiftLeftHLData() }},
	&Instruction{0x27, "SLA A", 0, 2, func(cpu *CPU) { cpu.ShiftLeftReg(&cpu.Reg.A) }},

	&Instruction{0x28, "SRA B", 0, 2, func(cpu *CPU) { cpu.ShiftRightReg(&cpu.Reg.B) }},
	&Instruction{0x29, "SRA C", 0, 2, func(cpu *CPU) { cpu.ShiftRightReg(&cpu.Reg.C) }},
	&Instruction{0x2A, "SRA D", 0, 2, func(cpu *CPU) { cpu.ShiftRightReg(&cpu.Reg.D) }},
	&Instruction{0x2B, "SRA E", 0, 2, func(cpu *CPU) { cpu.ShiftRightReg(&cpu.Reg.E) }},
	&Instruction{0x2C, "SRA H", 0, 2, func(cpu *CPU) { cpu.ShiftRightReg(&cpu.Reg.H) }},
	&Instruction{0x2D, "SRA L", 0, 2, func(cpu *CPU) { cpu.ShiftRightReg(&cpu.Reg.L) }},
	&Instruction{0x2E, "SRA (HL)", 0, 4, func(cpu *CPU) { cpu.ShiftRightHLData() }},
	&Instruction{0x2F, "SRA A", 0, 2, func(cpu *CPU) { cpu.ShiftRightReg(&cpu.Reg.A) }},

	&Instruction{0x30, "SWAP B", 0, 2, func(cpu *CPU) { cpu.SwapReg(&cpu.Reg.B) }},
	&Instruction{0x31, "SWAP C", 0, 2, func(cpu *CPU) { cpu.SwapReg(&cpu.Reg.C) }},
	&Instruction{0x32, "SWAP D", 0, 2, func(cpu *CPU) { cpu.SwapReg(&cpu.Reg.D) }},
	&Instruction{0x33, "SWAP E", 0, 2, func(cpu *CPU) { cpu.SwapReg(&cpu.Reg.E) }},
	&Instruction{0x34, "SWAP H", 0, 2, func(cpu *CPU) { cpu.SwapReg(&cpu.Reg.H) }},
	&Instruction{0x35, "SWAP L", 0, 2, func(cpu *CPU) { cpu.SwapReg(&cpu.Reg.L) }},
	&Instruction{0x36, "SWAP (HL)", 0, 4, func(cpu *CPU) { cpu.SwapHLData() }},
	&Instruction{0x37, "SWAP A", 0, 2, func(cpu *CPU) { cpu.SwapReg(&cpu.Reg.A) }},

	&Instruction{0x38, "SRL B", 0, 2, func(cpu *CPU) { cpu.ShiftRightRegLogical(&cpu.Reg.B) }},
	&Instruction{0x39, "SRL C", 0, 2, func(cpu *CPU) { cpu.ShiftRightRegLogical(&cpu.Reg.C) }},
	&Instruction{0x3A, "SRL D", 0, 2, func(cpu *CPU) { cpu.ShiftRightRegLogical(&cpu.Reg.D) }},
	&Instruction{0x3B, "SRL E", 0, 2, func(cpu *CPU) { cpu.ShiftRightRegLogical(&cpu.Reg.E) }},
	&Instruction{0x3C, "SRL H", 0, 2, func(cpu *CPU) { cpu.ShiftRightRegLogical(&cpu.Reg.H) }},
	&Instruction{0x3D, "SRL L", 0, 2, func(cpu *CPU) { cpu.ShiftRightRegLogical(&cpu.Reg.L) }},
	&Instruction{0x3E, "SRL (HL)", 0, 4, func(cpu *CPU) { cpu.ShiftRightHLDataLogical() }},
	&Instruction{0x3F, "SRL A", 0, 2, func(cpu *CPU) { cpu.ShiftRightRegLogical(&cpu.Reg.A) }},

	&Instruction{0x40, "BIT 0, B", 0, 2, func(cpu *CPU) { cpu.BitTestReg(0, &cpu.Reg.B) }},
	&Instruction{0x41, "BIT 0, C", 0, 2, func(cpu *CPU) { cpu.BitTestReg(0, &cpu.Reg.C) }},
	&Instruction{0x42, "BIT 0, D", 0, 2, func(cpu *CPU) { cpu.BitTestReg(0, &cpu.Reg.D) }},
	&Instruction{0x43, "BIT 0, E", 0, 2, func(cpu *CPU) { cpu.BitTestReg(0, &cpu.Reg.E) }},
	&Instruction{0x44, "BIT 0, H", 0, 2, func(cpu *CPU) { cpu.BitTestReg(0, &cpu.Reg.H) }},
	&Instruction{0x45, "BIT 0, L", 0, 2, func(cpu *CPU) { cpu.BitTestReg(0, &cpu.Reg.L) }},
	&Instruction{0x46, "BIT 0, (HL)", 0, 3, func(cpu *CPU) { cpu.BitTestHLData(0) }},
	&Instruction{0x47, "BIT 0, A", 0, 2, func(cpu *CPU) { cpu.BitTestReg(0, &cpu.Reg.A) }},

	&Instruction{0x48, "BIT 1, B", 0, 2, func(cpu *CPU) { cpu.BitTestReg(1, &cpu.Reg.B) }},
	&Instruction{0x49, "BIT 1, C", 0, 2, func(cpu *CPU) { cpu.BitTestReg(1, &cpu.Reg.C) }},
	&Instruction{0x4A, "BIT 1, D", 0, 2, func(cpu *CPU) { cpu.BitTestReg(1, &cpu.Reg.D) }},
	&Instruction{0x4B, "BIT 1, E", 0, 2, func(cpu *CPU) { cpu.BitTestReg(1, &cpu.Reg.E) }},
	&Instruction{0x4C, "BIT 1, H", 0, 2, func(cpu *CPU) { cpu.BitTestReg(1, &cpu.Reg.H) }},
	&Instruction{0x4D, "BIT 1, L", 0, 2, func(cpu *CPU) { cpu.BitTestReg(1, &cpu.Reg.L) }},
	&Instruction{0x4E, "BIT 1, (HL)", 0, 3, func(cpu *CPU) { cpu.BitTestHLData(1) }},
	&Instruction{0x4F, "BIT 1, A", 0, 2, func(cpu *CPU) { cpu.BitTestReg(1, &cpu.Reg.A) }},

	&Instruction{0x50, "BIT 2, B", 0, 2, func(cpu *CPU) { cpu.BitTestReg(2, &cpu.Reg.B) }},
	&Instruction{0x51, "BIT 2, C", 0, 2, func(cpu *CPU) { cpu.BitTestReg(2, &cpu.Reg.C) }},
	&Instruction{0x52, "BIT 2, D", 0, 2, func(cpu *CPU) { cpu.BitTestReg(2, &cpu.Reg.D) }},
	&Instruction{0x53, "BIT 2, E", 0, 2, func(cpu *CPU) { cpu.BitTestReg(2, &cpu.Reg.E) }},
	&Instruction{0x54, "BIT 2, H", 0, 2, func(cpu *CPU) { cpu.BitTestReg(2, &cpu.Reg.H) }},
	&Instruction{0x55, "BIT 2, L", 0, 2, func(cpu *CPU) { cpu.BitTestReg(2, &cpu.Reg.L) }},
	&Instruction{0x56, "BIT 2, (HL)", 0, 3, func(cpu *CPU) { cpu.BitTestHLData(2) }},
	&Instruction{0x57, "BIT 2, A", 0, 2, func(cpu *CPU) { cpu.BitTestReg(2, &cpu.Reg.A) }},

	&Instruction{0x58, "BIT 3, B", 0, 2, func(cpu *CPU) { cpu.BitTestReg(3, &cpu.Reg.B) }},
	&Instruction{0x59, "BIT 3, C", 0, 2, func(cpu *CPU) { cpu.BitTestReg(3, &cpu.Reg.C) }},
	&Instruction{0x5A, "BIT 3, D", 0, 2, func(cpu *CPU) { cpu.BitTestReg(3, &cpu.Reg.D) }},
	&Instruction{0x5B, "BIT 3, E", 0, 2, func(cpu *CPU) { cpu.BitTestReg(3, &cpu.Reg.E) }},
	&Instruction{0x5C, "BIT 3, H", 0, 2, func(cpu *CPU) { cpu.BitTestReg(3, &cpu.Reg.H) }},
	&Instruction{0x5D, "BIT 3, L", 0, 2, func(cpu *CPU) { cpu.BitTestReg(3, &cpu.Reg.L) }},
	&Instruction{0x5E, "BIT 3, (HL)", 0, 3, func(cpu *CPU) { cpu.BitTestHLData(3) }},
	&Instruction{0x5F, "BIT 3, A", 0, 2, func(cpu *CPU) { cpu.BitTestReg(3, &cpu.Reg.A) }},

	&Instruction{0x60, "BIT 4, B", 0, 2, func(cpu *CPU) { cpu.BitTestReg(4, &cpu.Reg.B) }},
	&Instruction{0x61, "BIT 4, C", 0, 2, func(cpu *CPU) { cpu.BitTestReg(4, &cpu.Reg.C) }},
	&Instruction{0x62, "BIT 4, D", 0, 2, func(cpu *CPU) { cpu.BitTestReg(4, &cpu.Reg.D) }},
	&Instruction{0x63, "BIT 4, E", 0, 2, func(cpu *CPU) { cpu.BitTestReg(4, &cpu.Reg.E) }},
	&Instruction{0x64, "BIT 4, H", 0, 2, func(cpu *CPU) { cpu.BitTestReg(4, &cpu.Reg.H) }},
	&Instruction{0x65, "BIT 4, L", 0, 2, func(cpu *CPU) { cpu.BitTestReg(4, &cpu.Reg.L) }},
	&Instruction{0x66, "BIT 4, (HL)", 0, 3, func(cpu *CPU) { cpu.BitTestHLData(4) }},
	&Instruction{0x67, "BIT 4, A", 0, 2, func(cpu *CPU) { cpu.BitTestReg(4, &cpu.Reg.A) }},

	&Instruction{0x68, "BIT 5, B", 0, 2, func(cpu *CPU) { cpu.BitTestReg(5, &cpu.Reg.B) }},
	&Instruction{0x69, "BIT 5, C", 0, 2, func(cpu *CPU) { cpu.BitTestReg(5, &cpu.Reg.C) }},
	&Instruction{0x6A, "BIT 5, D", 0, 2, func(cpu *CPU) { cpu.BitTestReg(5, &cpu.Reg.D) }},
	&Instruction{0x6B, "BIT 5, E", 0, 2, func(cpu *CPU) { cpu.BitTestReg(5, &cpu.Reg.E) }},
	&Instruction{0x6C, "BIT 5, H", 0, 2, func(cpu *CPU) { cpu.BitTestReg(5, &cpu.Reg.H) }},
	&Instruction{0x6D, "BIT 5, L", 0, 2, func(cpu *CPU) { cpu.BitTestReg(5, &cpu.Reg.L) }},
	&Instruction{0x6E, "BIT 5, (HL)", 0, 3, func(cpu *CPU) { cpu.BitTestHLData(5) }},
	&Instruction{0x6F, "BIT 5, A", 0, 2, func(cpu *CPU) { cpu.BitTestReg(5, &cpu.Reg.A) }},

	&Instruction{0x70, "BIT 6, B", 0, 2, func(cpu *CPU) { cpu.BitTestReg(6, &cpu.Reg.B) }},
	&Instruction{0x71, "BIT 6, C", 0, 2, func(cpu *CPU) { cpu.BitTestReg(6, &cpu.Reg.C) }},
	&Instruction{0x72, "BIT 6, D", 0, 2, func(cpu *CPU) { cpu.BitTestReg(6, &cpu.Reg.D) }},
	&Instruction{0x73, "BIT 6, E", 0, 2, func(cpu *CPU) { cpu.BitTestReg(6, &cpu.Reg.E) }},
	&Instruction{0x74, "BIT 6, H", 0, 2, func(cpu *CPU) { cpu.BitTestReg(6, &cpu.Reg.H) }},
	&Instruction{0x75, "BIT 6, L", 0, 2, func(cpu *CPU) { cpu.BitTestReg(6, &cpu.Reg.L) }},
	&Instruction{0x76, "BIT 6, (HL)", 0, 3, func(cpu *CPU) { cpu.BitTestHLData(6) }},
	&Instruction{0x77, "BIT 6, A", 0, 2, func(cpu *CPU) { cpu.BitTestReg(6, &cpu.Reg.A) }},

	&Instruction{0x78, "BIT 7, B", 0, 2, func(cpu *CPU) { cpu.BitTestReg(7, &cpu.Reg.B) }},
	&Instruction{0x79, "BIT 7, C", 0, 2, func(cpu *CPU) { cpu.BitTestReg(7, &cpu.Reg.C) }},
	&Instruction{0x7A, "BIT 7, D", 0, 2, func(cpu *CPU) { cpu.BitTestReg(7, &cpu.Reg.D) }},
	&Instruction{0x7B, "BIT 7, E", 0, 2, func(cpu *CPU) { cpu.BitTestReg(7, &cpu.Reg.E) }},
	&Instruction{0x7C, "BIT 7, H", 0, 2, func(cpu *CPU) { cpu.BitTestReg(7, &cpu.Reg.H) }},
	&Instruction{0x7D, "BIT 7, L", 0, 2, func(cpu *CPU) { cpu.BitTestReg(7, &cpu.Reg.L) }},
	&Instruction{0x7E, "BIT 7, (HL)", 0, 3, func(cpu *CPU) { cpu.BitTestHLData(7) }},
	&Instruction{0x7F, "BIT 7, A", 0, 2, func(cpu *CPU) { cpu.BitTestReg(7, &cpu.Reg.A) }},

	&Instruction{0x80, "RES 0, B", 0, 2, func(cpu *CPU) { cpu.BitResetReg(0, &cpu.Reg.B) }},
	&Instruction{0x81, "RES 0, C", 0, 2, func(cpu *CPU) { cpu.BitResetReg(0, &cpu.Reg.C) }},
	&Instruction{0x82, "RES 0, D", 0, 2, func(cpu *CPU) { cpu.BitResetReg(0, &cpu.Reg.D) }},
	&Instruction{0x83, "RES 0, E", 0, 2, func(cpu *CPU) { cpu.BitResetReg(0, &cpu.Reg.E) }},
	&Instruction{0x84, "RES 0, H", 0, 2, func(cpu *CPU) { cpu.BitResetReg(0, &cpu.Reg.H) }},
	&Instruction{0x85, "RES 0, L", 0, 2, func(cpu *CPU) { cpu.BitResetReg(0, &cpu.Reg.L) }},
	&Instruction{0x86, "RES 0, (HL)", 0, 4, func(cpu *CPU) { cpu.BitResetHLData(0) }},
	&Instruction{0x87, "RES 0, A", 0, 2, func(cpu *CPU) { cpu.BitResetReg(0, &cpu.Reg.A) }},

	&Instruction{0x88, "RES 1, B", 0, 2, func(cpu *CPU) { cpu.BitResetReg(1, &cpu.Reg.B) }},
	&Instruction{0x89, "RES 1, C", 0, 2, func(cpu *CPU) { cpu.BitResetReg(1, &cpu.Reg.C) }},
	&Instruction{0x8A, "RES 1, D", 0, 2, func(cpu *CPU) { cpu.BitResetReg(1, &cpu.Reg.D) }},
	&Instruction{0x8B, "RES 1, E", 0, 2, func(cpu *CPU) { cpu.BitResetReg(1, &cpu.Reg.E) }},
	&Instruction{0x8C, "RES 1, H", 0, 2, func(cpu *CPU) { cpu.BitResetReg(1, &cpu.Reg.H) }},
	&Instruction{0x8D, "RES 1, L", 0, 2, func(cpu *CPU) { cpu.BitResetReg(1, &cpu.Reg.L) }},
	&Instruction{0x8E, "RES 1, (HL)", 0, 4, func(cpu *CPU) { cpu.BitResetHLData(1) }},
	&Instruction{0x8F, "RES 1, A", 0, 2, func(cpu *CPU) { cpu.BitResetReg(1, &cpu.Reg.A) }},

	&Instruction{0x90, "RES 2, B", 0, 2, func(cpu *CPU) { cpu.BitResetReg(2, &cpu.Reg.B) }},
	&Instruction{0x91, "RES 2, C", 0, 2, func(cpu *CPU) { cpu.BitResetReg(2, &cpu.Reg.C) }},
	&Instruction{0x92, "RES 2, D", 0, 2, func(cpu *CPU) { cpu.BitResetReg(2, &cpu.Reg.D) }},
	&Instruction{0x93, "RES 2, E", 0, 2, func(cpu *CPU) { cpu.BitResetReg(2, &cpu.Reg.E) }},
	&Instruction{0x94, "RES 2, H", 0, 2, func(cpu *CPU) { cpu.BitResetReg(2, &cpu.Reg.H) }},
	&Instruction{0x95, "RES 2, L", 0, 2, func(cpu *CPU) { cpu.BitResetReg(2, &cpu.Reg.L) }},
	&Instruction{0x96, "RES 2, (HL)", 0, 4, func(cpu *CPU) { cpu.BitResetHLData(2) }},
	&Instruction{0x97, "RES 2, A", 0, 2, func(cpu *CPU) { cpu.BitResetReg(2, &cpu.Reg.A) }},

	&Instruction{0x98, "RES 3, B", 0, 2, func(cpu *CPU) { cpu.BitResetReg(3, &cpu.Reg.B) }},
	&Instruction{0x99, "RES 3, C", 0, 2, func(cpu *CPU) { cpu.BitResetReg(3, &cpu.Reg.C) }},
	&Instruction{0x9A, "RES 3, D", 0, 2, func(cpu *CPU) { cpu.BitResetReg(3, &cpu.Reg.D) }},
	&Instruction{0x9B, "RES 3, E", 0, 2, func(cpu *CPU) { cpu.BitResetReg(3, &cpu.Reg.E) }},
	&Instruction{0x9C, "RES 3, H", 0, 2, func(cpu *CPU) { cpu.BitResetReg(3, &cpu.Reg.H) }},
	&Instruction{0x9D, "RES 3, L", 0, 2, func(cpu *CPU) { cpu.BitResetReg(3, &cpu.Reg.L) }},
	&Instruction{0x9E, "RES 3, (HL)", 0, 4, func(cpu *CPU) { cpu.BitResetHLData(3) }},
	&Instruction{0x9F, "RES 3, A", 0, 2, func(cpu *CPU) { cpu.BitResetReg(3, &cpu.Reg.A) }},

	&Instruction{0xA0, "RES 4, B", 0, 2, func(cpu *CPU) { cpu.BitResetReg(4, &cpu.Reg.B) }},
	&Instruction{0xA1, "RES 4, C", 0, 2, func(cpu *CPU) { cpu.BitResetReg(4, &cpu.Reg.C) }},
	&Instruction{0xA2, "RES 4, D", 0, 2, func(cpu *CPU) { cpu.BitResetReg(4, &cpu.Reg.D) }},
	&Instruction{0xA3, "RES 4, E", 0, 2, func(cpu *CPU) { cpu.BitResetReg(4, &cpu.Reg.E) }},
	&Instruction{0xA4, "RES 4, H", 0, 2, func(cpu *CPU) { cpu.BitResetReg(4, &cpu.Reg.H) }},
	&Instruction{0xA5, "RES 4, L", 0, 2, func(cpu *CPU) { cpu.BitResetReg(4, &cpu.Reg.L) }},
	&Instruction{0xA6, "RES 4, (HL)", 0, 4, func(cpu *CPU) { cpu.BitResetHLData(4) }},
	&Instruction{0xA7, "RES 4, A", 0, 2, func(cpu *CPU) { cpu.BitResetReg(4, &cpu.Reg.A) }},

	&Instruction{0xA8, "RES 5, B", 0, 2, func(cpu *CPU) { cpu.BitResetReg(5, &cpu.Reg.B) }},
	&Instruction{0xA9, "RES 5, C", 0, 2, func(cpu *CPU) { cpu.BitResetReg(5, &cpu.Reg.C) }},
	&Instruction{0xAA, "RES 5, D", 0, 2, func(cpu *CPU) { cpu.BitResetReg(5, &cpu.Reg.D) }},
	&Instruction{0xAB, "RES 5, E", 0, 2, func(cpu *CPU) { cpu.BitResetReg(5, &cpu.Reg.E) }},
	&Instruction{0xAC, "RES 5, H", 0, 2, func(cpu *CPU) { cpu.BitResetReg(5, &cpu.Reg.H) }},
	&Instruction{0xAD, "RES 5, L", 0, 2, func(cpu *CPU) { cpu.BitResetReg(5, &cpu.Reg.L) }},
	&Instruction{0xAE, "RES 5, (HL)", 0, 4, func(cpu *CPU) { cpu.BitResetHLData(5) }},
	&Instruction{0xAF, "RES 5, A", 0, 2, func(cpu *CPU) { cpu.BitResetReg(5, &cpu.Reg.A) }},

	&Instruction{0xB0, "RES 6, B", 0, 2, func(cpu *CPU) { cpu.BitResetReg(6, &cpu.Reg.B) }},
	&Instruction{0xB1, "RES 6, C", 0, 2, func(cpu *CPU) { cpu.BitResetReg(6, &cpu.Reg.C) }},
	&Instruction{0xB2, "RES 6, D", 0, 2, func(cpu *CPU) { cpu.BitResetReg(6, &cpu.Reg.D) }},
	&Instruction{0xB3, "RES 6, E", 0, 2, func(cpu *CPU) { cpu.BitResetReg(6, &cpu.Reg.E) }},
	&Instruction{0xB4, "RES 6, H", 0, 2, func(cpu *CPU) { cpu.BitResetReg(6, &cpu.Reg.H) }},
	&Instruction{0xB5, "RES 6, L", 0, 2, func(cpu *CPU) { cpu.BitResetReg(6, &cpu.Reg.L) }},
	&Instruction{0xB6, "RES 6, (HL)", 0, 4, func(cpu *CPU) { cpu.BitResetHLData(6) }},
	&Instruction{0xB7, "RES 6, A", 0, 2, func(cpu *CPU) { cpu.BitResetReg(6, &cpu.Reg.A) }},

	&Instruction{0xB8, "RES 7, B", 0, 2, func(cpu *CPU) { cpu.BitResetReg(7, &cpu.Reg.B) }},
	&Instruction{0xB9, "RES 7, C", 0, 2, func(cpu *CPU) { cpu.BitResetReg(7, &cpu.Reg.C) }},
	&Instruction{0xBA, "RES 7, D", 0, 2, func(cpu *CPU) { cpu.BitResetReg(7, &cpu.Reg.D) }},
	&Instruction{0xBB, "RES 7, E", 0, 2, func(cpu *CPU) { cpu.BitResetReg(7, &cpu.Reg.E) }},
	&Instruction{0xBC, "RES 7, H", 0, 2, func(cpu *CPU) { cpu.BitResetReg(7, &cpu.Reg.H) }},
	&Instruction{0xBD, "RES 7, L", 0, 2, func(cpu *CPU) { cpu.BitResetReg(7, &cpu.Reg.L) }},
	&Instruction{0xBE, "RES 7, (HL)", 0, 4, func(cpu *CPU) { cpu.BitResetHLData(7) }},
	&Instruction{0xBF, "RES 7, A", 0, 2, func(cpu *CPU) { cpu.BitResetReg(7, &cpu.Reg.A) }},

	&Instruction{0xC0, "SET 0, B", 0, 2, func(cpu *CPU) { cpu.BitSetReg(0, &cpu.Reg.B) }},
	&Instruction{0xC1, "SET 0, C", 0, 2, func(cpu *CPU) { cpu.BitSetReg(0, &cpu.Reg.C) }},
	&Instruction{0xC2, "SET 0, D", 0, 2, func(cpu *CPU) { cpu.BitSetReg(0, &cpu.Reg.D) }},
	&Instruction{0xC3, "SET 0, E", 0, 2, func(cpu *CPU) { cpu.BitSetReg(0, &cpu.Reg.E) }},
	&Instruction{0xC4, "SET 0, H", 0, 2, func(cpu *CPU) { cpu.BitSetReg(0, &cpu.Reg.H) }},
	&Instruction{0xC5, "SET 0, L", 0, 2, func(cpu *CPU) { cpu.BitSetReg(0, &cpu.Reg.L) }},
	&Instruction{0xC6, "SET 0, (HL)", 0, 4, func(cpu *CPU) { cpu.BitSetHLData(0) }},
	&Instruction{0xC7, "SET 0, A", 0, 2, func(cpu *CPU) { cpu.BitSetReg(0, &cpu.Reg.A) }},

	&Instruction{0xC8, "SET 1, B", 0, 2, func(cpu *CPU) { cpu.BitSetReg(1, &cpu.Reg.B) }},
	&Instruction{0xC9, "SET 1, C", 0, 2, func(cpu *CPU) { cpu.BitSetReg(1, &cpu.Reg.C) }},
	&Instruction{0xCA, "SET 1, D", 0, 2, func(cpu *CPU) { cpu.BitSetReg(1, &cpu.Reg.D) }},
	&Instruction{0xCB, "SET 1, E", 0, 2, func(cpu *CPU) { cpu.BitSetReg(1, &cpu.Reg.E) }},
	&Instruction{0xCC, "SET 1, H", 0, 2, func(cpu *CPU) { cpu.BitSetReg(1, &cpu.Reg.H) }},
	&Instruction{0xCD, "SET 1, L", 0, 2, func(cpu *CPU) { cpu.BitSetReg(1, &cpu.Reg.L) }},
	&Instruction{0xCE, "SET 1, (HL)", 0, 4, func(cpu *CPU) { cpu.BitSetHLData(1) }},
	&Instruction{0xCF, "SET 1, A", 0, 2, func(cpu *CPU) { cpu.BitSetReg(1, &cpu.Reg.A) }},

	&Instruction{0xD0, "SET 2, B", 0, 2, func(cpu *CPU) { cpu.BitSetReg(2, &cpu.Reg.B) }},
	&Instruction{0xD1, "SET 2, C", 0, 2, func(cpu *CPU) { cpu.BitSetReg(2, &cpu.Reg.C) }},
	&Instruction{0xD2, "SET 2, D", 0, 2, func(cpu *CPU) { cpu.BitSetReg(2, &cpu.Reg.D) }},
	&Instruction{0xD3, "SET 2, E", 0, 2, func(cpu *CPU) { cpu.BitSetReg(2, &cpu.Reg.E) }},
	&Instruction{0xD4, "SET 2, H", 0, 2, func(cpu *CPU) { cpu.BitSetReg(2, &cpu.Reg.H) }},
	&Instruction{0xD5, "SET 2, L", 0, 2, func(cpu *CPU) { cpu.BitSetReg(2, &cpu.Reg.L) }},
	&Instruction{0xD6, "SET 2, (HL)", 0, 4, func(cpu *CPU) { cpu.BitSetHLData(2) }},
	&Instruction{0xD7, "SET 2, A", 0, 2, func(cpu *CPU) { cpu.BitSetReg(2, &cpu.Reg.A) }},

	&Instruction{0xD8, "SET 3, B", 0, 2, func(cpu *CPU) { cpu.BitSetReg(3, &cpu.Reg.B) }},
	&Instruction{0xD9, "SET 3, C", 0, 2, func(cpu *CPU) { cpu.BitSetReg(3, &cpu.Reg.C) }},
	&Instruction{0xDA, "SET 3, D", 0, 2, func(cpu *CPU) { cpu.BitSetReg(3, &cpu.Reg.D) }},
	&Instruction{0xDB, "SET 3, E", 0, 2, func(cpu *CPU) { cpu.BitSetReg(3, &cpu.Reg.E) }},
	&Instruction{0xDC, "SET 3, H", 0, 2, func(cpu *CPU) { cpu.BitSetReg(3, &cpu.Reg.H) }},
	&Instruction{0xDD, "SET 3, L", 0, 2, func(cpu *CPU) { cpu.BitSetReg(3, &cpu.Reg.L) }},
	&Instruction{0xDE, "SET 3, (HL)", 0, 4, func(cpu *CPU) { cpu.BitSetHLData(3) }},
	&Instruction{0xDF, "SET 3, A", 0, 2, func(cpu *CPU) { cpu.BitSetReg(3, &cpu.Reg.A) }},

	&Instruction{0xE0, "SET 4, B", 0, 2, func(cpu *CPU) { cpu.BitSetReg(4, &cpu.Reg.B) }},
	&Instruction{0xE1, "SET 4, C", 0, 2, func(cpu *CPU) { cpu.BitSetReg(4, &cpu.Reg.C) }},
	&Instruction{0xE2, "SET 4, D", 0, 2, func(cpu *CPU) { cpu.BitSetReg(4, &cpu.Reg.D) }},
	&Instruction{0xE3, "SET 4, E", 0, 2, func(cpu *CPU) { cpu.BitSetReg(4, &cpu.Reg.E) }},
	&Instruction{0xE4, "SET 4, H", 0, 2, func(cpu *CPU) { cpu.BitSetReg(4, &cpu.Reg.H) }},
	&Instruction{0xE5, "SET 4, L", 0, 2, func(cpu *CPU) { cpu.BitSetReg(4, &cpu.Reg.L) }},
	&Instruction{0xE6, "SET 4, (HL)", 0, 4, func(cpu *CPU) { cpu.BitSetHLData(4) }},
	&Instruction{0xE7, "SET 4, A", 0, 2, func(cpu *CPU) { cpu.BitSetReg(4, &cpu.Reg.A) }},

	&Instruction{0xE8, "SET 5, B", 0, 2, func(cpu *CPU) { cpu.BitSetReg(5, &cpu.Reg.B) }},
	&Instruction{0xE9, "SET 5, C", 0, 2, func(cpu *CPU) { cpu.BitSetReg(5, &cpu.Reg.C) }},
	&Instruction{0xEA, "SET 5, D", 0, 2, func(cpu *CPU) { cpu.BitSetReg(5, &cpu.Reg.D) }},
	&Instruction{0xEB, "SET 5, E", 0, 2, func(cpu *CPU) { cpu.BitSetReg(5, &cpu.Reg.E) }},
	&Instruction{0xEC, "SET 5, H", 0, 2, func(cpu *CPU) { cpu.BitSetReg(5, &cpu.Reg.H) }},
	&Instruction{0xED, "SET 5, L", 0, 2, func(cpu *CPU) { cpu.BitSetReg(5, &cpu.Reg.L) }},
	&Instruction{0xEE, "SET 5, (HL)", 0, 4, func(cpu *CPU) { cpu.BitSetHLData(5) }},
	&Instruction{0xEF, "SET 5, A", 0, 2, func(cpu *CPU) { cpu.BitSetReg(5, &cpu.Reg.A) }},

	&Instruction{0xF0, "SET 6, B", 0, 2, func(cpu *CPU) { cpu.BitSetReg(6, &cpu.Reg.B) }},
	&Instruction{0xF1, "SET 6, C", 0, 2, func(cpu *CPU) { cpu.BitSetReg(6, &cpu.Reg.C) }},
	&Instruction{0xF2, "SET 6, D", 0, 2, func(cpu *CPU) { cpu.BitSetReg(6, &cpu.Reg.D) }},
	&Instruction{0xF3, "SET 6, E", 0, 2, func(cpu *CPU) { cpu.BitSetReg(6, &cpu.Reg.E) }},
	&Instruction{0xF4, "SET 6, H", 0, 2, func(cpu *CPU) { cpu.BitSetReg(6, &cpu.Reg.H) }},
	&Instruction{0xF5, "SET 6, L", 0, 2, func(cpu *CPU) { cpu.BitSetReg(6, &cpu.Reg.L) }},
	&Instruction{0xF6, "SET 6, (HL)", 0, 4, func(cpu *CPU) { cpu.BitSetHLData(6) }},
	&Instruction{0xF7, "SET 6, A", 0, 2, func(cpu *CPU) { cpu.BitSetReg(6, &cpu.Reg.A) }},

	&Instruction{0xF8, "SET 7, B", 0, 2, func(cpu *CPU) { cpu.BitSetReg(7, &cpu.Reg.B) }},
	&Instruction{0xF9, "SET 7, C", 0, 2, func(cpu *CPU) { cpu.BitSetReg(7, &cpu.Reg.C) }},
	&Instruction{0xFA, "SET 7, D", 0, 2, func(cpu *CPU) { cpu.BitSetReg(7, &cpu.Reg.D) }},
	&Instruction{0xFB, "SET 7, E", 0, 2, func(cpu *CPU) { cpu.BitSetReg(7, &cpu.Reg.E) }},
	&Instruction{0xFC, "SET 7, H", 0, 2, func(cpu *CPU) { cpu.BitSetReg(7, &cpu.Reg.H) }},
	&Instruction{0xFD, "SET 7, L", 0, 2, func(cpu *CPU) { cpu.BitSetReg(7, &cpu.Reg.L) }},
	&Instruction{0xFE, "SET 7, (HL)", 0, 4, func(cpu *CPU) { cpu.BitSetHLData(7) }},
	&Instruction{0xFF, "SET 7, A", 0, 2, func(cpu *CPU) { cpu.BitSetReg(7, &cpu.Reg.A) }},
}
