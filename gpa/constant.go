package gpa

import (
	"github.com/thftgr/go-utils/gpa/influxRepository"
	"reflect"
)

var INFLUX_ENTITY_ENCODER_TYPE = reflect.TypeOf((*influxRepository.InfluxEntityEncoder)(nil)).Elem()
var INFLUX_ENTITY_DECODER_TYPE = reflect.TypeOf((*influxRepository.InfluxEntityDecoder)(nil)).Elem()
