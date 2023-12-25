package influxRepository

import (
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"reflect"
	"strings"
	"time"
)

var p_time_type = reflect.TypeOf((*time.Time)(nil))
var time_type = reflect.TypeOf((*time.Time)(nil)).Elem()
var p_string_type = reflect.TypeOf((*string)(nil))
var string_type = reflect.TypeOf((*string)(nil)).Elem()

type InfluxEntityTagHelper[E InfluxEntity] struct {
	measurement string
	tagIndex    map[string]*reflect.StructField
	fieldIndex  map[string]*reflect.StructField
	timeIndex   *reflect.StructField
}

// NewInfluxEntityTagHelper 요구사항에 부합하지 않은경우 error 대신 panic 을 발생시킵니다.
// 반복해서 호출하도록 설계하지 말고 엔티티당 1회 호출할수있도록 설계하는것을 권장
func NewInfluxEntityTagHelper[E InfluxEntity]() (r *InfluxEntityTagHelper[E]) {
	r = &InfluxEntityTagHelper[E]{
		tagIndex:   make(map[string]*reflect.StructField),
		fieldIndex: make(map[string]*reflect.StructField),
	}
	rtype := reflect.TypeOf((*E)(nil))
	for rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem() // 대상을 가져옴 = **E => *E => E
	}
	if rtype.Kind() != reflect.Struct {
		panic(fmt.Errorf("type %s is not struct or *struct", rtype.String()))
	}

	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		if !field.IsExported() { // skip private Field
			continue
		}

		stag := strings.Split(strings.TrimSpace(field.Tag.Get("influxdb")), ":")
		for i2 := range stag {
			stag[i2] = strings.TrimSpace(stag[i2])
		}

		switch stag[0] {
		case "", "-": // `influxdb:""`, `influxdb:"-"`
			// ignore case

		case "measurement": // `influxdb:"measurement:${value}"`
			if len(stag) < 2 || len(stag[1]) == 0 { // `influxdb:"measurement"` 형태로 사용한경우
				panic(fmt.Errorf(`failed to parse measurement %s.%s Use the following form "measurement:snake_case"`, rtype.String(), field.Name))
			}
			r.measurement = stag[1]

		case "time": // `influxdb:"time"`,
			if !field.Type.AssignableTo(p_time_type) && !field.Type.AssignableTo(time_type) {
				panic(fmt.Errorf(`time tag cannot use %s.%s non time type "%s" use *time.Time or time.Time`, rtype.String(), field.Name, field.Type.String()))
			}
			if r.timeIndex != nil {
				panic(fmt.Errorf(`duplicate time tag ["%s.%s", "%s.%s"]`, rtype.String(), r.timeIndex.Name, rtype.String(), field.Name))
			}
			r.timeIndex = &field

		case "tag": // `influxdb:"tag:${value}"`
			if !field.Type.AssignableTo(string_type) {
				panic(fmt.Errorf(`tag tag cannot use %s.%s non string type "%s" use string`, rtype.String(), field.Name, field.Type.String()))
			}
			if len(stag) < 2 || len(stag[1]) == 0 { // `influxdb:"tag"` 형태로 사용한경우
				panic(fmt.Errorf(`failed to parse tag %s.%s Use the following form "tag:snake_case"`, rtype.String(), field.Name))
			}
			if ti := r.tagIndex[stag[1]]; ti != nil { // 태그명이 중복되는경우
				panic(fmt.Errorf(`duplicate tag name "tag:%s" ["%s.%s", "%s.%s"]`, stag[1], rtype.String(), ti.Name, rtype.String(), field.Name))
			}
			r.tagIndex[stag[1]] = &field

		case "field": // `influxdb:"field:${value}"`
			if len(stag) < 2 || len(stag[1]) == 0 { // `influxdb:"field"` 형태로 사용한경우
				panic(fmt.Errorf(`failed to parse tag %s.%s Use the following form "field:snake_case"`, rtype.String(), field.Name))
			}
			if fi := r.fieldIndex[stag[1]]; fi != nil { // 필드명이 중복되는경우
				panic(fmt.Errorf(`duplicate field name "field:%s" ["%s.%s", "%s.%s"]`, stag[1], rtype.String(), fi.Name, rtype.String(), field.Name))
			}
			r.fieldIndex[stag[1]] = &field
		}

	}
	if r.measurement == "" {
		panic(`entity should have one measurement tag. influxdb:"measurement:${measurement_name}"`)
	}
	if r.timeIndex == nil {
		panic(`entity should have one time field. influxdb:"time"`)
	}
	if r.tagIndex == nil || len(r.tagIndex) < 1 {
		panic(`entity should have one or more tag field. influxdb:"tag:${tag_name}"`)
	}
	if r.fieldIndex == nil || len(r.fieldIndex) < 1 {
		panic(`entity should have one or more value field. influxdb:"field:${field_name}"`)
	}
	return
}

func (r *InfluxEntityTagHelper[E]) ToPoint(e *E) *write.Point {
	rvalue := reflect.ValueOf(e)
	for rvalue.Kind() == reflect.Pointer {
		rvalue = rvalue.Elem()
	}

	point := write.NewPointWithMeasurement(r.measurement)

	// parse time
	if r.timeIndex.Type.Kind() == reflect.Pointer {
		// Time *time.Time
		point.SetTime(*rvalue.FieldByIndex(r.timeIndex.Index).Interface().(*time.Time))
	} else {
		// Time time.Time
		point.SetTime(rvalue.FieldByIndex(r.timeIndex.Index).Interface().(time.Time))
	}

	for k, v := range r.tagIndex {
		point.AddTag(k, rvalue.FieldByIndex(v.Index).String())
	}

	for k, v := range r.fieldIndex {
		point.AddField(k, rvalue.FieldByIndex(v.Index).Interface())
	}
	return point
}

// FromRows 그룹 키가 변경될수 있는데 고려되지 않았음.
func (r *InfluxEntityTagHelper[E]) FromRows(rows *api.QueryTableResult) (res []E, err error) {
	if err = rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		record := rows.Record()
		values := record.Values()

		var row E
		rvalue := reflect.ValueOf(&row).Elem()

		// set time
		rvalue.FieldByIndex(r.timeIndex.Index).Set(reflect.ValueOf(record.Time()))

		for s, v := range values {
			fmt.Printf("key: %+v, value:%+v \n", s, v)
			// set tag fields if valid
			if ti := r.tagIndex[s]; ti != nil {
				if v, ok := v.(string); ok {
					rvalue.FieldByIndex(ti.Index).SetString(v)
				}
			}
			// set value fields if valid
			if fi := r.fieldIndex[s]; fi != nil {
				val := reflect.ValueOf(v)
				fieldVal := rvalue.FieldByIndex(fi.Index)

				if fieldVal.Type().Kind() == reflect.Ptr {
					// Create new pointer to value
					newVal := reflect.New(reflect.TypeOf(v))
					newVal.Elem().Set(val)
					val = newVal
				}

				switch fi.Type.Kind() {
				//case reflect.Invalid:
				//case reflect.Bool:

				case reflect.Int:
					fieldVal.Set(reflect.ValueOf(int(val.Int())))
				case reflect.Int8:
					fieldVal.Set(reflect.ValueOf(int8(val.Int())))
				case reflect.Int16:
					fieldVal.Set(reflect.ValueOf(int16(val.Int())))
				case reflect.Int32:
					fieldVal.Set(reflect.ValueOf(int32(val.Int())))

					//case reflect.Int64:

				case reflect.Uint:
					fieldVal.Set(reflect.ValueOf(uint(val.Uint())))
				case reflect.Uint8:
					fieldVal.Set(reflect.ValueOf(uint8(val.Uint())))
				case reflect.Uint16:
					fieldVal.Set(reflect.ValueOf(uint16(val.Uint())))
				case reflect.Uint32:
					fieldVal.Set(reflect.ValueOf(uint32(val.Uint())))

					//case reflect.Uint64:
					//case reflect.Uintptr:

				case reflect.Float32:
					fieldVal.Set(reflect.ValueOf(float32(val.Float())))

					//case reflect.Float64:

				//case reflect.Complex64:
				//case reflect.Complex128:
				//case reflect.Array:
				//case reflect.Chan:
				//case reflect.Func:
				//case reflect.Interface:
				//case reflect.Map:
				//case reflect.Pointer:
				//case reflect.Slice:
				//case reflect.String:
				//case reflect.Struct:
				//case reflect.UnsafePointer:
				default:
					fieldVal.Set(val)

				}

			}
		}
		res = append(res, row)
	}
	return
}
