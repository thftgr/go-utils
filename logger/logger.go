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

type LEVEL int

func (r LEVEL) String() string {
	return levelMap[r]
}

var levelMap = map[LEVEL]string{
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
