package queue

// BlockingSliceQueue is not goroutine-safe
type BlockingSliceQueue[E any] struct {
	queue []E
	off   int
}

func (q *BlockingSliceQueue[E]) Add(e ...E) {
	q.AddLast(e...)
}

func (q *BlockingSliceQueue[E]) AddFirst(e ...E) {
	es := len(e)
	_, ok := q.tryGrowByReslice(len(e))
	if !ok {
		_ = q.grow(len(e))
	}
	copy(q.queue[es:], q.queue[:q.Size()-es]) // [q1,q2,q3,nil,nil,nil] => [q1,q2,q3,q1,q2,q3]
	copy(q.queue[:es], e)                     // [e1,e2,e3] + [q1,q2,q3,q1,q2,q3] => [q1,q2,q3,q1,q2,q3]
}

func (q *BlockingSliceQueue[E]) AddLast(e ...E) {
	m, ok := q.tryGrowByReslice(len(e))
	if !ok {
		m = q.grow(len(e))
	}
	copy(q.queue[m:], e)
}

func (q *BlockingSliceQueue[E]) Poll(e []E) (n int) {
	if q.IsEmpty() {
		q.Clear()
		return
	}
	n = copy(e, q.queue[q.off:])
	q.off += n
	return
}

func (q *BlockingSliceQueue[E]) PollNext() *E {
	if q.IsEmpty() {
		q.Clear()
		return nil
	}

	data := q.queue[q.off : q.off+1]
	q.off++
	return &data[0]
}
func (q *BlockingSliceQueue[E]) PollNextN(n int) (e []E) {
	e = make([]E, n)
	pn := q.Poll(e)
	return e[:pn]
}

func (q *BlockingSliceQueue[E]) Peek(e []E) (n int) {
	if q.IsEmpty() {
		q.Clear()
		return
	}
	return copy(e, q.queue[q.off:])
}

func (q *BlockingSliceQueue[E]) PeekNext() *E {
	if q.IsEmpty() {
		q.Clear()
		return nil
	}

	data := q.queue[q.off : q.off+1]
	return &data[0]
}
func (q *BlockingSliceQueue[E]) PeekNextN(n int) (e []E) {
	e = make([]E, n)
	pn := q.Peek(e)
	return e[:pn]
}

func (q *BlockingSliceQueue[E]) IsEmpty() bool {
	return len(q.queue) <= q.off
}

func (q *BlockingSliceQueue[E]) Size() int {
	return len(q.queue) - q.off
}

func (q *BlockingSliceQueue[E]) Clear() {
	q.queue = q.queue[:0]
	q.off = 0
}

func (q *BlockingSliceQueue[E]) tryGrowByReslice(n int) (int, bool) {
	if l := len(q.queue); n <= cap(q.queue)-l {
		q.queue = q.queue[:l+n]
		return l, true
	}
	return 0, false
}

func (q *BlockingSliceQueue[E]) grow(n int) int {
	m := q.Size()
	// If buffer is empty, reset to recover space.
	if m == 0 && q.off != 0 {
		q.Clear()
	}
	// Try to grow by means of a reslice.
	if i, ok := q.tryGrowByReslice(n); ok {
		return i
	}
	if q.queue == nil && n <= smallBufferSize {
		q.queue = make([]E, n, smallBufferSize)
		return 0
	}
	c := cap(q.queue)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(q.queue, q.queue[q.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Add b.off to account for b.buf[:b.off] being sliced off the front.
		q.queue = q.growSlice(q.queue[q.off:], q.off+n)
	}
	// Restore b.off and len(b.buf).
	q.off = 0
	q.queue = q.queue[:m+n]
	return m
}

func (q *BlockingSliceQueue[E]) growSlice(b []E, n int) []E {
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	c := len(b) + n
	if c < 2*cap(b) {
		c = 2 * cap(b)
	}
	b2 := append([]E(nil), make([]E, c)...)
	copy(b2, b)
	return b2[:len(b)]
}
