package logger

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/thftgr/go-utils/env"
	"testing"
	"time"
)

var logger = NewConsoleLogger(TRACE)
var writer = func() api.WriteAPI {
	e, err := env.InfluxDB{}.Parse()
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
	client := influxdb2.NewClient(e.URL, e.TOKEN)

	ctx, cancle := context.WithTimeout(context.Background(), time.Second*2)
	defer cancle()

	ok, err := client.Ping(ctx)
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
	if !ok {
		logger.Fatal("failed to ping influxdb")
		panic("failed to ping influxdb")
	}
	logger.Info("success to connected influxdb")
	return client.WriteAPI(e.ORG, e.BUCKET)
}()

func TestInfluxLoggerImpl_Trace(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Trace("logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Debug(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Debug("logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Info(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Info("logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Warn(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Warn("logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Error(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Error("logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Fatal(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Fatal("logmessage")
}
func TestInfluxLoggerImpl_Tracef(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Tracef("msg : [%s]", "logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Debugf(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Debugf("msg : [%s]", "logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Infof(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Infof("msg : [%s]", "logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Warnf(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Warnf("msg : [%s]", "logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Errorf(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Errorf("msg : [%s]", "logmessage")
	l.Flush()
}
func TestInfluxLoggerImpl_Fatalf(t *testing.T) {
	l := NewInfluxLoggerImpl2(writer, TRACE, "testcode")
	l.Fatalf("msg : [%s]", "logmessage")
	l.Flush()
}
