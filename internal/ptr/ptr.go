package ptr

// To returns a pointer to the given value.
//
//go:fix inline
func To[T any](v T) *T {
	return new(v)
}
