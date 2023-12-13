package logger

import (
	"io"
	"os"
	"testing"
)

func TestConsoleLoggerImpl_default(t *testing.T) {
	type fields struct {
		Out    io.Writer
		Err    io.Writer
		Prefix string
		Level  LEVEL
	}
	type args struct {
		v []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TEST_LEVEL: FATAL", fields{os.Stdout, os.Stderr, "PREFIX", FATAL}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: ERROR", fields{os.Stdout, os.Stderr, "PREFIX", ERROR}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: WARN", fields{os.Stdout, os.Stderr, "PREFIX", WARN}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: INFO", fields{os.Stdout, os.Stderr, "PREFIX", INFO}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: DEBUG", fields{os.Stdout, os.Stderr, "PREFIX", DEBUG}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: TRACE", fields{os.Stdout, os.Stderr, "PREFIX", TRACE}, args{[]any{"HELLO!"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ConsoleLoggerImpl{
				Out:    tt.fields.Out,
				Err:    tt.fields.Err,
				Prefix: tt.fields.Prefix,
				Level:  tt.fields.Level,
			}
			l.Fatal(tt.args.v...)
			l.Error(tt.args.v...)
			l.Warn(tt.args.v...)
			l.Info(tt.args.v...)
			l.Debug(tt.args.v...)
			l.Trace(tt.args.v...)
		})
	}
}

func TestConsoleLoggerImpl_Format(t *testing.T) {
	type fields struct {
		Out    io.Writer
		Err    io.Writer
		Prefix string
		Level  LEVEL
	}
	type args struct {
		farmat string
		v      []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TEST_LEVEL: FATAL", fields{os.Stdout, os.Stderr, "PREFIX", FATAL}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: ERROR", fields{os.Stdout, os.Stderr, "PREFIX", ERROR}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: WARN", fields{os.Stdout, os.Stderr, "PREFIX", WARN}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: INFO", fields{os.Stdout, os.Stderr, "PREFIX", INFO}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: DEBUG", fields{os.Stdout, os.Stderr, "PREFIX", DEBUG}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: TRACE", fields{os.Stdout, os.Stderr, "PREFIX", TRACE}, args{"arg:[%s]", []any{"HELLO!"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ConsoleLoggerImpl{
				Out:    tt.fields.Out,
				Err:    tt.fields.Err,
				Prefix: tt.fields.Prefix,
				Level:  tt.fields.Level,
			}
			l.Fatalf(tt.args.farmat, tt.args.v...)
			l.Errorf(tt.args.farmat, tt.args.v...)
			l.Warnf(tt.args.farmat, tt.args.v...)
			l.Infof(tt.args.farmat, tt.args.v...)
			l.Debugf(tt.args.farmat, tt.args.v...)
			l.Tracef(tt.args.farmat, tt.args.v...)
		})
	}
}

func TestConsoleLoggerImpl_Skip(t *testing.T) {
	type fields struct {
		Out    io.Writer
		Err    io.Writer
		Prefix string
		Level  LEVEL
	}
	type args struct {
		v []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TEST_LEVEL: FATAL", fields{os.Stdout, os.Stderr, "PREFIX", FATAL}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: ERROR", fields{os.Stdout, os.Stderr, "PREFIX", ERROR}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: WARN", fields{os.Stdout, os.Stderr, "PREFIX", WARN}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: INFO", fields{os.Stdout, os.Stderr, "PREFIX", INFO}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: DEBUG", fields{os.Stdout, os.Stderr, "PREFIX", DEBUG}, args{[]any{"HELLO!"}}},
		{"TEST_LEVEL: TRACE", fields{os.Stdout, os.Stderr, "PREFIX", TRACE}, args{[]any{"HELLO!"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ConsoleLoggerImpl{
				Out:    tt.fields.Out,
				Err:    tt.fields.Err,
				Prefix: tt.fields.Prefix,
				Level:  tt.fields.Level,
			}
			l.SFatal(0, tt.args.v...)
			l.SError(0, tt.args.v...)
			l.SWarn(0, tt.args.v...)
			l.SInfo(0, tt.args.v...)
			l.SDebug(0, tt.args.v...)
			l.STrace(0, tt.args.v...)
		})
	}
}

func TestConsoleLoggerImpl_SkipFormat(t *testing.T) {
	type fields struct {
		Out    io.Writer
		Err    io.Writer
		Prefix string
		Level  LEVEL
	}
	type args struct {
		farmat string
		v      []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TEST_LEVEL: FATAL", fields{os.Stdout, os.Stderr, "PREFIX", FATAL}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: ERROR", fields{os.Stdout, os.Stderr, "PREFIX", ERROR}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: WARN", fields{os.Stdout, os.Stderr, "PREFIX", WARN}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: INFO", fields{os.Stdout, os.Stderr, "PREFIX", INFO}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: DEBUG", fields{os.Stdout, os.Stderr, "PREFIX", DEBUG}, args{"arg:[%s]", []any{"HELLO!"}}},
		{"TEST_LEVEL: TRACE", fields{os.Stdout, os.Stderr, "PREFIX", TRACE}, args{"arg:[%s]", []any{"HELLO!"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ConsoleLoggerImpl{
				Out:    tt.fields.Out,
				Err:    tt.fields.Err,
				Prefix: tt.fields.Prefix,
				Level:  tt.fields.Level,
			}
			l.SFatalf(0, tt.args.farmat, tt.args.v...)
			l.SErrorf(0, tt.args.farmat, tt.args.v...)
			l.SWarnf(0, tt.args.farmat, tt.args.v...)
			l.SInfof(0, tt.args.farmat, tt.args.v...)
			l.SDebugf(0, tt.args.farmat, tt.args.v...)
			l.STracef(0, tt.args.farmat, tt.args.v...)
		})
	}
}
