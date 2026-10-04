package collection

// mayHavePointers conservatively identifies types needing retired-slot clearing.
func mayHavePointers[T any]() bool {
	switch any((*T)(nil)).(type) {
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *uintptr, *float32, *float64, *complex64, *complex128, *bool:
		return false
	}
	return true
}
