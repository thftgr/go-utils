package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	out io.Writer = os.Stdout
	err io.Writer = os.Stderr
)

type LoggerImpl struct {
	// ALL > TRACE > DEBUG > INFO > WARN > ERROR > FATAL > OFF
	// default YYYY-MM-DD HH:mm:ss.sss | ${prefix} | ${file} | ${level} :
	Out    io.Writer
	Err    io.Writer
	Prefix string
	Level  LEVEL
}

func (l *LoggerImpl) Fatal(v ...any) {
	if l.Level >= FATAL {
		l.print(l.Err, 1, "FATAL", fmt.Sprint(v...))
	}
}
func (l *LoggerImpl) Error(v ...any) {
	if l.Level >= WARN {
		l.print(l.Err, 1, "ERROR", fmt.Sprint(v...))
	}
}
func (l *LoggerImpl) Warn(v ...any) {
	if l.Level >= ERROR {
		l.print(l.Out, 1, "WARN", fmt.Sprint(v...))
	}
}
func (l *LoggerImpl) Info(v ...any) {
	if l.Level >= INFO {
		l.print(l.Out, 1, "INFO", fmt.Sprint(v...))
	}
}
func (l *LoggerImpl) Debug(v ...any) {
	if l.Level >= DEBUG {
		l.print(l.Out, 1, "DEBUG", fmt.Sprint(v...))
	}
}
func (l *LoggerImpl) Trace(v ...any) {
	if l.Level >= TRACE {
		l.print(l.Out, 1, "TRACE", fmt.Sprint(v...))
	}
}

func (l *LoggerImpl) Fatalf(format string, a ...any) {
	if l.Level >= FATAL {
		l.printf(l.Err, 1, "FATAL", format, a...)
	}
}
func (l *LoggerImpl) Errorf(format string, a ...any) {
	if l.Level >= ERROR {
		l.printf(l.Err, 1, "ERROR", format, a...)
	}
}
func (l *LoggerImpl) Warnf(format string, a ...any) {
	if l.Level >= WARN {
		l.printf(l.Out, 1, "WARN", format, a...)
	}
}
func (l *LoggerImpl) Infof(format string, a ...any) {
	if l.Level >= INFO {
		l.printf(l.Out, 1, "INFO", format, a...)
	}
}
func (l *LoggerImpl) Debugf(format string, a ...any) {
	if l.Level >= DEBUG {
		l.printf(l.Out, 1, "DEBUG", format, a...)
	}
}
func (l *LoggerImpl) Tracef(format string, a ...any) {
	if l.Level >= TRACE {
		l.printf(l.Out, 1, "TRACE", format, a...)
	}
}

func (l *LoggerImpl) Flush() {}

func (l *LoggerImpl) print(w io.Writer, skip int, level string, v string) {
	buf := bytes.Buffer{}
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05.999"))
	if l.Prefix != "" {
		buf.WriteString(" | ")
		buf.WriteString(l.Prefix)
	}
	if skip > -1 {
		buf.WriteString(" | ")
		buf.WriteString(l.getCodeLine(skip + 1))
	}
	buf.WriteString(" | ")
	buf.WriteString(level)
	buf.WriteString(" | ")
	buf.WriteString(v)
	buf.WriteString("\n")
	_, _ = buf.WriteTo(w)
}

func (l *LoggerImpl) printf(w io.Writer, skip int, level, format string, args ...any) {
	l.print(w, skip+1, level, fmt.Sprintf(format, args...))
}

func (l *LoggerImpl) getCodeLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "???:0"
	}
	pathParts := strings.Split(file, "/")
	n := len(pathParts)
	if n > 1 {
		file = pathParts[n-2] + "/" + pathParts[n-1]
	}
	return file + ":" + strconv.Itoa(line)
}

//=================================================

var defaultLoggerImpl = LoggerImpl{
	Out:    out,
	Err:    err,
	Prefix: "",
}

func NewLogger(level LEVEL) Logger {
	l := defaultLoggerImpl
	l.Level = level
	return &l
}

func NewLoggerWithWriter(out, err io.Writer, level LEVEL) Logger {
	l := defaultLoggerImpl
	l.Out = out
	l.Err = err
	l.Level = level
	return &l
}
