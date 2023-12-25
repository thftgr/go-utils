package utils

import "testing"

type ITest interface {
}

type Measurement struct{}

type Test struct {
	Measurement `influxdb:"measurement:test"`
	f1          string    `influxdb:"tag:f1"`
	f2          [5]string `influxdb:"field:f1"`
}

func TestParseStructTag(t *testing.T) {
	_, err := ParseStructTag[Test]("influxdb")
	if err != nil {
		t.Error(err)
	}
}
