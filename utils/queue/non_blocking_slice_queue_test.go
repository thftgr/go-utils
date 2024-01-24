package queue

import (
	"reflect"
	"testing"
)

func TestNonBlockingQueue_Add(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{99, 98, 1}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{99, 98, 2}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{99, 98, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.Add(tt.args...)
			if peek := tt.q.PeekNextN(len(tt.want)); !reflect.DeepEqual(tt.want, peek) {
				t.Errorf("Add() want = %+v, get %+v", tt.want, peek)
			}
		})
	}
}

func TestNonBlockingQueue_AddFirst(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{1, 99, 98}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{2, 99, 98}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{3, 99, 98}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.AddFirst(tt.args...)
			if peek := tt.q.PeekNextN(len(tt.want)); !reflect.DeepEqual(tt.want, peek) {
				t.Errorf("AddFirst() want = %+v, get %+v", tt.want, peek)
			}
		})
	}
}

func TestNonBlockingQueue_AddLast(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{99, 98, 1}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{99, 98, 2}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{99, 98, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.AddLast(tt.args...)
			if peek := tt.q.PeekNextN(len(tt.want)); !reflect.DeepEqual(tt.want, peek) {
				t.Errorf("AddLast() want = %+v, get %+v", tt.want, peek)
			}
		})
	}
}

func TestNonBlockingQueue_Clear(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.Clear()
			if !tt.q.IsEmpty() {
				t.Errorf("Clear() want = true, get %+v", tt.q.IsEmpty())
			}
		})
	}
}

func TestNonBlockingQueue_IsEmpty(t *testing.T) {
	type testCase[E any] struct {
		name    string
		q       *NonBlockingSliceQueue[E]
		init    []E
		isEmpty bool
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{}, isEmpty: true},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{99, 98}, isEmpty: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			if tt.q.IsEmpty() != tt.isEmpty {
				t.Errorf("IsEmpty() want = %+v, get %+v", tt.isEmpty, tt.q.IsEmpty())
			}
		})
	}
}

func TestNonBlockingQueue_Peek(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{1, 99, 98}, want: []int{1, 99}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{2, 99, 98}, want: []int{2, 99}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{3, 99, 98}, want: []int{3, 99}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			peek := make([]int, len(tt.want))
			tt.q.PeekNextN(len(tt.want))
			tt.q.Peek(peek)
			if !reflect.DeepEqual(tt.want, peek) {
				t.Errorf("Peek() want = %+v, get %+v", tt.want, peek)
			}
		})
	}
}
func TestNonBlockingQueue_PeekNext(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		want E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{1, 99, 98}, want: 1},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{2, 99, 98}, want: 2},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{3, 99, 98}, want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.PeekNext()
			peek := tt.q.PeekNext()
			if !reflect.DeepEqual(&tt.want, peek) {
				if peek == nil {
					t.Errorf("PeekNext() want = %+v, get nil", tt.want)
				} else {
					t.Errorf("PeekNext() want = %+v, get %+v", tt.want, *peek)
				}
			}
		})
	}
}

func TestNonBlockingQueue_PeekNextN(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{1, 99, 98}, want: []int{1, 99}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{2, 99, 98}, want: []int{2, 99}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{3, 99, 98}, want: []int{3, 99}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.PeekNextN(len(tt.want))
			peek := tt.q.PeekNextN(len(tt.want))
			if !reflect.DeepEqual(tt.want, peek) {
				t.Errorf("PeekNextN() want = %+v, get %+v", tt.want, peek)
			}
		})
	}
}

func TestNonBlockingQueue_Poll(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{1, 99, 98}, want: []int{1, 99}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{2, 99, 98}, want: []int{2, 99}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{3, 99, 98}, want: []int{3, 99}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			poll := make([]int, len(tt.want))
			tt.q.Poll(poll)
			if !reflect.DeepEqual(tt.want, poll) {
				t.Errorf("Poll() want = %+v, get %+v", tt.want, poll)
			} else {
				peek := tt.q.PeekNextN(len(tt.want))
				if reflect.DeepEqual(tt.want, peek) {
					t.Errorf("PeekNextN() want != %+v, get %+v", tt.want, peek)
				}
			}
		})
	}
}

func TestNonBlockingQueue_PollNext(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{1, 99, 98}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{2, 99, 98}},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{3, 99, 98}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			on := tt.q.PollNext()
			pn := tt.q.PeekNext()
			if reflect.DeepEqual(on, pn) {
				t.Errorf("PollNext() get = %+v, peek = %+v", on, pn)
			}
		})
	}
}

func TestNonBlockingQueue_PollNextN(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		n    int
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{1, 99, 98}, n: 1},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{2, 99, 98}, n: 2},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{3, 99, 98}, n: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			on := tt.q.PollNextN(tt.n)
			pn := tt.q.PeekNextN(tt.n)
			if reflect.DeepEqual(on, pn) {
				t.Errorf("PollNextN() get = %+v, peek = %+v", on, pn)
			}
		})
	}
}

func TestNonBlockingQueue_Size(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingSliceQueue[E]
		init []E
		size int
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{1}, size: 1},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{2, 99}, size: 2},
		{name: "", q: &NonBlockingSliceQueue[int]{}, init: []int{3, 99, 98}, size: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			if tt.q.Size() != tt.size {
				t.Errorf("Size() get = %+v, want = %+v", tt.q.Size(), tt.size)
			}
		})
	}
}
