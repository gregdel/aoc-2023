package aoc

import "container/heap"

type QueueElement[T comparable] struct {
	e        T
	index    int
	priority int
}

type PriorityQueue[T comparable] struct {
	pq priorityQueue[T]
}

func NewPriorityQueue[T comparable]() PriorityQueue[T] {
	pq := priorityQueue[T]{}
	heap.Init(&pq)
	return PriorityQueue[T]{
		pq: pq,
	}
}

func (pq *PriorityQueue[T]) Len() int {
	return pq.pq.Len()
}

func (pq *PriorityQueue[T]) Push(e T, p int) {
	qe := &QueueElement[T]{
		e:        e,
		priority: p,
	}
	heap.Push(&pq.pq, qe)
}

func (pq *PriorityQueue[T]) Pop() (T, int) {
	x := heap.Pop(&pq.pq)
	item := x.(*QueueElement[T])
	return item.e, item.priority
}

type priorityQueue[T comparable] []*QueueElement[T]

func (pq priorityQueue[T]) Len() int { return len(pq) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the lowest priority first
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*QueueElement[T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
