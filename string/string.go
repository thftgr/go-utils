package string

import "strings"

func SplitTrim(s, sep string) (res []string) {
	ss := strings.Split(s, sep)
	for _, i := range ss {
		res = append(res, strings.TrimSpace(i))
	}
	return
}
