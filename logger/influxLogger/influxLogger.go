package influxLogger

import (
	"bytes"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	protocol "github.com/influxdata/line-protocol"
	"github.com/thftgr/go-utils/logger"
	"github.com/thftgr/go-utils/utils"
	"time"
)

type InfluxLogger interface {
	logger.SkipLogger
}

type InfluxLoggerImpl struct {
	Tags        []protocol.Tag
	Writer      api.WriteAPI
	GroupName   string
	Level       logger.LEVEL
	ServiceName string
}

func (l *InfluxLoggerImpl) Fatal(v ...any) { l.SFatal(1, v...) }
func (l *InfluxLoggerImpl) Error(v ...any) { l.SError(1, v...) }
func (l *InfluxLoggerImpl) Warn(v ...any)  { l.SWarn(1, v...) }
func (l *InfluxLoggerImpl) Info(v ...any)  { l.SInfo(1, v...) }
func (l *InfluxLoggerImpl) Debug(v ...any) { l.SDebug(1, v...) }
func (l *InfluxLoggerImpl) Trace(v ...any) { l.STrace(1, v...) }

func (l *InfluxLoggerImpl) Fatalf(f string, a ...any) { l.SFatalf(1, f, a...) }
func (l *InfluxLoggerImpl) Errorf(f string, a ...any) { l.SErrorf(1, f, a...) }
func (l *InfluxLoggerImpl) Warnf(f string, a ...any)  { l.SWarnf(1, f, a...) }
func (l *InfluxLoggerImpl) Infof(f string, a ...any)  { l.SInfof(1, f, a...) }
func (l *InfluxLoggerImpl) Debugf(f string, a ...any) { l.SDebugf(1, f, a...) }
func (l *InfluxLoggerImpl) Tracef(f string, a ...any) { l.STracef(1, f, a...) }

func (l *InfluxLoggerImpl) SFatal(s int, v ...any) { l.print(s+1, logger.FATAL, v...) }
func (l *InfluxLoggerImpl) SError(s int, v ...any) { l.print(s+1, logger.ERROR, v...) }
func (l *InfluxLoggerImpl) SWarn(s int, v ...any)  { l.print(s+1, logger.WARN, v...) }
func (l *InfluxLoggerImpl) SInfo(s int, v ...any)  { l.print(s+1, logger.INFO, v...) }
func (l *InfluxLoggerImpl) SDebug(s int, v ...any) { l.print(s+1, logger.DEBUG, v...) }
func (l *InfluxLoggerImpl) STrace(s int, v ...any) { l.print(s+1, logger.TRACE, v...) }

func (l *InfluxLoggerImpl) SFatalf(s int, f string, a ...any) { l.printf(s+1, logger.FATAL, f, a...) }
func (l *InfluxLoggerImpl) SErrorf(s int, f string, a ...any) { l.printf(s+1, logger.ERROR, f, a...) }
func (l *InfluxLoggerImpl) SWarnf(s int, f string, a ...any)  { l.printf(s+1, logger.WARN, f, a...) }
func (l *InfluxLoggerImpl) SInfof(s int, f string, a ...any)  { l.printf(s+1, logger.INFO, f, a...) }
func (l *InfluxLoggerImpl) SDebugf(s int, f string, a ...any) { l.printf(s+1, logger.DEBUG, f, a...) }
func (l *InfluxLoggerImpl) STracef(s int, f string, a ...any) { l.printf(s+1, logger.TRACE, f, a...) }

func (l *InfluxLoggerImpl) Flush() {
	l.Writer.Flush()
}

func (l *InfluxLoggerImpl) NewGroup(name string) logger.GroupLogger {
	return &InfluxLoggerImpl{
		Tags:        l.Tags,
		Writer:      l.Writer,
		GroupName:   name,
		Level:       l.Level,
		ServiceName: l.ServiceName,
	}
}

func (l *InfluxLoggerImpl) print(skip int, level logger.LEVEL, v ...any) {
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
	l.post(level, buf.String())
}

func (l *InfluxLoggerImpl) printf(skip int, level logger.LEVEL, format string, args ...any) {
	if !l.Level.IsLevelAtLeast(level) {
		return
	}
	l.print(skip+1, level, fmt.Sprintf(format, args...))
}

func (l *InfluxLoggerImpl) post(level logger.LEVEL, data string) {
	point := influxdb2.NewPointWithMeasurement("log").SetTime(time.Now())
	for i := range l.Tags {
		point.AddTag(l.Tags[i].Key, l.Tags[i].Value)
	}
	if l.GroupName != "" {
		point.AddTag("group", l.GroupName)
	}
	point.AddTag("level", level.String())
	point.AddTag("service_name", l.ServiceName)
	point.AddField("log", data)
	l.Writer.WritePoint(point)
}
