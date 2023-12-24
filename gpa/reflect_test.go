package gpa

import (
	"reflect"
	"testing"
	"time"
)

type ITest interface {
}

type Measurement struct{}

type Test struct {
	Measurement `influxdb:"measurement:test"`
	T1          string    `influxdb:"tag:t1"`
	T2          *string   `influxdb:"tag:t2"` // 처리 불가능한게 맞음.
	F3          string    `influxdb:"field:f3"`
	F4          int       `influxdb:"field:f4"`
	F5          time.Time `influxdb:"time"`
	//F6          *time.Time `influxdb:"time"`
}

func TestInfluxEntityTagHelper_ParseTagField(t *testing.T) {
	h := &InfluxEntityTagHelper[Test]{
		TagIndex:   make(map[string]*reflect.StructField),
		FieldIndex: make(map[string]*reflect.StructField),
	}
	h.ParseField()
}
