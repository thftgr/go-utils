package header

import (
	"regexp"
)

var REGEXP_CONTENT_TYPE, _ = regexp.Compile(`(([A-z-*]+)\/([A-z0-9.+*-]+))`)                                 // ?/?
var REGEXP_CONTENT_TYPE_CHARSET, _ = regexp.Compile(`(?:(?:[cC][hH][aA][rR][sS][eE][tT])(?:=)([A-z0-9-]+))`) // charset

type ContentType struct {
	MimeType MimeType
	Type     string
	Subtype  string
	Charset  string
}

func ParseContentType(str string) *ContentType {
	contentType := &ContentType{}

	// Content-Type 파싱
	if matches := REGEXP_CONTENT_TYPE.FindStringSubmatch(str); len(matches) > 0 {
		contentType.MimeType = MimeType(matches[1])
		contentType.Type = matches[2]
		contentType.Subtype = matches[3]
	}

	// Charset 파싱
	if matches := REGEXP_CONTENT_TYPE_CHARSET.FindStringSubmatch(str); len(matches) > 0 {
		contentType.Charset = matches[1]
	}

	return contentType
}
