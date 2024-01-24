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

func If[E any](b bool, t E, f E) E {
	if b {
		return t
	} else {
		return f
	}
}
