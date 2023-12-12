package logger

import (
	"bytes"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	protocol "github.com/influxdata/line-protocol"
	"github.com/thftgr/go-utils/utils"
	"os"
	"time"
)

type InfluxLogger interface {
	SkipLogger
}

type InfluxLoggerImpl struct {
	tags        []protocol.Tag
	writer      api.WriteAPI
	Prefix      string
	Level       LEVEL
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

func (l *InfluxLoggerImpl) SFatal(s int, v ...any) { l.print(s+1, FATAL, v...) }
func (l *InfluxLoggerImpl) SError(s int, v ...any) { l.print(s+1, ERROR, v...) }
func (l *InfluxLoggerImpl) SWarn(s int, v ...any)  { l.print(s+1, WARN, v...) }
func (l *InfluxLoggerImpl) SInfo(s int, v ...any)  { l.print(s+1, INFO, v...) }
func (l *InfluxLoggerImpl) SDebug(s int, v ...any) { l.print(s+1, DEBUG, v...) }
func (l *InfluxLoggerImpl) STrace(s int, v ...any) { l.print(s+1, TRACE, v...) }

func (l *InfluxLoggerImpl) SFatalf(s int, f string, a ...any) { l.printf(s+1, FATAL, f, a...) }
func (l *InfluxLoggerImpl) SErrorf(s int, f string, a ...any) { l.printf(s+1, ERROR, f, a...) }
func (l *InfluxLoggerImpl) SWarnf(s int, f string, a ...any)  { l.printf(s+1, WARN, f, a...) }
func (l *InfluxLoggerImpl) SInfof(s int, f string, a ...any)  { l.printf(s+1, INFO, f, a...) }
func (l *InfluxLoggerImpl) SDebugf(s int, f string, a ...any) { l.printf(s+1, DEBUG, f, a...) }
func (l *InfluxLoggerImpl) STracef(s int, f string, a ...any) { l.printf(s+1, TRACE, f, a...) }

func (l *InfluxLoggerImpl) Flush() {
	l.writer.Flush()
}

func (l *InfluxLoggerImpl) print(skip int, lvl LEVEL, v ...any) {
	if l.Level >= lvl {
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
	buf.WriteString(lvl.String())
	buf.WriteString(" | ")
	buf.WriteString(fmt.Sprint(v...))
	buf.WriteString("\n")
	l.post(lvl, buf.String())
}

func (l *InfluxLoggerImpl) printf(skip int, lvl LEVEL, format string, args ...any) {
	l.print(skip+1, lvl, fmt.Sprintf(format, args...))
}

func (l *InfluxLoggerImpl) post(level LEVEL, data string) {
	point := influxdb2.NewPointWithMeasurement("log").SetTime(time.Now())
	for i := range l.tags {
		point.AddTag(l.tags[i].Key, l.tags[i].Value)
	}
	point.AddTag("level", level.String())
	point.AddTag("service_name", l.ServiceName)
	point.AddField("log", data)
	l.writer.WritePoint(point)
}

//=================================================

func NewInfluxLoggerImpl(tags []protocol.Tag, writer api.WriteAPI, level LEVEL, serviceName string) *InfluxLoggerImpl {
	return &InfluxLoggerImpl{
		tags:        tags,
		writer:      writer,
		Level:       level,
		ServiceName: serviceName,
	}
}
func NewInfluxLoggerImpl2(writer api.WriteAPI, level LEVEL, serviceName string) *InfluxLoggerImpl {
	return &InfluxLoggerImpl{
		tags: []protocol.Tag{
			{"hostname", os.Getenv("HOSTNAME")},
		},
		writer:      writer,
		Level:       level,
		ServiceName: serviceName,
	}
}
func NewInfluxLoggerImpl3(writer api.WriteAPI, serviceName string) *InfluxLoggerImpl {
	return &InfluxLoggerImpl{
		tags: []protocol.Tag{
			{"hostname", os.Getenv("HOSTNAME")},
		},
		writer:      writer,
		Level:       INFO,
		ServiceName: serviceName,
	}
}
