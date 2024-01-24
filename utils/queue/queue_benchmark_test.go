package queue

import (
	"reflect"
	"testing"
)

func Benchmark_sliceGrow_copy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 := make([]int, 10)
		s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
		s2 := []int{1, 2, 3, 4, 5}
		copy(s1[len(s2):], s1[:len(s1)-len(s2)])
		copy(s1[:len(s2)], s2)
	}
}

func Benchmark_sliceGrow_append(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 := make([]int, 10)
		s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
		s2 := []int{1, 2, 3, 4, 5}
		s1 = append(s2, s1[:5]...)
	}
}

func Benchmark_sliceGrow_copy_N_int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 := make([]int, 14500)
		s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
		s2 := []int{1, 2, 3, 4, 5}
		copy(s1[len(s2):], s1[:len(s1)-len(s2)])
		copy(s1[:len(s2)], s2)
	}
}

func Benchmark_sliceGrow_append_N_int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 := make([]int, 14500)
		s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
		s2 := []int{1, 2, 3, 4, 5}
		s1 = append(s2, s1[:5]...)
	}
}
func Benchmark_sliceGrow_copy_N_int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 := make([]int64, 15000)
		s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
		s2 := []int64{1, 2, 3, 4, 5}
		copy(s1[len(s2):], s1[:len(s1)-len(s2)])
		copy(s1[:len(s2)], s2)
	}
}

func Benchmark_sliceGrow_append_N_int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 := make([]int64, 15000)
		s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
		s2 := []int64{1, 2, 3, 4, 5}
		s1 = append(s2, s1[:5]...)
	}
}
func Benchmark_sliceGrow_copy_N_string(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 := make([]string, 15000)
		s1[0], s1[1], s1[2], s1[3], s1[4] = "6", "7", "8", "9", "10"
		s2 := []string{"1", "2", "3", "4", "5"}
		copy(s1[len(s2):], s1[:len(s1)-len(s2)])
		copy(s1[:len(s2)], s2)
	}
}

func Benchmark_sliceGrow_append_N_string(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 := make([]string, 15000)
		s1[0], s1[1], s1[2], s1[3], s1[4] = "6", "7", "8", "9", "10"
		s2 := []string{"1", "2", "3", "4", "5"}
		s1 = append(s2, s1[:5]...)
	}
}

func Test_sliceGrow_copy_1(t *testing.T) {
	s1 := make([]int, 10)
	s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
	s2 := []int{1}
	copy(s1[len(s2):], s1[:len(s1)-len(s2)])
	copy(s1[:len(s2)], s2)

	want := []int{1, 6, 7, 8, 9, 10}
	if !reflect.DeepEqual(s1[:len(want)], want) {
		t.Errorf("want = %+v, get = %+v", want, s1[:len(want)])
	} else {
		t.Logf("s1: %+v", s1[:len(want)])
	}
}
func Test_sliceGrow_copy_2(t *testing.T) {
	s1 := make([]int, 10)
	s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
	s2 := []int{1, 2}
	copy(s1[len(s2):], s1[:len(s1)-len(s2)])
	copy(s1[:len(s2)], s2)

	want := []int{1, 2, 6, 7, 8, 9, 10}
	if !reflect.DeepEqual(s1[:len(want)], want) {
		t.Errorf("want = %+v, get = %+v", want, s1[:len(want)])
	} else {
		t.Logf("s1: %+v", s1[:len(want)])
	}
}
func Test_sliceGrow_copy_5(t *testing.T) {
	s1 := make([]int, 10)
	s1[0], s1[1], s1[2], s1[3], s1[4] = 6, 7, 8, 9, 10
	s2 := []int{1, 2, 3, 4, 5}
	copy(s1[len(s2):], s1[:len(s1)-len(s2)])
	copy(s1[:len(s2)], s2)
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if !reflect.DeepEqual(s1[:len(want)], want) {
		t.Errorf("want = %+v, get = %+v", want, s1[:len(want)])
	} else {
		t.Logf("s1: %+v", s1[:len(want)])
	}
}
