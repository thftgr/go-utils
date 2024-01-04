package header

import "strings"

type MimeType string

func (r MimeType) Equals(a any) bool {
	if v, ok := a.(string); !ok {
		return false
	} else {
		return string(r) == strings.ToLower(v)
	}
}

const (
	CONTENT_TYPE_APPLICATION_JSON MimeType = `application/json`
	CONTENT_TYPE_APPLICATION_XML  MimeType = `application/xml`
	CONTENT_TYPE_TEXT_HTML        MimeType = `text/html`
)
