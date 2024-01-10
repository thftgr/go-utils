package header

import "strings"

const (
	CONTENT_TYPE_APPLICATION_ALL          MimeType = `application/*`
	CONTENT_TYPE_APPLICATION_JSON         MimeType = `application/json`
	CONTENT_TYPE_APPLICATION_XML          MimeType = `application/xml`
	CONTENT_TYPE_APPLICATION_OCTET_STREAM MimeType = `application/octet-stream`
	CONTENT_TYPE_TEXT_ALL                 MimeType = `text/*`
	CONTENT_TYPE_TEXT_HTML                MimeType = `text/html`
	CONTENT_TYPE_ALL                      MimeType = `*/*`
)

// MimeType "application/json" 형태
type MimeType string

func (r MimeType) String() string {
	return string(r)
}
func (r MimeType) Valid() bool {
	s := strings.Split(strings.TrimSpace(string(r)), "/")
	if len(s) != 2 {
		return false
	}
	if s[0] == "" || s[1] == "" { // "/", "application/", "/json"
		return false
	}
	if s[0] == "*" && s[1] != "*" { // "*/json"
		return false
	}
	return true

}
func (r MimeType) Equals(a any) bool {
	if v, ok := a.(string); !ok {
		return false
	} else {
		return strings.ToLower(r.String()) == strings.ToLower(v)
	}
}
func (r MimeType) Type() string {
	s := strings.Split(strings.TrimSpace(string(r)), "/")
	if len(s) >= 2 {
		return s[0]
	}
	return "*"
}

func (r MimeType) SubType() string {
	s := strings.Split(strings.TrimSpace(string(r)), "/")
	if len(s) >= 2 {
		return s[1]
	}
	return "*"
}

// Match nil 인경우 매칭되는 값이 없거나 유효하지 않다고 보면 됨,
func (r MimeType) Match(a ...MimeType) *MimeType {
	if !r.Valid() {
		return nil
	}

	if string(r) == "*/*" { // 와일드카드인경우 가장 먼저 매칭되는 값을 리턴
		for i := range a {
			if a[i].Valid() {
				return &a[i]
			}
		}
		return nil
	}

	t := strings.ToLower(r.Type())
	st := strings.ToLower(r.SubType())
	wt := t == "*"
	wst := st == "*"
	for _, a := range a {

		if !a.Valid() {
			continue
		}

		at := strings.ToLower(a.Type())
		ast := strings.ToLower(a.SubType())
		if (wt || t == at) && (wst || st == ast) { // 타입과 서브타입이 일치하거나 와일드카드인 경우
			return &a
		}

	}
	return nil
}
