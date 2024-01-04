package header

import (
	"regexp"
	"sort"
	"strconv"
)

var REGEXP_ACCEPT, _ = regexp.Compile(`(?:(([A-z-*]+)\/([A-z0-9.+*-]+))(?:;(?:\s)?(?:q=([01](?:[.][0-9])?)))?)`) // accept

type Accept struct {
	MimeType MimeType
	Type     string
	Subtype  string
	Weight   float32
}

func ParseAccept(str string) []Accept {
	var accepts []Accept

	// Accept 헤더 파싱
	matches := REGEXP_ACCEPT.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		if len(match) > 0 {
			accept := Accept{
				MimeType: MimeType(match[1]),
				Type:     match[2],
				Subtype:  match[3],
			}

			if len(match) > 4 && match[4] != "" {
				if weight, err := strconv.ParseFloat(match[4], 32); err == nil {
					accept.Weight = float32(weight)
				}
			} else {
				accept.Weight = 1.0 // 기본 가중치
			}

			accepts = append(accepts, accept)
		}
	}
	sort.Slice(accepts, func(i, j int) bool {
		return accepts[i].Weight > accepts[j].Weight
	})
	return accepts
}
