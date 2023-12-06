package logger

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	protocol "github.com/influxdata/line-protocol"
	"os"
	"time"
)

type InfluxLogger interface {
	Logger
}

type InfluxLoggerImpl struct {
	tags        []protocol.Tag
	writer      api.WriteAPI
	Level       LEVEL
	ServiceName string
}

func (l *InfluxLoggerImpl) Trace(v ...any) {
	if l.Level >= TRACE {
		l.post("TRACE", fmt.Sprint(v...))
	}
}
func (l *InfluxLoggerImpl) Debug(v ...any) {
	if l.Level >= DEBUG {
		l.post("DEBUG", fmt.Sprint(v...))
	}
}
func (l *InfluxLoggerImpl) Info(v ...any) {
	if l.Level >= INFO {
		l.post("INFO", fmt.Sprint(v...))
	}
}
func (l *InfluxLoggerImpl) Warn(v ...any) {
	if l.Level >= WARN {
		l.post("WARN", fmt.Sprint(v...))
	}
}
func (l *InfluxLoggerImpl) Error(v ...any) {
	if l.Level >= ERROR {
		l.post("ERROR", fmt.Sprint(v...))
	}
}
func (l *InfluxLoggerImpl) Fatal(v ...any) {
	if l.Level >= FATAL {
		l.post("FATAL", fmt.Sprint(v...))
	}
}

func (l *InfluxLoggerImpl) Tracef(format string, v ...any) {
	if l.Level >= TRACE {
		l.post("TRACE", fmt.Sprintf(format, v...))
	}
}

func (l *InfluxLoggerImpl) Debugf(format string, v ...any) {
	if l.Level >= DEBUG {
		l.post("DEBUG", fmt.Sprintf(format, v...))
	}
}

func (l *InfluxLoggerImpl) Infof(format string, v ...any) {
	if l.Level >= INFO {
		l.post("INFO", fmt.Sprintf(format, v...))
	}
}

func (l *InfluxLoggerImpl) Warnf(format string, v ...any) {
	if l.Level >= WARN {
		l.post("WARN", fmt.Sprintf(format, v...))
	}
}

func (l *InfluxLoggerImpl) Errorf(format string, v ...any) {
	if l.Level >= ERROR {
		l.post("ERROR", fmt.Sprintf(format, v...))
	}
}

func (l *InfluxLoggerImpl) Fatalf(format string, v ...any) {
	if l.Level >= FATAL {
		l.post("FATAL", fmt.Sprintf(format, v...))
	}
}

func (l *InfluxLoggerImpl) Flush() {
	l.writer.Flush()
}

func (l *InfluxLoggerImpl) post(level string, data string) {
	point := influxdb2.NewPointWithMeasurement("log").SetTime(time.Now())
	for i := range l.tags {
		point.AddTag(l.tags[i].Key, l.tags[i].Value)
	}
	point.AddTag("level", level)
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
