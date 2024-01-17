package utils

import (
	"reflect"
	"testing"
)

func Test_cap_len(t *testing.T) {
	s1 := make([]int, 10)
	t.Logf("cap: %d", cap(s1))
	t.Logf("len: %d", len(s1))
}

func TestNonBlockingQueue_Add(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingQueue[E]
		init []E
		e    []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1}, want: []int{99, 98, 1}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 3}, want: []int{99, 98, 1, 2, 3}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 4, 5}, want: []int{99, 98, 1, 2, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.Add(tt.e...)
			if queued := tt.q.PeekNextN(len(tt.want)); !reflect.DeepEqual(queued, tt.want) {
				t.Errorf("Add() want = %+v, get %+v", tt.want, queued)
			}
		})
	}
}

func TestNonBlockingQueue_AddFirst(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingQueue[E]
		init []E
		e    []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1}, want: []int{1, 99, 98}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 3}, want: []int{1, 2, 3, 99, 98}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 4, 5}, want: []int{1, 2, 4, 5, 99, 98}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.AddFirst(tt.e...)
			if queued := tt.q.PeekNextN(len(tt.want)); !reflect.DeepEqual(queued, tt.want) {
				t.Errorf("AddFirst() want = %+v, get %+v", tt.want, queued)
			}
		})
	}
}

func TestNonBlockingQueue_AddLast(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingQueue[E]
		init []E
		e    []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1}, want: []int{99, 98, 1}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 3}, want: []int{99, 98, 1, 2, 3}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 4, 5}, want: []int{99, 98, 1, 2, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.AddLast(tt.e...)
			if queued := tt.q.PeekNextN(len(tt.want)); !reflect.DeepEqual(queued, tt.want) {
				t.Errorf("AddLast() want = %+v, get %+v", tt.want, queued)
			}
		})
	}
}

func TestNonBlockingQueue_Clear(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingQueue[E]
		init []E
		e    []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 3}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			tt.q.Clear()
			if e := tt.q.PeekNext(); e != nil {
				t.Errorf("Clear() want = nil, get %+v", e)
			}
		})
	}
}

func TestNonBlockingQueue_IsEmpty(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingQueue[E]
		init []E
		e    []E
		want bool
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1}, want: false},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 3}, want: false},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{99, 98}, e: []int{1, 2, 4, 5}, want: false},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{}, e: []int{}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			if e := tt.q.IsEmpty(); e != tt.want {
				t.Errorf("IsEmpty() want = true, get %+v", e)
			}
		})
	}
}

func TestNonBlockingQueue_Peek(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingQueue[E]
		init []E
		want []E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{1}, want: []int{1}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{1, 2, 4, 5}, want: []int{1, 2, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Add(tt.init...)
			p1 := make([]int, len(tt.want))
			p2 := make([]int, len(tt.want))
			tt.q.Peek(p1)
			n2 := tt.q.Peek(p2)
			p2 = p2[:n2]
			if !reflect.DeepEqual(tt.want, p2) {
				t.Errorf("Add() want = %+v, get %+v", tt.want, p2)
			}
		})
	}
}
func TestNonBlockingQueue_PeekNext(t *testing.T) {
	type testCase[E any] struct {
		name string
		q    *NonBlockingQueue[E]
		init []E
		want E
	}
	tests := []testCase[int]{
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{1}, want: 1},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{2, 2, 3}, want: 2},
		{name: "", q: &NonBlockingQueue[int]{}, init: []int{3, 2, 4, 5}, want: 3},
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

//
//func TestNonBlockingQueue_Poll(t *testing.T) {
//	type args[E any] struct {
//		e []E
//	}
//	type testCase[E any] struct {
//		name    string
//		q       NonBlockingQueue[E]
//		args    args[E]
//		wantN   int
//		wantErr bool
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotN, err := tt.q.Poll(tt.args.e)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Poll() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if gotN != tt.wantN {
//				t.Errorf("Poll() gotN = %v, want %v", gotN, tt.wantN)
//			}
//		})
//	}
//}
//
//func TestNonBlockingQueue_PollNext(t *testing.T) {
//	type args struct {
//		n int
//	}
//	type testCase[E any] struct {
//		name string
//		q    NonBlockingQueue[E]
//		args args
//		want []E
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.q.PollNext(tt.args.n); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("PollNext() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNonBlockingQueue_Size(t *testing.T) {
//	type testCase[E any] struct {
//		name string
//		q    NonBlockingQueue[E]
//		want int
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.q.Size(); got != tt.want {
//				t.Errorf("Size() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNonBlockingQueue_grow(t *testing.T) {
//	type args struct {
//		n int
//	}
//	type testCase[E any] struct {
//		name string
//		q    NonBlockingQueue[E]
//		args args
//		want int
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.q.grow(tt.args.n); got != tt.want {
//				t.Errorf("grow() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNonBlockingQueue_growSlice(t *testing.T) {
//	type args[E any] struct {
//		b []E
//		n int
//	}
//	type testCase[E any] struct {
//		name string
//		q    NonBlockingQueue[E]
//		args args[E]
//		want []E
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.q.growSlice(tt.args.b, tt.args.n); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("growSlice() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNonBlockingQueue_tryGrowByReslice(t *testing.T) {
//	type args struct {
//		n int
//	}
//	type testCase[E any] struct {
//		name  string
//		q     NonBlockingQueue[E]
//		args  args
//		want  int
//		want1 bool
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := tt.q.tryGrowByReslice(tt.args.n)
//			if got != tt.want {
//				t.Errorf("tryGrowByReslice() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("tryGrowByReslice() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
