package influxrepository

import (
	protocol "github.com/influxdata/line-protocol"
	"time"
)

// Measurement 태그 지정용으로 사용.
// example
//
//		type Cpu struct{
//		    measurement `influxdb:"measurement:cpu"`
//		    Usage       float `influxdb:"field:usage"`
//	     ...
//		}
type Measurement struct{}

type InfluxEntityEncoder interface { // from Entity
	GetMeasurement() string
	GetTags() []*protocol.Tag
	GetField() []*protocol.Field
	GetTime() time.Time
}

type InfluxEntityDecoder interface { // to Entity
	SetValue(map[string]interface{}) error
	SetTime(time.Time)
}
