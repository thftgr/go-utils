package bean

import (
	"reflect"
)

var bean = map[string]any{}

func Register[E any](a E) {
	if reflect.TypeOf(a).Kind() != reflect.Pointer {
		panic("not allowed non pointer type")
	}
	bean[reflect.TypeOf((*E)(nil)).String()] = a
}

func Get[E any]() E {
	return bean[reflect.TypeOf((*E)(nil)).String()].(E)
}
