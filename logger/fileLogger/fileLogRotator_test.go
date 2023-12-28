package fileLogger

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimeBaseFileLogRotatorImpl(t *testing.T) {
	log := NewTimeBaseFileLogRotatorImpl("./application")
	fmt.Fprintln(log, time.Now().Format(time.RFC3339))
	time.Sleep(time.Second * 61)
	fmt.Fprintln(log, time.Now().Format(time.RFC3339))
	time.Sleep(time.Second * 61)
	fmt.Fprintln(log, time.Now().Format(time.RFC3339))
	time.Sleep(time.Second * 61)
}
