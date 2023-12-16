package fileLogger

import (
	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"math"
)

type FileLogRotator interface {
	Rotate() error
	io.WriteCloser
}

type TimeBaseFileLogRotatorImpl struct {
	File      *lumberjack.Logger
	Scheduler *cron.Cron
}

func NewTimeBaseFileLogRotatorImpl(fileName string) (res *TimeBaseFileLogRotatorImpl) {
	res = &TimeBaseFileLogRotatorImpl{
		File: &lumberjack.Logger{
			Filename:   fileName,
			MaxBackups: math.MaxInt,
			LocalTime:  true,
			Compress:   true,
		},
		Scheduler: cron.New(),
	}
	_, _ = res.Scheduler.AddFunc("0 0 * * *", func() {
		_ = res.Rotate()
	})
	res.Scheduler.Start() // async
	return
}

func (f *TimeBaseFileLogRotatorImpl) Rotate() error {
	return f.File.Rotate()
}

func (f *TimeBaseFileLogRotatorImpl) Write(p []byte) (n int, err error) {
	return f.File.Write(p)
}

func (f *TimeBaseFileLogRotatorImpl) Close() error {
	f.Scheduler.Stop()
	return f.File.Close()
}
