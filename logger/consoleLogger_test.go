package logger

import (
	"testing"
)

func TestConsoleLoggerImpl_Trace(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Trace("logmessage")
}
func TestConsoleLoggerImpl_Debug(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Debug("logmessage")
}
func TestConsoleLoggerImpl_Info(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Info("logmessage")
}
func TestConsoleLoggerImpl_Warn(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Warn("logmessage")
}
func TestConsoleLoggerImpl_Error(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Error("logmessage")
}
func TestConsoleLoggerImpl_Fatal(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Fatal("logmessage")
}
func TestConsoleLoggerImpl_Tracef(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Tracef("msg : [%s]", "logmessage")
}
func TestConsoleLoggerImpl_Debugf(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Debugf("msg : [%s]", "logmessage")
}
func TestConsoleLoggerImpl_Infof(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Infof("msg : [%s]", "logmessage")
}
func TestConsoleLoggerImpl_Warnf(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Warnf("msg : [%s]", "logmessage")
}
func TestConsoleLoggerImpl_Errorf(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Errorf("msg : [%s]", "logmessage")
}
func TestConsoleLoggerImpl_Fatalf(t *testing.T) {
	log := NewConsoleLogger(TRACE)
	log.Fatalf("msg : [%s]", "logmessage")
}
