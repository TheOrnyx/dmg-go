package main


// Constants for the register flags positions
const (
	ZeroFlagBytePos uint8 = 7
	SubtractFlagBytePos uint8 = 6
	HalfCarryFlagBytePos uint8 = 5
	CarryFlagBytePos uint8 = 4
)

// FlagsRegister struct for flags from the register
type FlagsRegister struct {
	zero bool
	subtract bool
	half_carry bool
	carry bool
}

// toByte convert the FlagsRegister f to a uint8 byte value
func (f *FlagsRegister) toByte() uint8 {
	return uint8(
		(boolToBit(f.zero) << ZeroFlagBytePos) |
			(boolToBit(f.subtract) << SubtractFlagBytePos) |
			(boolToBit(f.half_carry) << HalfCarryFlagBytePos) |
			(boolToBit(f.carry) << CarryFlagBytePos),
	)
}

// byteToFlagsRegister convert a uint8 byte value to FlagsRegister instance
func byteToFlagsRegister(b uint8) FlagsRegister {
	return FlagsRegister{
		zero: ((b >> ZeroFlagBytePos) & 0b1) != 0,
		subtract:   ((b >> SubtractFlagBytePos) & 0b1) != 0,
		half_carry: ((b >> HalfCarryFlagBytePos) & 0b1) != 0,
		carry:      ((b >> CarryFlagBytePos) & 0b1) != 0,
	}
}


// boolToBit convert a bool to a bit value (1 for true, 0 for false)
func boolToBit(b bool) uint8 {
	if b {
		return 1
	}

	return 0
}

// Registers the memory Registers for the CPU
type Registers struct {
	a uint8
	b uint8
	c uint8
	d uint8
	e uint8
	f uint8
	h uint8
	l uint8
}

// Combined register methods
// These methods are used to get the combination of two registers

// GetBC get the combination register of b and c
func (r *Registers) GetBC() uint16 {
	return (uint16(r.b) << 8) | uint16(r.c)
}

// SetBC set the value stored in the b and c combination register
func (r *Registers) SetBC(value uint16)  {
	r.b = uint8((value & 0xFF00) >> 8)
	r.c = uint8(value & 0xFF)
}
