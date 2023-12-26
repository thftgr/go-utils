package utils

func In[E comparable](t E, v ...E) bool {
	for i := range v {
		if t == v[i] {
			return true
		}
	}
	return false
}

func And[E comparable](t E, v ...E) bool {
	for i := range v {
		if t != v[i] {
			return false
		}
	}
	return true
}
