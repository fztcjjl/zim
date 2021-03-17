package pqueue

import "container/heap"

type Comparable interface {
	CompareTo(Comparable) int
}

type heapImpl []Comparable

func (h heapImpl) Len() int {
	return len(h)
}

func (h heapImpl) Less(i, j int) bool {
	return h[i].CompareTo(h[j]) < 0
}

func (h heapImpl) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heapImpl) Push(x interface{}) {
	n := len(*h)
	c := cap(*h)
	if n+1 > c {
		npq := make(heapImpl, n, c*2)
		copy(npq, *h)
		*h = npq
	}
	*h = (*h)[0 : n+1]
	(*h)[n] = x.(Comparable)
}

func (h *heapImpl) Pop() interface{} {
	n := len(*h)
	c := cap(*h)
	if n < (c/2) && c > 25 {
		npq := make(heapImpl, n, c/2)
		copy(npq, *h)
		*h = npq
	}
	e := (*h)[n-1]
	*h = (*h)[0 : n-1]
	return e
}

func NewPriorityQueue(capacity int) *PriorityQueue {
	if capacity == 0 {
		capacity = 1
	}
	q := &PriorityQueue{h: make(heapImpl, 0, capacity)}

	return q
}

type PriorityQueue struct {
	h heapImpl
}

func (pq *PriorityQueue) Push(x Comparable) {
	heap.Push(&pq.h, x)
}

func (pq *PriorityQueue) Pop() Comparable {
	return heap.Pop(&pq.h).(Comparable)
}

func (pq *PriorityQueue) Front() Comparable {
	if pq.h.Len() == 0 {
		return nil
	}

	return pq.h[0].(Comparable)
}

func (pq *PriorityQueue) Remove(i int) Comparable {
	return heap.Remove(&pq.h, i).(Comparable)
}
