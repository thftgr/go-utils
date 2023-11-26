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
	Out     io.Writer
	Err     io.Writer
	Prefix  string
	Ftrace  func(v ...any)
	Fdebug  func(v ...any)
	Finfo   func(v ...any)
	Fwarn   func(v ...any)
	Ferror  func(v ...any)
	Ffatal  func(v ...any)
	Ftracef func(f string, a ...any)
	Fdebugf func(f string, a ...any)
	Finfof  func(f string, a ...any)
	Fwarnf  func(f string, a ...any)
	Ferrorf func(f string, a ...any)
	Ffatalf func(f string, a ...any)
}

func (l *LoggerImpl) Trace(v ...any)                    { l.Ftrace(v...) }
func (l *LoggerImpl) Debug(v ...any)                    { l.Fdebug(v...) }
func (l *LoggerImpl) Info(v ...any)                     { l.Finfo(v...) }
func (l *LoggerImpl) Warn(v ...any)                     { l.Fwarn(v...) }
func (l *LoggerImpl) Error(v ...any)                    { l.Ferror(v...) }
func (l *LoggerImpl) Fatal(v ...any)                    { l.Ffatal(v...) }
func (l *LoggerImpl) Tracef(format string, args ...any) { l.Ftracef(format, args...) }
func (l *LoggerImpl) Debugf(format string, args ...any) { l.Fdebugf(format, args...) }
func (l *LoggerImpl) Infof(format string, args ...any)  { l.Finfof(format, args...) }
func (l *LoggerImpl) Warnf(format string, args ...any)  { l.Fwarnf(format, args...) }
func (l *LoggerImpl) Errorf(format string, args ...any) { l.Ferrorf(format, args...) }
func (l *LoggerImpl) Fatalf(format string, args ...any) { l.Ffatalf(format, args...) }
func (l *LoggerImpl) Flush()                            {}

func (l *LoggerImpl) _Trace(v ...any)                 { l.print(l.Out, 1, "TRACE", fmt.Sprint(v...)) }
func (l *LoggerImpl) _Debug(v ...any)                 { l.print(l.Out, 1, "DEBUG", fmt.Sprint(v...)) }
func (l *LoggerImpl) _Info(v ...any)                  { l.print(l.Out, 1, "INFO", fmt.Sprint(v...)) }
func (l *LoggerImpl) _Warn(v ...any)                  { l.print(l.Out, 1, "WARN", fmt.Sprint(v...)) }
func (l *LoggerImpl) _Error(v ...any)                 { l.print(l.Err, 1, "ERROR", fmt.Sprint(v...)) }
func (l *LoggerImpl) _Fatal(v ...any)                 { l.print(l.Err, 1, "FATAL", fmt.Sprint(v...)) }
func (l *LoggerImpl) _Tracef(format string, a ...any) { l.printf(l.Out, 2, "TRACE", format, a...) }
func (l *LoggerImpl) _Debugf(format string, a ...any) { l.printf(l.Out, 2, "DEBUG", format, a...) }
func (l *LoggerImpl) _Infof(format string, a ...any)  { l.printf(l.Out, 2, "INFO", format, a...) }
func (l *LoggerImpl) _Warnf(format string, a ...any)  { l.printf(l.Out, 2, "WARN", format, a...) }
func (l *LoggerImpl) _Errorf(format string, a ...any) { l.printf(l.Err, 2, "ERROR", format, a...) }
func (l *LoggerImpl) _Fatalf(format string, a ...any) { l.printf(l.Err, 2, "FATAL", format, a...) }

func (l *LoggerImpl) print(w io.Writer, skip int, level string, v string) {
	buf := bytes.Buffer{}
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05.999"))
	if l.Prefix != "" {
		buf.WriteString(" | ")
		buf.WriteString(l.Prefix)
	}
	if skip > -1 {
		buf.WriteString(" | ")
		buf.WriteString(l.getCodeLine(skip))
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
	_, file, line, ok := runtime.Caller(skip)
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
	Out:     out,
	Err:     err,
	Prefix:  "",
	Ftrace:  func(v ...any) {},
	Fdebug:  func(v ...any) {},
	Finfo:   func(v ...any) {},
	Fwarn:   func(v ...any) {},
	Ferror:  func(v ...any) {},
	Ffatal:  func(v ...any) {},
	Ftracef: func(f string, a ...any) {},
	Fdebugf: func(f string, a ...any) {},
	Finfof:  func(f string, a ...any) {},
	Fwarnf:  func(f string, a ...any) {},
	Ferrorf: func(f string, a ...any) {},
	Ffatalf: func(f string, a ...any) {},
}

func NewLogger(level LEVEL) Logger {
	l := defaultLoggerImpl
	switch level {
	case ALL, TRACE:
		l.Ftrace = l._Trace
		l.Ftracef = l._Tracef
		fallthrough
	case DEBUG:
		l.Fdebug = l._Debug
		l.Fdebugf = l._Debugf
		fallthrough
	case INFO:
		l.Finfo = l._Info
		l.Finfof = l._Infof
		fallthrough
	case ERROR:
		l.Ferror = l._Error
		l.Ferrorf = l._Errorf
		fallthrough
	case WARN:
		l.Fwarn = l._Warn
		l.Fwarnf = l._Warnf
		fallthrough
	case FATAL:
		l.Ffatal = l._Fatal
		l.Ffatalf = l._Fatalf
		fallthrough
	case OFF:
	}
	return &l
}

func NewLoggerWithWriter(out, err io.Writer, level LEVEL) Logger {
	l := defaultLoggerImpl
	l.Out = out
	l.Err = err
	switch level {
	case ALL, TRACE:
		l.Ftrace = l._Trace
		l.Ftracef = l._Tracef
	case DEBUG:
		l.Fdebugf = l._Debugf
		l.Fdebug = l._Debug
		fallthrough
	case INFO:
		l.Finfo = l._Info
		l.Finfof = l._Infof
		fallthrough
	case ERROR:
		l.Ferror = l._Error
		l.Ferrorf = l._Errorf
		fallthrough
	case WARN:
		l.Fwarn = l._Warn
		l.Fwarnf = l._Warnf
		fallthrough
	case FATAL:
		l.Ffatal = l._Fatal
		l.Ffatalf = l._Fatalf
		fallthrough
	case OFF:
	}
	return &l
}
