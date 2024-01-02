package filelogger

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"os"
	"sync"
	"time"
)

type FileLogRotator interface {
	io.WriteCloser
}

type TimeBaseFileLogRotatorImpl struct {
	File      io.WriteCloser
	Scheduler *cron.Cron
	lock      sync.RWMutex
	filename  string
}

// NewTimeBaseFileLogRotatorImpl
// example NewTimeBaseFileLogRotatorImpl()
func NewTimeBaseFileLogRotatorImpl(fileName string) FileLogRotator {
	res := &TimeBaseFileLogRotatorImpl{
		File:      nil,
		Scheduler: cron.New(),
		filename:  fileName,
	}
	res.setFile()

	_, _ = res.Scheduler.AddFunc("0 0 * * *", func() {
		res.setFile()
	})
	res.Scheduler.Start() // async
	return res
}
func (f *TimeBaseFileLogRotatorImpl) setFile() {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.File != nil {
		f.File.Close()
		fmt.Println("file closed.")
	}
	filename := fmt.Sprintf("%s-%s.log", f.filename, time.Now().Local().Format("2006-01-02"))
	if file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		panic(err) // 최초를 제외하고는 사실상 오류가 발생할 가능성이 없음.
	} else {
		f.File = file
		fmt.Println("file rolled.")
	}

}

func (f *TimeBaseFileLogRotatorImpl) Write(p []byte) (n int, err error) {
	return f.File.Write(p)
}

func (f *TimeBaseFileLogRotatorImpl) Close() error {
	f.Scheduler.Stop()
	return f.File.Close()
}
