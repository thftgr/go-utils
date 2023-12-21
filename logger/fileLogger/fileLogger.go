package fileLogger

import (
	"bytes"
	"fmt"
	"github.com/thftgr/go-utils/logger"
	"github.com/thftgr/go-utils/utils"
	"io"
	"time"
)

type RotateFileLogger interface {
	logger.SkipLogger
	io.Closer
	io.Writer
	Rotate() error
	RotateWithWriteCloser(wc io.WriteCloser) error
}

type RotateFileLoggerImpl struct {
	// ALL > TRACE > DEBUG > INFO > WARN > ERROR > FATAL > OFF
	// default YYYY-MM-DD HH:mm:ss.sss | ${prefix} | ${File} | ${level} :
	File      FileLogRotator
	GroupName string
	Level     logger.LEVEL
}

func (l *RotateFileLoggerImpl) Fatal(v ...any) { l.SFatal(1, v...) }
func (l *RotateFileLoggerImpl) Error(v ...any) { l.SError(1, v...) }
func (l *RotateFileLoggerImpl) Warn(v ...any)  { l.SWarn(1, v...) }
func (l *RotateFileLoggerImpl) Info(v ...any)  { l.SInfo(1, v...) }
func (l *RotateFileLoggerImpl) Debug(v ...any) { l.SDebug(1, v...) }
func (l *RotateFileLoggerImpl) Trace(v ...any) { l.STrace(1, v...) }

func (l *RotateFileLoggerImpl) Fatalf(f string, a ...any) { l.SFatalf(1, f, a...) }
func (l *RotateFileLoggerImpl) Errorf(f string, a ...any) { l.SErrorf(1, f, a...) }
func (l *RotateFileLoggerImpl) Warnf(f string, a ...any)  { l.SWarnf(1, f, a...) }
func (l *RotateFileLoggerImpl) Infof(f string, a ...any)  { l.SInfof(1, f, a...) }
func (l *RotateFileLoggerImpl) Debugf(f string, a ...any) { l.SDebugf(1, f, a...) }
func (l *RotateFileLoggerImpl) Tracef(f string, a ...any) { l.STracef(1, f, a...) }

func (l *RotateFileLoggerImpl) SFatal(s int, v ...any) { l.print(s+1, logger.FATAL, v...) }
func (l *RotateFileLoggerImpl) SError(s int, v ...any) { l.print(s+1, logger.ERROR, v...) }
func (l *RotateFileLoggerImpl) SWarn(s int, v ...any)  { l.print(s+1, logger.WARN, v...) }
func (l *RotateFileLoggerImpl) SInfo(s int, v ...any)  { l.print(s+1, logger.INFO, v...) }
func (l *RotateFileLoggerImpl) SDebug(s int, v ...any) { l.print(s+1, logger.DEBUG, v...) }
func (l *RotateFileLoggerImpl) STrace(s int, v ...any) { l.print(s+1, logger.TRACE, v...) }

func (l *RotateFileLoggerImpl) SFatalf(s int, f string, a ...any) {
	l.printf(s+1, logger.FATAL, f, a...)
}
func (l *RotateFileLoggerImpl) SErrorf(s int, f string, a ...any) {
	l.printf(s+1, logger.ERROR, f, a...)
}
func (l *RotateFileLoggerImpl) SWarnf(s int, f string, a ...any) {
	l.printf(s+1, logger.WARN, f, a...)
}
func (l *RotateFileLoggerImpl) SInfof(s int, f string, a ...any) {
	l.printf(s+1, logger.INFO, f, a...)
}
func (l *RotateFileLoggerImpl) SDebugf(s int, f string, a ...any) {
	l.printf(s+1, logger.DEBUG, f, a...)
}
func (l *RotateFileLoggerImpl) STracef(s int, f string, a ...any) {
	l.printf(s+1, logger.TRACE, f, a...)
}
func (l *RotateFileLoggerImpl) Flush() {}

func (l *RotateFileLoggerImpl) Group(name string) (res logger.GroupLogger) {
	return &RotateFileLoggerImpl{
		File:      l.File,
		GroupName: name,
		Level:     l.Level,
	}
}

func (l *RotateFileLoggerImpl) Close() error {
	return l.File.Close()
}

func (l *RotateFileLoggerImpl) Write(p []byte) (int, error) {
	return l.File.Write(p)
}

func (l *RotateFileLoggerImpl) print(skip int, level logger.LEVEL, v ...any) {
	if !l.Level.IsLevelAtLeast(level) {
		return
	}
	buf := bytes.Buffer{}
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05.999"))
	if l.GroupName != "" {
		buf.WriteString(" | ")
		buf.WriteString(l.GroupName)
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
	_, _ = buf.WriteTo(l.File)
}

func (l *RotateFileLoggerImpl) printf(skip int, level logger.LEVEL, format string, args ...any) {
	if !l.Level.IsLevelAtLeast(level) {
		return
	}
	l.print(skip+1, level, fmt.Sprintf(format, args...))
}
