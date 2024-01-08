package loggeradapter

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/thftgr/go-utils/logger"
	"io"
)

func Logger2EchoLogger(l logger.SkipLogger) echo.Logger {
	return &EchoAdapter{l}
}

type EchoAdapter struct {
	Logger logger.SkipLogger
}

func (e *EchoAdapter) Write(p []byte) (n int, err error) {
	e.Info(string(p))
	return len(p), nil
}

func (e *EchoAdapter) Output() io.Writer { return e }
func (e *EchoAdapter) SetOutput(w io.Writer) {
	panic("Do not use this method.SetOutput must be handled in the Logger constructor.")
}
func (e *EchoAdapter) Prefix() string {
	panic("Do not use this method.Prefix must be handled in the Logger constructor.")
}
func (e *EchoAdapter) SetPrefix(p string) {
	panic("Do not use this method.SetPrefix must be handled in the Logger constructor.")
}
func (e *EchoAdapter) Level() log.Lvl {
	panic("Do not use this method.Level must be handled in the Logger constructor.")
}
func (e *EchoAdapter) SetLevel(v log.Lvl) {
	panic("Do not use this method.SetLevel must be handled in the Logger constructor.")
}
func (e *EchoAdapter) SetHeader(h string) {
	panic("Do not use this method.SetHeader must be handled in the Logger constructor.")
}

func (e *EchoAdapter) Print(a ...interface{}) {
	e.Logger.SInfo(1, a...)
}

func (e *EchoAdapter) Printf(format string, a ...interface{}) {
	e.Logger.SInfof(1, format, a...)
}

func (e *EchoAdapter) Printj(j log.JSON) {
	b, _ := json.Marshal(j)
	e.Logger.SInfo(1, string(b))
}

func (e *EchoAdapter) Debug(a ...interface{}) {
	e.Logger.SDebug(1, a...)
}

func (e *EchoAdapter) Debugf(f string, a ...interface{}) {
	e.Logger.SDebugf(1, f, a...)
}

func (e *EchoAdapter) Debugj(j log.JSON) {
	b, _ := json.Marshal(j)
	e.Logger.SDebug(1, string(b))
}

func (e *EchoAdapter) Info(a ...interface{}) {
	e.Logger.SInfo(1, a...)
}

func (e *EchoAdapter) Infof(f string, a ...interface{}) {
	e.Logger.SInfof(1, f, a...)

}

func (e *EchoAdapter) Infoj(j log.JSON) {
	b, _ := json.Marshal(j)
	e.Logger.SInfo(1, string(b))
}

func (e *EchoAdapter) Warn(a ...interface{}) {
	e.Logger.SWarn(1, a...)
}

func (e *EchoAdapter) Warnf(f string, a ...interface{}) {
	e.Logger.SWarnf(1, f, a...)
}

func (e *EchoAdapter) Warnj(j log.JSON) {
	b, _ := json.Marshal(j)
	e.Logger.SWarn(1, string(b))
}

func (e *EchoAdapter) Error(a ...interface{}) {
	e.Logger.SError(1, a...)
}

func (e *EchoAdapter) Errorf(f string, a ...interface{}) {
	e.Logger.SErrorf(1, f, a...)
}

func (e *EchoAdapter) Errorj(j log.JSON) {
	b, _ := json.Marshal(j)
	e.Logger.SError(1, string(b))
}

func (e *EchoAdapter) Fatal(a ...interface{}) {
	e.Logger.SFatal(1, a...)
}

func (e *EchoAdapter) Fatalj(j log.JSON) {
	b, _ := json.Marshal(j)
	e.Logger.SFatal(1, string(b))
}

func (e *EchoAdapter) Fatalf(f string, a ...interface{}) {
	e.Logger.SFatalf(1, f, a...)
}

func (e *EchoAdapter) Panic(a ...interface{}) {
	e.Logger.SFatal(1, a...)
}

func (e *EchoAdapter) Panicj(j log.JSON) {
	b, _ := json.Marshal(j)
	e.Logger.SFatal(1, string(b))
}

func (e *EchoAdapter) Panicf(f string, a ...interface{}) {
	e.Logger.SFatalf(1, f, a...)
}
