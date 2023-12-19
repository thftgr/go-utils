package multiLogger

import "github.com/thftgr/go-utils/logger"

type MultiLogRouter interface {
	logger.SkipLogger
}

type MultiLogRouterImpl struct {
	logger.AbstractSkipLoggerImpl
	writer []logger.SkipLogger
}

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

//=================================================

func NewMultiLogRouterImpl(writer ...logger.SkipLogger) *MultiLogRouterImpl {
	return &MultiLogRouterImpl{writer: writer}
}
