package cpu

import "fmt"

// Constants for the register flags positions
const (
	ZeroFlagBytePos      uint8 = 7
	SubtractFlagBytePos  uint8 = 6
	HalfCarryFlagBytePos uint8 = 5
	CarryFlagBytePos     uint8 = 4
)

const (
	Z = iota // zero flag
	N        // subtraction flag
	H        // half carry flag
	C        // carry flag
)

// FlagsRegister struct for flags from the register
type FlagsRegister struct {
	zero       bool
	subtract   bool
	half_carry bool
	carry      bool
}

// toByte convert the FlagsRegister f to a uint8 byte value
func (f *FlagsRegister) toByte() byte {
	return byte(
		(boolToBit(f.zero) << ZeroFlagBytePos) |
			(boolToBit(f.subtract) << SubtractFlagBytePos) |
			(boolToBit(f.half_carry) << HalfCarryFlagBytePos) |
			(boolToBit(f.carry) << CarryFlagBytePos),
	)
}

// byteToFlagsRegister convert a uint8 byte value to FlagsRegister instance
func byteToFlagsRegister(b uint8) FlagsRegister {
	return FlagsRegister{
		zero:       ((b >> ZeroFlagBytePos) & 0b1) != 0,
		subtract:   ((b >> SubtractFlagBytePos) & 0b1) != 0,
		half_carry: ((b >> HalfCarryFlagBytePos) & 0b1) != 0,
		carry:      ((b >> CarryFlagBytePos) & 0b1) != 0,
	}
}



// Registers the memory Registers for the CPU
type Registers struct {
	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	F FlagsRegister
	H uint8
	L uint8
}

// String register string
func (r *Registers) String() string {
	return fmt.Sprintf("A:%v F:%v B:%v C:%v D:%v E:%v H:%v L:%v", r.A, r.F.toByte(), r.B, r.C, r.D, r.E, r.H, r.L)
}

// HL return the address of the HL register pair
func (r *Registers) HL() uint16 {
	return JoinBytes(r.H, r.L)
}

// HLByte returns the HL register pair as two byte pointers
// TODO - maybe switch my instruction calls to use this instead of individual stuff but idk
func (r *Registers) HLByte() (*byte, *byte) {
	return &r.H, &r.L
}

// reset reset the Registers to their default state
func (r *Registers) reset()  {
	r.A = 0
	r.B = 0
	r.C = 0
	r.D = 0
	r.E = 0
	r.H = 0
	r.A = 0

	// TODO maybe move the flag to be here
}

// Combined register methods
// These methods are used to get the combination of two registers

// GetBC get the combination register of b and c
func (r *Registers) GetBC() uint16 {
	return (uint16(r.B) << 8) | uint16(r.C)
}

// SetBC set the value stored in the b and c combination register
func (r *Registers) SetBC(value uint16) {
	r.B = uint8((value & 0xFF00) >> 8)
	r.C = uint8(value & 0xFF)
}
