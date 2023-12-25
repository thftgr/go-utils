package gpa

import (
    "fmt"
    "github.com/influxdata/influxdb-client-go/v2/api"
    "github.com/influxdata/influxdb-client-go/v2/api/write"
    protocol "github.com/influxdata/line-protocol"
    "reflect"
    "strings"
    "time"
)

var (
    p_time_type   = reflect.TypeOf((*time.Time)(nil))
    time_type     = reflect.TypeOf((*time.Time)(nil)).Elem()
    p_string_type = reflect.TypeOf((*string)(nil))
    string_type   = reflect.TypeOf((*string)(nil)).Elem()
)



type InfluxEntityTagHelper[E any] struct {
    measurement string
    tagIndex    map[string]*reflect.StructField
    fieldIndex  map[string]*reflect.StructField
    timeIndex   *reflect.StructField
}


// NewInfluxEntityTagHelper 요구사항에 부합하지 않은경우 error 대신 panic 을 발생시킵니다.
func NewInfluxEntityTagHelper[E any]() (r *InfluxEntityTagHelper[E]) {
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
        //fmt.Printf("name: %+v, tag %+v, public:%+v \n", field.Name, field.Tag.Get("influxdb"), field.IsExported())
        stag := strings.Split(strings.TrimSpace(field.Tag.Get("influxdb")), ":")
        for i2 := range stag {
            stag[i2] = strings.TrimSpace(stag[i2])
        }

        switch stag[0] {
        case "":
            // ignore case

        case "measurement":
            if len(stag) < 2 || len(stag[1]) == 0 { // `influxdb:"tag"` 형태로 사용한경우
                panic(fmt.Errorf(`failed to parse measurement %s.%s Use the following form "measurement:snake_case"`, rtype.String(), field.Name))
            }
            r.measurement = stag[1]

        case "time":
            if !field.Type.AssignableTo(p_time_type) && !field.Type.AssignableTo(time_type) {
                panic(fmt.Errorf(`time tag cannot use %s.%s non time type "%s" use *time.Time or time.Time`, rtype.String(), field.Name, field.Type.String()))
            }
            if r.timeIndex != nil {
                panic(fmt.Errorf(`duplicate time tag ["%s.%s", "%s.%s"]`, rtype.String(), r.timeIndex.Name, rtype.String(), field.Name))
            }
            r.timeIndex = &field

        case "tag":
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

        case "field":
            if len(stag) < 2 || len(stag[1]) == 0 { // `influxdb:"field"` 형태로 사용한경우
                panic(fmt.Errorf(`failed to parse tag %s.%s Use the following form "field:snake_case"`, rtype.String(), field.Name))
            }
            if fi := r.fieldIndex[stag[1]]; fi != nil { // 필드명이 중복되는경우
                panic(fmt.Errorf(`duplicate field name "field:%s" ["%s.%s", "%s.%s"]`, stag[1], rtype.String(), fi.Name, rtype.String(), field.Name))
            }
            r.fieldIndex[stag[1]] = &field
        }

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
        rvalue := reflect.ValueOf(&row)
        rvalue.Set(reflect.)
        for s := range values {
            if ti := r.tagIndex[s]; ti != nil {
                rvalue.FieldByIndex(ti.Index).SetString(values[s].(string))
            }
            if fi := r.fieldIndex[s]; fi != nil {
                rvalue.FieldByIndex(fi.Index).Set(reflect.ValueOf(values[s]))
            }
        }

    }
}
