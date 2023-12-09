package logger

import (
	"bytes"
	"fmt"
	"github.com/thftgr/go-utils/utils"
	"io"
	"os"
	"time"
)

var (
	out io.Writer = os.Stdout
	err io.Writer = os.Stderr
)

type ConsoleLogger interface {
	Logger
}

type ConsoleLoggerImpl struct {
	// ALL > TRACE > DEBUG > INFO > WARN > ERROR > FATAL > OFF
	// default YYYY-MM-DD HH:mm:ss.sss | ${prefix} | ${file} | ${level} :
	Out    io.Writer
	Err    io.Writer
	Prefix string
	Level  LEVEL
}

func (l *ConsoleLoggerImpl) Fatal(v ...any) {
	if l.Level >= FATAL {
		l.print(l.Err, 1, "FATAL", fmt.Sprint(v...))
	}
}
func (l *ConsoleLoggerImpl) Error(v ...any) {
	if l.Level >= ERROR {
		l.print(l.Err, 1, "ERROR", fmt.Sprint(v...))
	}
}
func (l *ConsoleLoggerImpl) Warn(v ...any) {
	if l.Level >= WARN {
		l.print(l.Out, 1, "WARN", fmt.Sprint(v...))
	}
}
func (l *ConsoleLoggerImpl) Info(v ...any) {
	if l.Level >= INFO {
		l.print(l.Out, 1, "INFO", fmt.Sprint(v...))
	}
}
func (l *ConsoleLoggerImpl) Debug(v ...any) {
	if l.Level >= DEBUG {
		l.print(l.Out, 1, "DEBUG", fmt.Sprint(v...))
	}
}
func (l *ConsoleLoggerImpl) Trace(v ...any) {
	if l.Level >= TRACE {
		l.print(l.Out, 1, "TRACE", fmt.Sprint(v...))
	}
}

func (l *ConsoleLoggerImpl) Fatalf(format string, a ...any) {
	if l.Level >= FATAL {
		l.printf(l.Err, 1, "FATAL", format, a...)
	}
}
func (l *ConsoleLoggerImpl) Errorf(format string, a ...any) {
	if l.Level >= ERROR {
		l.printf(l.Err, 1, "ERROR", format, a...)
	}
}
func (l *ConsoleLoggerImpl) Warnf(format string, a ...any) {
	if l.Level >= WARN {
		l.printf(l.Out, 1, "WARN", format, a...)
	}
}
func (l *ConsoleLoggerImpl) Infof(format string, a ...any) {
	if l.Level >= INFO {
		l.printf(l.Out, 1, "INFO", format, a...)
	}
}
func (l *ConsoleLoggerImpl) Debugf(format string, a ...any) {
	if l.Level >= DEBUG {
		l.printf(l.Out, 1, "DEBUG", format, a...)
	}
}
func (l *ConsoleLoggerImpl) Tracef(format string, a ...any) {
	if l.Level >= TRACE {
		l.printf(l.Out, 1, "TRACE", format, a...)
	}
}

func (l *ConsoleLoggerImpl) Flush() {}

func (l *ConsoleLoggerImpl) print(w io.Writer, skip int, level string, v string) {
	buf := bytes.Buffer{}
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05.999"))
	if l.Prefix != "" {
		buf.WriteString(" | ")
		buf.WriteString(l.Prefix)
	}
	if skip > -1 {
		buf.WriteString(" | ")
		buf.WriteString(utils.GetSourceLine(skip + 1))
	}
	buf.WriteString(" | ")
	buf.WriteString(level)
	buf.WriteString(" | ")
	buf.WriteString(v)
	buf.WriteString("\n")
	_, _ = buf.WriteTo(w)
}

func (l *ConsoleLoggerImpl) printf(w io.Writer, skip int, level, format string, args ...any) {
	l.print(w, skip+1, level, fmt.Sprintf(format, args...))
}

//=================================================

var defaultConsoleLoggerImpl = ConsoleLoggerImpl{
	Out:    out,
	Err:    err,
	Prefix: "",
}

func NewConsoleLogger(level LEVEL) Logger {
	l := defaultConsoleLoggerImpl
	l.Level = level
	return &l
}

func NewConsoleLoggerWithWriter(out, err io.Writer, level LEVEL) Logger {
	l := defaultConsoleLoggerImpl
	l.Out = out
	l.Err = err
	l.Level = level
	return &l
}
