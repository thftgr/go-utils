package queue

type nonBlockingLinkedQueueNode[E any] struct {
	element *E
	next    *nonBlockingLinkedQueueNode[E]
}

// NonBlockingLinkedQueue is not goroutine-safe
type NonBlockingLinkedQueue[E any] struct {
	first *nonBlockingLinkedQueueNode[E]
	last  *nonBlockingLinkedQueueNode[E]
	size  int
}

func (q *NonBlockingLinkedQueue[E]) Add(e ...E) {
	q.AddLast(e...)
}

func (q *NonBlockingLinkedQueue[E]) AddFirst(e ...E) {
	for i := len(e) - 1; i >= 0; i-- {
		q.first = &nonBlockingLinkedQueueNode[E]{element: &e[i], next: q.first}
	}

	if q.last == nil {
		q.last = q.first
	}

	q.size += len(e)
}

func (q *NonBlockingLinkedQueue[E]) AddLast(e ...E) {
	if len(e) < 1 {
		return
	}

	if q.IsEmpty() {
		q.Clear()
	}

	// 임시 리스트 생성
	first := &nonBlockingLinkedQueueNode[E]{element: &e[0]}
	last := first

	for i := range e[1:] {
		last.next = &nonBlockingLinkedQueueNode[E]{element: &e[i+1]}
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

func (q *NonBlockingLinkedQueue[E]) Poll(p []E) (n int) {
	if q.IsEmpty() {
		q.Clear()
		return
	}
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

func (q *NonBlockingLinkedQueue[E]) PollNext() (e *E) {
	if q.IsEmpty() {
		q.Clear()
		return nil
	}

	q.size--
	e = q.first.element
	q.first = q.first.next
	return
}

func (q *NonBlockingLinkedQueue[E]) PollNextN(n int) (e []E) {
	e = make([]E, n)
	pn := q.Poll(e)
	return e[:pn]
}

func (q *NonBlockingLinkedQueue[E]) Peek(p []E) (n int) {
	if q.IsEmpty() {
		q.Clear()
		return
	}
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

func (q *NonBlockingLinkedQueue[E]) PeekNext() *E {
	if q.IsEmpty() {
		q.Clear()
		return nil
	}
	return q.first.element
}

func (q *NonBlockingLinkedQueue[E]) PeekNextN(n int) (e []E) {
	e = make([]E, n)
	pn := q.Peek(e)
	return e[:pn]
}

func (q *NonBlockingLinkedQueue[E]) IsEmpty() bool {
	return q.size < 1
}

func (q *NonBlockingLinkedQueue[E]) Size() int {
	return q.size
}

func (q *NonBlockingLinkedQueue[E]) Clear() {
	q.first = nil
	q.last = nil
	q.size = 0
}
