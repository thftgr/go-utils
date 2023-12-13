package logger

import "testing"

func TestLEVEL_String(t *testing.T) {
	tests := []struct {
		name string
		r    LEVEL
		want string
	}{
		{name: "TRACE_test_01", r: TRACE, want: "TRACE"},
		{name: "DEBUG_test_01", r: DEBUG, want: "DEBUG"},
		{name: "INFO_test_01", r: INFO, want: "INFO"},
		{name: "WARN_test_01", r: WARN, want: "WARN"},
		{name: "ERROR_test_01", r: ERROR, want: "ERROR"},
		{name: "FATAL_test_01", r: FATAL, want: "FATAL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLEVEL_IsLevelAtLeast(t *testing.T) {
	type args struct {
		level LEVEL
	}
	tests := []struct {
		name string
		r    LEVEL
		args args
		want bool
	}{
		{"01_TRACE_TRACE", TRACE, args{TRACE}, true},
		{"02_TRACE_DEBUG", TRACE, args{DEBUG}, true},
		{"03_TRACE_INFO", TRACE, args{INFO}, true},
		{"04_TRACE_WARN", TRACE, args{WARN}, true},
		{"05_TRACE_ERROR", TRACE, args{ERROR}, true},
		{"06_TRACE_FATAL", TRACE, args{FATAL}, true},

		{"07_DEBUG_TRACE", DEBUG, args{TRACE}, false},
		{"08_DEBUG_DEBUG", DEBUG, args{DEBUG}, true},
		{"09_DEBUG_INFO", DEBUG, args{INFO}, true},
		{"10_DEBUG_WARN", DEBUG, args{WARN}, true},
		{"11_DEBUG_ERROR", DEBUG, args{ERROR}, true},
		{"12_DEBUG_FATAL", DEBUG, args{FATAL}, true},

		{"13_INFO_TRACE", INFO, args{TRACE}, false},
		{"14_INFO_DEBUG", INFO, args{DEBUG}, false},
		{"15_INFO_INFO", INFO, args{INFO}, true},
		{"16_INFO_WARN", INFO, args{WARN}, true},
		{"17_INFO_ERROR", INFO, args{ERROR}, true},
		{"18_INFO_FATAL", INFO, args{FATAL}, true},

		{"19_WARN_TRACE", WARN, args{TRACE}, false},
		{"20_WARN_DEBUG", WARN, args{DEBUG}, false},
		{"21_WARN_INFO", WARN, args{INFO}, false},
		{"22_WARN_WARN", WARN, args{WARN}, true},
		{"23_WARN_ERROR", WARN, args{ERROR}, true},
		{"24_WARN_FATAL", WARN, args{FATAL}, true},

		{"25_ERROR_TRACE", ERROR, args{TRACE}, false},
		{"26_ERROR_DEBUG", ERROR, args{DEBUG}, false},
		{"27_ERROR_INFO", ERROR, args{INFO}, false},
		{"28_ERROR_WARN", ERROR, args{WARN}, false},
		{"29_ERROR_ERROR", ERROR, args{ERROR}, true},
		{"30_ERROR_FATAL", ERROR, args{FATAL}, true},

		{"31_FATAL_TRACE", FATAL, args{TRACE}, false},
		{"32_FATAL_DEBUG", FATAL, args{DEBUG}, false},
		{"33_FATAL_INFO", FATAL, args{INFO}, false},
		{"34_FATAL_WARN", FATAL, args{WARN}, false},
		{"35_FATAL_ERROR", FATAL, args{ERROR}, false},
		{"36_FATAL_FATAL", FATAL, args{FATAL}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsLevelAtLeast(tt.args.level); got != tt.want {
				t.Errorf("IsLevelAtLeast() = %v, want %v", got, tt.want)
			}
		})
	}
}
