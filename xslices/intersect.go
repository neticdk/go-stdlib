package xslices

// Intersect returns the intersection of two comparable slices.
func Intersect[T comparable](a, b []T) []T {
	seen := make(map[T]struct{})
	set := make([]T, 0)

	for _, v := range a {
		seen[v] = struct{}{}
	}

	for _, v := range b {
		if _, exists := seen[v]; exists {
			set = append(set, v)
			// Remove the element from the seen map to avoid duplicates
			delete(seen, v)
		}
	}

	return set
}
