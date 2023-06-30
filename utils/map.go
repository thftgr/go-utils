package utils

// GetAllInMap It returns the value corresponding to the key existing in the map as a slice.
func GetAllInMap[K comparable, V comparable](m map[K]V, keys []K) (res []V) {
	if m == nil || keys == nil || len(m) < 1 || len(keys) < 1 {
		return
	}
	for _, key := range keys {
		res = append(res, m[key])
	}
	return
}

type Map[K comparable, V any] struct {
	m map[K]V
}

func (m *Map[K, V]) Put(key K, value V) *Map[K, V] {
	m.m[key] = value
	return m
}
func (m *Map[K, V]) Get(key K) V {
	return m.m[key]
}
