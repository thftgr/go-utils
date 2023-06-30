package utils

func GetOrElse[V any](value V, orElse V) V {
	if value != nil {
		return value
	}
	return orElse
}
