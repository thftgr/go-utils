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
