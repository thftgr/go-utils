package header

import (
	"reflect"
	"testing"
)

var (
	APPLICATION_JSON  MimeType = "application/json"
	APPLICATION_JSON2 MimeType = "application/json"
	APPLICATION_ALL   MimeType = "application/*"
	APPLICATION_XML   MimeType = "application/xml"
	TEXT_HTML         MimeType = "text/html"
	ALL_ALL           MimeType = "*/*"
	ALL_JSON          MimeType = "*/json"
	APPLICATION_VOID  MimeType = "application/"
	VOID_JSON         MimeType = "/json"
)

func TestMimeType_Match(t *testing.T) {
	tests := []struct {
		name string
		r    MimeType
		args []MimeType
		want *MimeType
	}{
		{name: "", r: APPLICATION_JSON, args: []MimeType{APPLICATION_JSON, APPLICATION_XML, TEXT_HTML}, want: &APPLICATION_JSON},
		{name: "", r: APPLICATION_JSON, args: []MimeType{APPLICATION_JSON, APPLICATION_XML, TEXT_HTML}, want: &APPLICATION_JSON2},

		{name: "", r: ALL_ALL, args: []MimeType{APPLICATION_JSON, APPLICATION_XML, TEXT_HTML}, want: &APPLICATION_JSON},
		{name: "", r: ALL_ALL, args: []MimeType{APPLICATION_XML, APPLICATION_JSON, TEXT_HTML}, want: &APPLICATION_XML},
		{name: "", r: ALL_ALL, args: []MimeType{APPLICATION_XML, TEXT_HTML, APPLICATION_JSON}, want: &APPLICATION_XML},

		{name: "", r: APPLICATION_JSON, args: []MimeType{APPLICATION_XML, TEXT_HTML}, want: nil},
		{name: "", r: APPLICATION_XML, args: []MimeType{APPLICATION_JSON, TEXT_HTML}, want: nil},
		{name: "", r: TEXT_HTML, args: []MimeType{APPLICATION_JSON, APPLICATION_XML}, want: nil},

		{name: "", r: ALL_JSON, args: []MimeType{APPLICATION_JSON, APPLICATION_XML, TEXT_HTML}, want: nil},

		{name: "", r: APPLICATION_ALL, args: []MimeType{APPLICATION_JSON, APPLICATION_XML, TEXT_HTML}, want: &APPLICATION_JSON},
		{name: "", r: APPLICATION_ALL, args: []MimeType{APPLICATION_XML, APPLICATION_JSON, TEXT_HTML}, want: &APPLICATION_XML},
		{name: "", r: APPLICATION_ALL, args: []MimeType{TEXT_HTML, APPLICATION_JSON, APPLICATION_XML}, want: &APPLICATION_JSON},
		{name: "", r: APPLICATION_ALL, args: []MimeType{TEXT_HTML}, want: nil},

		{name: "", r: APPLICATION_VOID, args: []MimeType{APPLICATION_JSON, APPLICATION_XML, TEXT_HTML}, want: nil},
		{name: "", r: VOID_JSON, args: []MimeType{APPLICATION_JSON, APPLICATION_XML, TEXT_HTML}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Match(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
