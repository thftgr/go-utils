package consoleLogger

import (
	"bytes"
	"fmt"
	"github.com/thftgr/go-utils/logger"
	"github.com/thftgr/go-utils/utils"
	"io"
	"os"
	"time"
)

type ConsoleLogger interface {
	logger.SkipLogger
}

type ConsoleLoggerImpl struct {
	// ALL > TRACE > DEBUG > INFO > WARN > ERROR > FATAL > OFF
	// default YYYY-MM-DD HH:mm:ss.sss | ${prefix} | ${file} | ${level} :
	Out    io.Writer
	Err    io.Writer
	Prefix string
	Level  logger.LEVEL
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

func (l *ConsoleLoggerImpl) SFatal(s int, v ...any) { l.print(l.Err, s+1, logger.FATAL, v...) }
func (l *ConsoleLoggerImpl) SError(s int, v ...any) { l.print(l.Err, s+1, logger.ERROR, v...) }
func (l *ConsoleLoggerImpl) SWarn(s int, v ...any)  { l.print(l.Out, s+1, logger.WARN, v...) }
func (l *ConsoleLoggerImpl) SInfo(s int, v ...any)  { l.print(l.Out, s+1, logger.INFO, v...) }
func (l *ConsoleLoggerImpl) SDebug(s int, v ...any) { l.print(l.Out, s+1, logger.DEBUG, v...) }
func (l *ConsoleLoggerImpl) STrace(s int, v ...any) { l.print(l.Out, s+1, logger.TRACE, v...) }

func (l *ConsoleLoggerImpl) SFatalf(s int, f string, a ...any) {
	l.printf(l.Err, s+1, logger.FATAL, f, a...)
}
func (l *ConsoleLoggerImpl) SErrorf(s int, f string, a ...any) {
	l.printf(l.Err, s+1, logger.ERROR, f, a...)
}
func (l *ConsoleLoggerImpl) SWarnf(s int, f string, a ...any) {
	l.printf(l.Out, s+1, logger.WARN, f, a...)
}
func (l *ConsoleLoggerImpl) SInfof(s int, f string, a ...any) {
	l.printf(l.Out, s+1, logger.INFO, f, a...)
}
func (l *ConsoleLoggerImpl) SDebugf(s int, f string, a ...any) {
	l.printf(l.Out, s+1, logger.DEBUG, f, a...)
}
func (l *ConsoleLoggerImpl) STracef(s int, f string, a ...any) {
	l.printf(l.Out, s+1, logger.TRACE, f, a...)
}

func (l *ConsoleLoggerImpl) Flush() {}

func (l *ConsoleLoggerImpl) print(w io.Writer, skip int, level logger.LEVEL, v ...any) {
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

func (l *ConsoleLoggerImpl) printf(w io.Writer, skip int, level logger.LEVEL, format string, args ...any) {
	if !l.Level.IsLevelAtLeast(level) {
		return
	}
	l.print(w, skip+1, level, fmt.Sprintf(format, args...))
}

//=================================================

func NewConsoleLoggerImpl(level logger.LEVEL) *ConsoleLoggerImpl {
	return &ConsoleLoggerImpl{
		Out:   os.Stdout,
		Err:   os.Stderr,
		Level: level,
	}
}

func NewConsoleLoggerImplWithWriter(out, err io.Writer, level logger.LEVEL) *ConsoleLoggerImpl {
	return &ConsoleLoggerImpl{
		Out:   out,
		Err:   err,
		Level: level,
	}
}
