package utils

import (
	"regexp"
	"sort"
	"strconv"
)

var REGEXP_CONTENT_TYPE, _ = regexp.Compile(`(([A-z-*]+)\/([A-z0-9.+*-]+))`)                                     // ?/?
var REGEXP_CONTENT_TYPE_CHARSET, _ = regexp.Compile(`(?:(?:[cC][hH][aA][rR][sS][eE][tT])(?:=)([A-z0-9-]+))`)     // charset
var REGEXP_ACCEPT, _ = regexp.Compile(`(?:(([A-z-*]+)\/([A-z0-9.+*-]+))(?:;(?:\s)?(?:q=([01](?:[.][0-9])?)))?)`) // accept

type ContentType struct {
	ContentType string
	Type        string
	Subtype     string
	Charset     string
}

type Accept struct {
	ContentType string
	Type        string
	Subtype     string
	Weight      float32
}

func ParseContentType(str string) *ContentType {
	contentType := &ContentType{}

	// Content-Type 파싱
	if matches := REGEXP_CONTENT_TYPE.FindStringSubmatch(str); len(matches) > 0 {
		contentType.ContentType = matches[1]
		contentType.Type = matches[2]
		contentType.Subtype = matches[3]
	}

	// Charset 파싱
	if matches := REGEXP_CONTENT_TYPE_CHARSET.FindStringSubmatch(str); len(matches) > 0 {
		contentType.Charset = matches[1]
	}

	return contentType
}

func ParseAccept(str string) []Accept {
	var accepts []Accept

	// Accept 헤더 파싱
	matches := REGEXP_ACCEPT.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		if len(match) > 0 {
			accept := Accept{
				ContentType: match[1],
				Type:        match[2],
				Subtype:     match[3],
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
