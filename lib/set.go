package aoc

// Set represents a set of data
type Set[T comparable] map[T]struct{}

// NewSet returns a new set
func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

// Add adds an element to the set
func (s *Set[T]) Reset() {
	*s = map[T]struct{}{}
}

// Add adds an element to the set
func (s Set[T]) Add(k T) {
	s[k] = struct{}{}
}

// Remove removes an element from the set
func (s Set[T]) Remove(k T) {
	delete(s, k)
}

// Len returns the length of the set
func (s Set[T]) Len() int {
	return len(s)
}

// Has returns true if the element is in the set
func (s Set[T]) Has(e T) bool {
	_, ok := s[e]
	return ok
}

// Slice returns a slice of T
func (s Set[T]) Slice() []T {
	slice := make([]T, s.Len())
	i := 0
	for k := range s {
		slice[i] = k
		i++
	}
	return slice
}
