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

type LEVEL int

func (r LEVEL) String() string {
	return levelMap[r]
}

var levelMap = []string{
	FATAL: "FATAL",
	ERROR: "ERROR",
	WARN:  "WARN",
	INFO:  "INFO",
	DEBUG: "DEBUG",
	TRACE: "TRACE",
}

const (
	FATAL LEVEL = iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)
