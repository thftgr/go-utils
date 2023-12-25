package gpa

import (
	"testing"
	"time"
)

type ITest interface {
}

type Measurement struct{}

type Test struct {
	Measurement `influxdb:"measurement:test"`
	T1          string    `influxdb:"tag:t1"`
	F2          *string   `influxdb:"field:f2"`
	F3          string    `influxdb:"field:f3"`
	F4          int       `influxdb:"field:f4"`
	Time        time.Time `influxdb:"time"`
	//F6          *time.Time `influxdb:"time"`
}

func TestNewInfluxEntityTagHelper(t *testing.T) {
	//ienv, _ := env.InfluxDB{}.Parse()
	//client := influxdb2.NewClient(ienv.URL, ienv.TOKEN)
	//influxRepository.InfluxRepositoryImpl[Test]{}

	h := NewInfluxEntityTagHelper[Test]()
	ps := "*string-1-1-1-"
	p := h.ToPoint(&Test{
		T1:   "t1t11t1t11",
		F2:   &ps,
		F3:   "ffaa",
		F4:   12857,
		Time: time.Now(),
	})
	t.Logf("%+v", p.TagList()[0])
	t.Logf("%+v", p.FieldList()[0])
	t.Logf("%+v", p.FieldList()[1])
	t.Logf("%+v", p.FieldList()[2])
	t.Logf("%+v", p.Time().Format(time.RFC3339))
}
