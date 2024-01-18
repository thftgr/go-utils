package utils

import (
	"reflect"
	"testing"
)

// TODO lock 이 제대로 작동하고있는지에 대한 테스트 추가해야함

func TestBlockingQueue_Add(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{99, 98, 1}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{99, 98, 2}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{99, 98, 3}},
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

func TestBlockingQueue_AddFirst(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{1, 99, 98}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{2, 99, 98}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{3, 99, 98}},
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

func TestBlockingQueue_AddLast(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{99, 98, 1}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{99, 98, 2}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{99, 98, 3}},
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

func TestBlockingQueue_Clear(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}},
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

func TestBlockingQueue_IsEmpty(t *testing.T) {
	type testCase[E any] struct {
		name    string
		q       *BlockingSliceQueue[E]
		init    []E
		isEmpty bool
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{}, isEmpty: true},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{99, 98}, isEmpty: false},
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

func TestBlockingQueue_Peek(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{1, 99, 98}, want: []int{1, 99}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{2, 99, 98}, want: []int{2, 99}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{3, 99, 98}, want: []int{3, 99}},
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
func TestBlockingQueue_PeekNext(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		want E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{1, 99, 98}, want: 1},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{2, 99, 98}, want: 2},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{3, 99, 98}, want: 3},
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

func TestBlockingQueue_PeekNextN(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{1, 99, 98}, want: []int{1, 99}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{2, 99, 98}, want: []int{2, 99}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{3, 99, 98}, want: []int{3, 99}},
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

func TestBlockingQueue_Poll(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{1, 99, 98}, want: []int{1, 99}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{2, 99, 98}, want: []int{2, 99}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{3, 99, 98}, want: []int{3, 99}},
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

func TestBlockingQueue_PollNext(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{1, 99, 98}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{2, 99, 98}},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{3, 99, 98}},
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

func TestBlockingQueue_PollNextN(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		n    int
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{1, 99, 98}, n: 1},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{2, 99, 98}, n: 2},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{3, 99, 98}, n: 3},
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

func TestBlockingQueue_Size(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *BlockingSliceQueue[E]
		init []E
		size int
	}
	tests := []testCase[int]{
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{1}, size: 1},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{2, 99}, size: 2},
		{name: "", q: &BlockingSliceQueue[int]{}, init: []int{3, 99, 98}, size: 3},
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
