package queue

import "errors"

type Queue[E any] interface {
	Add(...E) // add last
	AddFirst(...E)
	AddLast(...E)

	Poll(e []E) (n int)
	PollNext() *E
	PollNextN(int) []E

	Peek(e []E) (n int)
	PeekNext() *E
	PeekNextN(int) []E

	IsEmpty() bool
	Size() int
	Clear()
}

// ===================================================================================================
// ===================================================================================================
// ===================================================================================================

var ErrTooLarge = errors.New("utils.NonBlockingSliceQueue: too large")

const maxInt = int(^uint(0) >> 1)
const smallBufferSize = 64
