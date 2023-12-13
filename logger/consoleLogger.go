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
	SkipLogger
}

type ConsoleLoggerImpl struct {
	// ALL > TRACE > DEBUG > INFO > WARN > ERROR > FATAL > OFF
	// default YYYY-MM-DD HH:mm:ss.sss | ${prefix} | ${file} | ${level} :
	Out    io.Writer
	Err    io.Writer
	Prefix string
	Level  LEVEL
}

func (l *ConsoleLoggerImpl) Fatal(v ...any) { l.SFatal(1, v...) }
func (l *ConsoleLoggerImpl) Error(v ...any) { l.SError(1, v...) }
func (l *ConsoleLoggerImpl) Warn(v ...any)  { l.SWarn(1, v...) }
func (l *ConsoleLoggerImpl) Info(v ...any)  { l.SInfo(1, v...) }
func (l *ConsoleLoggerImpl) Debug(v ...any) { l.SDebug(1, v...) }
func (l *ConsoleLoggerImpl) Trace(v ...any) { l.STrace(1, v...) }

func (l *ConsoleLoggerImpl) Fatalf(f string, a ...any) { l.SFatalf(1, f, a...) }
func (l *ConsoleLoggerImpl) Errorf(f string, a ...any) { l.SErrorf(1, f, a...) }
func (l *ConsoleLoggerImpl) Warnf(f string, a ...any)  { l.SWarnf(1, f, a...) }
func (l *ConsoleLoggerImpl) Infof(f string, a ...any)  { l.SInfof(1, f, a...) }
func (l *ConsoleLoggerImpl) Debugf(f string, a ...any) { l.SDebugf(1, f, a...) }
func (l *ConsoleLoggerImpl) Tracef(f string, a ...any) { l.STracef(1, f, a...) }

func (l *ConsoleLoggerImpl) SFatal(s int, v ...any) { l.print(l.Err, s+1, FATAL, v...) }
func (l *ConsoleLoggerImpl) SError(s int, v ...any) { l.print(l.Err, s+1, ERROR, v...) }
func (l *ConsoleLoggerImpl) SWarn(s int, v ...any)  { l.print(l.Out, s+1, WARN, v...) }
func (l *ConsoleLoggerImpl) SInfo(s int, v ...any)  { l.print(l.Out, s+1, INFO, v...) }
func (l *ConsoleLoggerImpl) SDebug(s int, v ...any) { l.print(l.Out, s+1, DEBUG, v...) }
func (l *ConsoleLoggerImpl) STrace(s int, v ...any) { l.print(l.Out, s+1, TRACE, v...) }

func (l *ConsoleLoggerImpl) SFatalf(s int, f string, a ...any) { l.printf(l.Err, s+1, FATAL, f, a...) }
func (l *ConsoleLoggerImpl) SErrorf(s int, f string, a ...any) { l.printf(l.Err, s+1, ERROR, f, a...) }
func (l *ConsoleLoggerImpl) SWarnf(s int, f string, a ...any)  { l.printf(l.Out, s+1, WARN, f, a...) }
func (l *ConsoleLoggerImpl) SInfof(s int, f string, a ...any)  { l.printf(l.Out, s+1, INFO, f, a...) }
func (l *ConsoleLoggerImpl) SDebugf(s int, f string, a ...any) { l.printf(l.Out, s+1, DEBUG, f, a...) }
func (l *ConsoleLoggerImpl) STracef(s int, f string, a ...any) { l.printf(l.Out, s+1, TRACE, f, a...) }

func (l *ConsoleLoggerImpl) Flush() {}

func (l *ConsoleLoggerImpl) print(w io.Writer, skip int, level LEVEL, v ...any) {
	if !l.Level.IsLevelAtLeast(level) {
		return
	}
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
	buf.WriteString(level.String())
	buf.WriteString(" | ")
	_, _ = fmt.Fprint(&buf, v...)
	buf.WriteString("\n")
	_, _ = buf.WriteTo(w)
}

func (l *ConsoleLoggerImpl) printf(w io.Writer, skip int, lvl LEVEL, format string, args ...any) {
	l.print(w, skip+1, lvl, fmt.Sprintf(format, args...))
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
