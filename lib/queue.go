package aoc

// Queue represents a queue of items
type Queue[T comparable] struct {
	s []T
}

// NewQueue returns a new queue
func NewQueue[T comparable]() Queue[T] {
	return Queue[T]{}
}

// Push adds an element at the end of the slice
func (q *Queue[T]) Push(e T) {
	q.s = append(q.s, e)
}

// Pop removes the first element of the queue
func (q *Queue[T]) Pop() T {
	e := q.s[0]
	if len(q.s) == 1 {
		q.s = []T{}
	} else {
		q.s = q.s[1:q.Len()]
	}
	return e
}

// Last removes the last element of the queue
func (q *Queue[T]) Last() T {
	lastIdx := q.Len() - 1
	e := q.s[lastIdx]
	if len(q.s) == 1 {
		q.s = []T{}
	} else {
		q.s = q.s[0:lastIdx]
	}
	return e
}

// Len returns the length of the set
func (q Queue[T]) Len() int {
	return len(q.s)
}
