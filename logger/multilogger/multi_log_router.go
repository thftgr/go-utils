package multilogger

import "github.com/thftgr/go-utils/logger"

type MultiLogRouter interface {
	logger.GroupLogger
}

type MultiLogRouterImpl struct {
	writer []logger.GroupLogger
}

func (l *MultiLogRouterImpl) Fatal(v ...any) { l.SFatal(1, v...) }
func (l *MultiLogRouterImpl) Error(v ...any) { l.SError(1, v...) }
func (l *MultiLogRouterImpl) Warn(v ...any)  { l.SWarn(1, v...) }
func (l *MultiLogRouterImpl) Info(v ...any)  { l.SInfo(1, v...) }
func (l *MultiLogRouterImpl) Debug(v ...any) { l.SDebug(1, v...) }
func (l *MultiLogRouterImpl) Trace(v ...any) { l.STrace(1, v...) }

func (l *MultiLogRouterImpl) Fatalf(f string, a ...any) { l.SFatalf(1, f, a...) }
func (l *MultiLogRouterImpl) Errorf(f string, a ...any) { l.SErrorf(1, f, a...) }
func (l *MultiLogRouterImpl) Warnf(f string, a ...any)  { l.SWarnf(1, f, a...) }
func (l *MultiLogRouterImpl) Infof(f string, a ...any)  { l.SInfof(1, f, a...) }
func (l *MultiLogRouterImpl) Debugf(f string, a ...any) { l.SDebugf(1, f, a...) }
func (l *MultiLogRouterImpl) Tracef(f string, a ...any) { l.STracef(1, f, a...) }

func (l *MultiLogRouterImpl) SFatal(s int, v ...any) {
	for i := range l.writer {
		l.writer[i].SFatal(s+1, v...)
	}
}
func (l *MultiLogRouterImpl) SError(s int, v ...any) {
	for i := range l.writer {
		l.writer[i].SError(s+1, v...)
	}
}
func (l *MultiLogRouterImpl) SWarn(s int, v ...any) {
	for i := range l.writer {
		l.writer[i].SWarn(s+1, v...)
	}
}
func (l *MultiLogRouterImpl) SInfo(s int, v ...any) {
	for i := range l.writer {
		l.writer[i].SInfo(s+1, v...)
	}
}
func (l *MultiLogRouterImpl) SDebug(s int, v ...any) {
	for i := range l.writer {
		l.writer[i].SDebug(s+1, v...)
	}
}
func (l *MultiLogRouterImpl) STrace(s int, v ...any) {
	for i := range l.writer {
		l.writer[i].STrace(s+1, v...)
	}
}

func (l *MultiLogRouterImpl) SFatalf(s int, f string, a ...any) {
	for i := range l.writer {
		l.writer[i].SFatalf(s+1, f, a...)
	}
}
func (l *MultiLogRouterImpl) SErrorf(s int, f string, a ...any) {
	for i := range l.writer {
		l.writer[i].SErrorf(s+1, f, a...)
	}
}
func (l *MultiLogRouterImpl) SWarnf(s int, f string, a ...any) {
	for i := range l.writer {
		l.writer[i].SWarnf(s+1, f, a...)
	}
}
func (l *MultiLogRouterImpl) SInfof(s int, f string, a ...any) {
	for i := range l.writer {
		l.writer[i].SInfof(s+1, f, a...)
	}
}
func (l *MultiLogRouterImpl) SDebugf(s int, f string, a ...any) {
	for i := range l.writer {
		l.writer[i].SDebugf(s+1, f, a...)
	}
}
func (l *MultiLogRouterImpl) STracef(s int, f string, a ...any) {
	for i := range l.writer {
		l.writer[i].STracef(s+1, f, a...)
	}
}
func (l *MultiLogRouterImpl) Flush() {
	for i := range l.writer {
		l.writer[i].Flush()
	}
}
func (l *MultiLogRouterImpl) NewGroup(name string) logger.GroupLogger {
	res := &MultiLogRouterImpl{writer: make([]logger.GroupLogger, len(l.writer))}
	for i := range l.writer {
		res.writer[i] = l.writer[i].NewGroup(name)
	}
	return res
}

//=================================================

func NewMultiLogRouterImpl(writer ...logger.GroupLogger) *MultiLogRouterImpl {
	return &MultiLogRouterImpl{writer: writer}
}
