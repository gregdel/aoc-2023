package aoc

// Result is a helper to get results
type Result struct {
	v []string
}

// NewResult returns a new result
func NewResult(rt1, r1, rt2, r2 string) *Result {
	return &Result{
		v: []string{rt1, r1, rt2, r2},
	}
}

// Expect returns the expected result of a test
func (r *Result) Expect(part int, test bool) string {
	idx := 0

	if part == 2 {
		idx += 2
	}

	if !test {
		idx++
	}

	return r.v[idx]
}
