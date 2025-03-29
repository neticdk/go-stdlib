package xslices

// Filter returns a new slice containing only the elements of the slice that
// satisfy the predicate.
func Filter[T any](data []T, f func(T) bool) []T {
	fltd := make([]T, 0, len(data))

	for _, e := range data {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}
