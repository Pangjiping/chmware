package hash

type uints []uint32

// Len returns the length of the uints array.
func (x uints) Len() int {
	return len(x)
}

// Less returns true if element i is less than element j.
func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

// Swap swap uints[i] and uints[j].
func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
