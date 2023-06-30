package utils

type List[T any] struct {
	list []T
}

func (l *List[T]) Size() int {
	return len(l.list)
}

func (l *List[T]) Add(value T) {
	l.list = append(l.list, value)
}

func (l *List[T]) AddAll(value ...T) {
	l.list = append(l.list, value...)
}

func (l *List[T]) Get(index int) (res T) {
	if len(l.list) < 1 {
		return
	}
	return l.list[index]
}

func (l *List[T]) Clear() {
	l.list = []T{}
}

func (l *List[T]) Poll() (res T) {
	if len(l.list) < 1 {
		return
	}
	res = l.list[0]
	l.list = l.list[1:]
	return
}
