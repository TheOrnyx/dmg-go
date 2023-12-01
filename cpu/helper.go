package cpu

type unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

// overflowAdd takes in two unsigned values, adds them and checks if the addition has resulted in an overflow
// Returns the result of the adition and a bool type representing whether the value overflowed
func overflowAdd[T unsigned](a, b T) (T, bool)  {
	c := a + b
	if c < a {
		return c, true
	}

	return c, false
}

// halfCarryAdd8b checks and returns if the addition result of a and b results in a half carry
func halfCarryAdd8b(a, b byte) bool {
	return (((a & 0xF) + (b & 0xF)) & 0x10) == 0x10
}

// halfCarrySub8b checks and returns if the subtraction result of a and b results in a half carry
func halfCarrySub8b(a, b byte) bool {
	return (((a & 0xF) - (b & 0xF)) & 0x10) == 0x10
}

// JoinBytes join 2 bytes together
// TODO - check that this actually works properly
func JoinBytes(high, low byte) uint16 {
	return (uint16(high) << 8) ^ uint16(low)
}
