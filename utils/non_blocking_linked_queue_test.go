package utils

import (
	"reflect"
	"testing"
)

func TestNonBlockingLinkedQueue_Add(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingLinkedQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{99, 98, 1}},
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{99, 98, 2}},
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{99, 98, 3}},
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

func TestNonBlockingLinkedQueue_AddFirst(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingLinkedQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{1, 99, 98}},
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{2, 99, 98}},
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{3, 99, 98}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.AddFirst(tt.args...)
			if peek := tt.q.PeekNextN(len(tt.want)); !reflect.DeepEqual(tt.want, peek) {
				t.Errorf("Add() want = %+v, get %+v", tt.want, peek)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_AddLast(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingLinkedQueue[E]
		init []E
		args []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{1}, want: []int{99, 98, 1}},
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{2}, want: []int{99, 98, 2}},
		{name: "", q: &NonBlockingLinkedQueue[int]{}, init: []int{99, 98}, args: []int{3}, want: []int{99, 98, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.AddLast(tt.args...)
			if peek := tt.q.PeekNextN(len(tt.want)); !reflect.DeepEqual(tt.want, peek) {
				t.Errorf("Add() want = %+v, get %+v", tt.want, peek)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_Clear(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingLinkedQueue[E]
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Clear()
		})
	}
}

func TestNonBlockingLinkedQueue_IsEmpty(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingLinkedQueue[E]
		want bool
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_Peek(t *testing.T) {
	type args[E any] struct {
		p []E
	}
	type testCase[E any] struct {
		name  string
		q     NonBlockingLinkedQueue[E]
		args  args[E]
		wantN int
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := tt.q.Peek(tt.args.p); gotN != tt.wantN {
				t.Errorf("Peek() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_PeekNext(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingLinkedQueue[E]
		want *E
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.PeekNext(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeekNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_PeekNextN(t *testing.T) {
	type args struct {
		n int
	}
	type testCase[E any] struct {
		name  string
		q     NonBlockingLinkedQueue[E]
		args  args
		wantE []E
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotE := tt.q.PeekNextN(tt.args.n); !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("PeekNextN() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_Poll(t *testing.T) {
	type args[E any] struct {
		p []E
	}
	type testCase[E any] struct {
		name  string
		q     NonBlockingLinkedQueue[E]
		args  args[E]
		wantN int
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := tt.q.Poll(tt.args.p); gotN != tt.wantN {
				t.Errorf("Poll() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_PollNext(t *testing.T) {
	type testCase[E any] struct {
		name  string
		q     NonBlockingLinkedQueue[E]
		wantE *E
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotE := tt.q.PollNext(); !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("PollNext() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_PollNextN(t *testing.T) {
	type args struct {
		n int
	}
	type testCase[E any] struct {
		name  string
		q     NonBlockingLinkedQueue[E]
		args  args
		wantE []E
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotE := tt.q.PollNextN(tt.args.n); !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("PollNextN() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNonBlockingLinkedQueue_Size(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingLinkedQueue[E]
		want int
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
