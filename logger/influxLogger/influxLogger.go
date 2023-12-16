package influxLogger

import (
	"bytes"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	protocol "github.com/influxdata/line-protocol"
	"github.com/thftgr/go-utils/logger"
	"github.com/thftgr/go-utils/utils"
	"os"
	"time"
)

type InfluxLogger interface {
	logger.SkipLogger
}

type InfluxLoggerImpl struct {
	logger.AbstractSkipLoggerImpl
	tags        []protocol.Tag
	writer      api.WriteAPI
	Prefix      string
	Level       logger.LEVEL
	ServiceName string
}

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
	l.writer.Flush()
}

func (l *InfluxLoggerImpl) print(skip int, level logger.LEVEL, v ...any) {
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
	l.post(level, buf.String())
}

func (l *InfluxLoggerImpl) printf(skip int, lvl logger.LEVEL, format string, args ...any) {
	l.print(skip+1, lvl, fmt.Sprintf(format, args...))
}

func (l *InfluxLoggerImpl) post(level logger.LEVEL, data string) {
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

func NewInfluxLoggerImpl(tags []protocol.Tag, writer api.WriteAPI, level logger.LEVEL, serviceName string) *InfluxLoggerImpl {
	return &InfluxLoggerImpl{
		tags:        tags,
		writer:      writer,
		Level:       level,
		ServiceName: serviceName,
	}
}
func NewInfluxLoggerImpl2(writer api.WriteAPI, level logger.LEVEL, serviceName string) *InfluxLoggerImpl {
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
		Level:       logger.INFO,
		ServiceName: serviceName,
	}
}
