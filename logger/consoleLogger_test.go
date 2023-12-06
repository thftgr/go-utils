package logger

import (
	"testing"
)

func TestLoggerImpl_Trace(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Trace("logmessage")
}
func TestLoggerImpl_Debug(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Debug("logmessage")
}
func TestLoggerImpl_Info(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Info("logmessage")
}
func TestLoggerImpl_Warn(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Warn("logmessage")
}
func TestLoggerImpl_Error(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Error("logmessage")
}
func TestLoggerImpl_Fatal(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Fatal("logmessage")
}
func TestLoggerImpl_Tracef(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Tracef("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Debugf(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Debugf("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Infof(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Infof("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Warnf(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Warnf("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Errorf(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Errorf("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Fatalf(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Fatalf("msg : [%s]", "logmessage")
}
