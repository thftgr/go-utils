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
	logger.AbstractSkipLoggerImpl
	// ALL > TRACE > DEBUG > INFO > WARN > ERROR > FATAL > OFF
	// default YYYY-MM-DD HH:mm:ss.sss | ${prefix} | ${File} | ${level} :
	File   FileLogRotator
	Prefix string
	Level  logger.LEVEL
}

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
	_, _ = buf.WriteTo(l.File)
}

func (l *RotateFileLoggerImpl) printf(skip int, lvl logger.LEVEL, format string, args ...any) {
	l.print(skip+1, lvl, fmt.Sprintf(format, args...))
}

//=================================================

func NewRotateFileLoggerImpl(file FileLogRotator, lvl logger.LEVEL) *RotateFileLoggerImpl {
	return &RotateFileLoggerImpl{
		File:   file,
		Prefix: "",
		Level:  lvl,
	}
}
func NewRotateFileLoggerImpl1(file FileLogRotator) *RotateFileLoggerImpl {
	return &RotateFileLoggerImpl{
		File:   file,
		Prefix: "",
		Level:  logger.INFO,
	}
}
func NewRotateFileLoggerImpl2() *RotateFileLoggerImpl {
	return &RotateFileLoggerImpl{
		File:   NewTimeBaseFileLogRotatorImpl("./logs/application.log"),
		Prefix: "",
		Level:  logger.INFO,
	}
}
