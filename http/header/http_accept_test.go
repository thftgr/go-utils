package header

import (
	"reflect"
	"testing"
)

func TestParseAccept(t *testing.T) {
	tests := []struct {
		name   string
		accept string
		want   []Accept
	}{
		{name: "TestParseAccept-001", accept: "text/html, application/xhtml+xml, */*;q=0.8, application/xml;q=0.9", want: []Accept{
			{"text/html", "text", "html", 1},
			{"application/xhtml+xml", "application", "xhtml+xml", 1},
			{"application/xml", "application", "xml", 0.9},
			{"*/*", "*", "*", 0.8},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseAccept(tt.accept); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAccept() = %v, want %v", got, tt.want)
			}
		})
	}
}
