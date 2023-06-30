package utils

// Coalesce (nil,nil,"aa",nil) => "aa"
func Coalesce[V any](v ...V) (res V) {
	for _, val := range v {
		if val != nil {
			return val
		}
	}
	return
}
