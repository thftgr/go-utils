package queue

import "sync"

type blockingLinkedQueueNode[E any] struct {
	element *E
	next    *blockingLinkedQueueNode[E]
}

// BlockingLinkedQueue is goroutine-safe
type BlockingLinkedQueue[E any] struct {
	first *blockingLinkedQueueNode[E]
	last  *blockingLinkedQueueNode[E]
	size  int

	mutex sync.Mutex
}

func (q *BlockingLinkedQueue[E]) Add(e ...E) {
	q.AddLast(e...) // addlast use mutex
}

func (q *BlockingLinkedQueue[E]) AddFirst(e ...E) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for i := len(e) - 1; i >= 0; i-- {
		q.first = &blockingLinkedQueueNode[E]{element: &e[i], next: q.first}
	}

	if q.last == nil {
		q.last = q.first
	}

	q.size += len(e)
}

func (q *BlockingLinkedQueue[E]) AddLast(e ...E) {
	if len(e) < 1 {
		return
	}

	if q.IsEmpty() {
		q.Clear()
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()
	// 임시 리스트 생성
	first := &blockingLinkedQueueNode[E]{element: &e[0]}
	last := first

	for i := range e[1:] {
		last.next = &blockingLinkedQueueNode[E]{element: &e[i+1]}
		last = last.next
	}
	if q.last == nil {
		q.last = last
		q.first = first
	} else {
		q.last.next = last
	}

	q.size += len(e)
}

func (q *BlockingLinkedQueue[E]) Poll(p []E) (n int) {
	if q.IsEmpty() {
		q.Clear()
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	m := len(p)
	if q.size < m {
		m = q.size
	}

	for i := 0; i < m; i++ {
		p[i] = *q.first.element
		q.first = q.first.next
	}
	q.size -= m
	return m
}

func (q *BlockingLinkedQueue[E]) PollNext() (e *E) {

	if q.IsEmpty() {
		q.Clear()
		return nil
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.size--
	e = q.first.element
	q.first = q.first.next
	return
}

func (q *BlockingLinkedQueue[E]) PollNextN(n int) (e []E) {
	e = make([]E, n)
	pn := q.Poll(e) // poll use mutex
	return e[:pn]
}

func (q *BlockingLinkedQueue[E]) Peek(p []E) (n int) {
	if q.IsEmpty() {
		q.Clear()
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	m := len(p)
	if q.size < m {
		m = q.size
	}
	first := q.first
	for i := 0; i < m; i++ {
		p[i] = *first.element
		first = first.next
	}
	return m
}

func (q *BlockingLinkedQueue[E]) PeekNext() *E {
	if q.IsEmpty() {
		q.Clear()
		return nil
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.first.element
}

func (q *BlockingLinkedQueue[E]) PeekNextN(n int) (e []E) {
	e = make([]E, n)
	pn := q.Peek(e) // peek use mutex
	return e[:pn]
}

func (q *BlockingLinkedQueue[E]) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.size < 1
}

func (q *BlockingLinkedQueue[E]) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.size
}

func (q *BlockingLinkedQueue[E]) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.first = nil
	q.last = nil
	q.size = 0
}
