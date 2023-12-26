package influxRepository

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/thftgr/go-utils/utils"
	"math"
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
// 반복해서 호출하도록 설계하지 말고 엔티티당 최초 1회만 호출할수있도록 설계하는것을 권장
// - repo 에서 초기화할때 가지고 있는것을 추천함.
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
		point.SetTime(rvalue.FieldByIndex(r.timeIndex.Index).Elem().Interface().(time.Time))
	} else {
		// Time time.Time
		point.SetTime(rvalue.FieldByIndex(r.timeIndex.Index).Interface().(time.Time))
	}

	for k, v := range r.tagIndex {
		point.AddTag(k, rvalue.FieldByIndex(v.Index).String())
	}

	for k, v := range r.fieldIndex {
		if v.Type.Kind() == reflect.Pointer {
			point.AddField(k, rvalue.FieldByIndex(v.Index).Elem().Interface())
		} else {
			point.AddField(k, rvalue.FieldByIndex(v.Index).Interface())
		}
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
		timeVal := rvalue.FieldByIndex(r.timeIndex.Index)
		if timeVal.Type().Kind() == reflect.Ptr {
			// Create new pointer to value
			newVal := reflect.New(time_type)
			newVal.Elem().Set(reflect.ValueOf(record.Time()))
			rvalue.FieldByIndex(r.timeIndex.Index).Set(newVal)
		} else {
			rvalue.FieldByIndex(r.timeIndex.Index).Set(reflect.ValueOf(record.Time()))
		}

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
				case reflect.Int:
					if v, e := toKindOfInt(reflect.Int, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetInt(v)
					}

				case reflect.Int8:
					if v, e := toKindOfInt(reflect.Int8, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetInt(v)
					}

				case reflect.Int16:
					if v, e := toKindOfInt(reflect.Int16, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetInt(v)
					}
				case reflect.Int32:
					if v, e := toKindOfInt(reflect.Int32, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetInt(v)
					}

				case reflect.Int64:
					if v, e := toKindOfInt(reflect.Int64, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetInt(v)
					}

				case reflect.Uint:
					if v, e := toKindOfUint(reflect.Uint, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetUint(v)
					}
				case reflect.Uint8:
					if v, e := toKindOfUint(reflect.Uint8, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetUint(v)
					}
				case reflect.Uint16:
					if v, e := toKindOfUint(reflect.Uint16, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetUint(v)
					}
				case reflect.Uint32:
					if v, e := toKindOfUint(reflect.Uint32, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetUint(v)
					}

				case reflect.Uint64:
					if v, e := toKindOfUint(reflect.Uint64, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetUint(v)
					}

				case reflect.Float32:
					if v, e := toKindOfFloat(reflect.Float32, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetFloat(v)
					}

				case reflect.Float64:
					if v, e := toKindOfFloat(reflect.Float64, &val); e != nil {
						err = errors.Join(err, e)
					} else {
						fieldVal.SetFloat(v)
					}
				default:
					fieldVal.Set(val)

				}

			}
		}
		res = append(res, row)
	}
	return
}

func toKindOfInt(kind reflect.Kind, value *reflect.Value) (i int64, err error) {
	if utils.In[reflect.Kind](kind, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
		i = value.Int()
		if kind == reflect.Int && math.MinInt <= i && i <= math.MaxInt {
			i = int64(int(i))
		} else if kind == reflect.Int64 {
			i = int64(i)
		} else if kind == reflect.Int32 && math.MinInt32 <= i && i <= math.MaxInt32 {
			i = int64(int32(i))
		} else if kind == reflect.Int16 && math.MinInt16 <= i && i <= math.MaxInt16 {
			i = int64(int16(i))
		} else if kind == reflect.Int8 && math.MinInt8 <= i && i <= math.MaxInt8 {
			i = int64(int8(i))
		} else {
			err = fmt.Errorf("cannot convert %s to %s", value.Kind().String(), kind.String())
		}
	} else {
		err = fmt.Errorf("cannot convert %s to %s", value.Kind().String(), kind.String())
	}
	return
}
func toKindOfUint(kind reflect.Kind, value *reflect.Value) (i uint64, err error) {
	if utils.In[reflect.Kind](kind, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64) {
		i = value.Uint()
		if kind == reflect.Uint && i <= math.MaxUint {
			i = uint64(uint(i))
		} else if kind == reflect.Uint64 {
			i = uint64(i)
		} else if kind == reflect.Uint32 && i <= math.MaxUint32 {
			i = uint64(uint32(i))
		} else if kind == reflect.Uint16 && i <= math.MaxUint16 {
			i = uint64(uint16(i))
		} else if kind == reflect.Uint8 && i <= math.MaxUint8 {
			i = uint64(uint8(i))
		} else {
			err = fmt.Errorf("cannot convert %s to %s", value.Kind().String(), kind.String())
		}
	} else {
		err = fmt.Errorf("cannot convert %s to %s", value.Kind().String(), kind.String())
	}
	return
}

func toKindOfFloat(kind reflect.Kind, value *reflect.Value) (i float64, err error) {
	if utils.In[reflect.Kind](kind, reflect.Float32, reflect.Float64) {
		i = value.Float()
		if kind == reflect.Float32 && math.SmallestNonzeroFloat32 <= i && i <= math.MaxFloat32 {
			i = float64(float32(i))
		} else if kind == reflect.Float64 {
			i = float64(i)
		} else {
			err = fmt.Errorf("cannot convert %s to %s", value.Kind().String(), kind.String())
		}
	} else {
		err = fmt.Errorf("cannot convert %s to %s", value.Kind().String(), kind.String())
	}
	return
}
