package logger

type Logger interface {
	Trace(v ...any)
	Debug(v ...any)
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
	Tracef(format string, v ...any)
	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
	Flush()
}

type SkipLogger interface {
	Logger
	STrace(skip int, v ...any)
	SDebug(skip int, v ...any)
	SInfo(skip int, v ...any)
	SWarn(skip int, v ...any)
	SError(skip int, v ...any)
	SFatal(skip int, v ...any)
	STracef(skip int, format string, v ...any)
	SDebugf(skip int, format string, v ...any)
	SInfof(skip int, format string, v ...any)
	SWarnf(skip int, format string, v ...any)
	SErrorf(skip int, format string, v ...any)
	SFatalf(skip int, format string, v ...any)
}

type AbstractSkipLoggerImpl struct {
	SkipLogger
}

func (l *AbstractSkipLoggerImpl) Fatal(v ...any) { l.SFatal(1, v...) }
func (l *AbstractSkipLoggerImpl) Error(v ...any) { l.SError(1, v...) }
func (l *AbstractSkipLoggerImpl) Warn(v ...any)  { l.SWarn(1, v...) }
func (l *AbstractSkipLoggerImpl) Info(v ...any)  { l.SInfo(1, v...) }
func (l *AbstractSkipLoggerImpl) Debug(v ...any) { l.SDebug(1, v...) }
func (l *AbstractSkipLoggerImpl) Trace(v ...any) { l.STrace(1, v...) }

func (l *AbstractSkipLoggerImpl) Fatalf(f string, a ...any) { l.SFatalf(1, f, a...) }
func (l *AbstractSkipLoggerImpl) Errorf(f string, a ...any) { l.SErrorf(1, f, a...) }
func (l *AbstractSkipLoggerImpl) Warnf(f string, a ...any)  { l.SWarnf(1, f, a...) }
func (l *AbstractSkipLoggerImpl) Infof(f string, a ...any)  { l.SInfof(1, f, a...) }
func (l *AbstractSkipLoggerImpl) Debugf(f string, a ...any) { l.SDebugf(1, f, a...) }
func (l *AbstractSkipLoggerImpl) Tracef(f string, a ...any) { l.STracef(1, f, a...) }

type LEVEL int

func (r LEVEL) String() string {
	return levelMap[r]
}

func (r LEVEL) IsLevelAtLeast(level LEVEL) bool {
	// INFO(2) <= DEBUG(1) -> false
	// INFO(2) <= INFO(2) -> true
	return r <= level
}

var levelMap = []string{
	TRACE: "TRACE",
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

const (
	TRACE LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)
