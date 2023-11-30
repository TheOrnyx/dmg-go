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
