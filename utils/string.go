package utils

import "strings"

func SplitTrim(s, sep string) (res []string) {
	ss := strings.Split(s, sep)
	for _, i := range ss {
		res = append(res, strings.TrimSpace(i))
	}
	return
}

func SplitTrimUpper(s, sep string) (res []string) {
	ss := strings.Split(strings.ToUpper(s), sep)
	for _, i := range ss {
		res = append(res, strings.TrimSpace(i))
	}
	return
}

func SplitTrimLower(s, sep string) (res []string) {
	ss := strings.Split(strings.ToLower(s), sep)
	for _, i := range ss {
		res = append(res, strings.TrimSpace(i))
	}
	return
}

func TrimLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func TrimUpper(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
