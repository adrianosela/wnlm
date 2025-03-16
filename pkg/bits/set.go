package bits

// Integer represents all signed and unsigned integer types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// AreSet returns true if all bits on b are set on base.
func AreSet[N Integer](base N, b N) bool { return base&b == b }
