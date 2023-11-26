package logger

import (
	"testing"
)

func TestLoggerImpl_Trace(t *testing.T) {
	log := NewLogger(ALL)
	log.Trace("logmessage")
}
func TestLoggerImpl_Debug(t *testing.T) {
	log := NewLogger(ALL)
	log.Debug("logmessage")
}
func TestLoggerImpl_Info(t *testing.T) {
	log := NewLogger(ALL)
	log.Info("logmessage")
}
func TestLoggerImpl_Warn(t *testing.T) {
	log := NewLogger(ALL)
	log.Warn("logmessage")
}
func TestLoggerImpl_Error(t *testing.T) {
	log := NewLogger(ALL)
	log.Error("logmessage")
}
func TestLoggerImpl_Fatal(t *testing.T) {
	log := NewLogger(ALL)
	log.Fatal("logmessage")
}
func TestLoggerImpl_Tracef(t *testing.T) {
	log := NewLogger(ALL)
	log.Tracef("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Debugf(t *testing.T) {
	log := NewLogger(ALL)
	log.Debugf("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Infof(t *testing.T) {
	log := NewLogger(ALL)
	log.Infof("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Warnf(t *testing.T) {
	log := NewLogger(ALL)
	log.Warnf("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Errorf(t *testing.T) {
	log := NewLogger(ALL)
	log.Errorf("msg : [%s]", "logmessage")
}
func TestLoggerImpl_Fatalf(t *testing.T) {
	log := NewLogger(ALL)
	log.Fatalf("msg : [%s]", "logmessage")
}
