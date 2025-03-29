package xslices

const (
	DefaultMaxIterations = 100000
	DefaultStep          = 1
)

type UnfoldConfig struct {
	Max  int
	Step int
}

type UnfoldOption func(*UnfoldConfig)

// WithStep sets the step of the unfold.
// This does not include the first step.
// If n is 5, the length of the result is <= 6.
func WithMax(n int) UnfoldOption {
	return func(c *UnfoldConfig) {
		c.Max = n
	}
}

// WithStep sets the step of the unfold.
// This does not include the first step.
// If step is 2, the result will include every second value.
func WithStep(step int) UnfoldOption {
	return func(c *UnfoldConfig) {
		c.Step = step
	}
}

// Unfold generates a slice by repeatedly applying a function to an accumulator.
// It includes the accumulator in the result as the first value.
// If the predicate is always false, it returns nil.
// It stops when the predicate returns false.
func Unfold[T any](acc T, f func(T) T, p func(T) bool, opts ...UnfoldOption) []T {
	config := &UnfoldConfig{
		Max:  DefaultMaxIterations,
		Step: DefaultStep,
	}

	for _, opt := range opts {
		opt(config)
	}

	current := 0
	var res []T
	for a := acc; p(a); a, current = f(a), current+1 {
		if current%config.Step == 0 {
			res = append(res, a)
		}
		if current >= config.Max {
			return res
		}
	}

	return res
}

// UnfoldI generates a slice by repeatedly applying a function to an accumulator.
// It includes the accumulator in the result as the first value.
// The length of the result is equal to i + 1.
// If i is negative, it returns nil.
// It stops after i iterations.
func UnfoldI[T any](acc T, f func(T) T, n int, opts ...UnfoldOption) []T {
	config := &UnfoldConfig{
		Max:  DefaultMaxIterations,
		Step: DefaultStep,
	}

	for _, opt := range opts {
		opt(config)
	}

	if n < 0 {
		return nil
	}

	n = min(n, config.Max)

	res := make([]T, 0, n)
	res = append(res, acc)
	for i := 1; i <= n; i++ {
		acc = f(acc)
		if i%config.Step == 0 {
			res = append(res, acc)
		}
		if i >= config.Max {
			return res
		}
	}

	return res
}
