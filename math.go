package stellar

func max[T int | int64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T int | int64](a, b T) T {
	if a < b {
		return a
	}
	return b
}
