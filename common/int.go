package common

// IntPtr returns a pointer to the given int.
func IntPtr(i int) *int {
	return &i
}
