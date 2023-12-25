package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func GetSourceLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "???:0"
	}
	pathParts := strings.Split(file, "/")
	n := len(pathParts)
	if n > 1 {
		file = pathParts[n-2] + "/" + pathParts[n-1]
	}
	return file + ":" + strconv.Itoa(line)
}

func GetFileName(skip int) string {
	_, file, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown"
	}
	pathParts := strings.Split(file, "/")
	n := len(pathParts)
	if n > 1 {
		file = pathParts[n-1]
	}
	return file
}

// =================================================================================
// =================================================================================
// =================================================================================

type StructMapper[E any] struct {
	Type   reflect.Type
	Fields []reflect.StructField
}

// ParseStructTag if E != struct return error
func ParseStructTag[E any](tagKey string) (res *StructMapper[E], err error) {
	rtype := reflect.TypeOf((*E)(nil))
	for rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	if rtype.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s is not struct or *struct", rtype.String())
	}
	res = &StructMapper[E]{Type: rtype, Fields: make([]reflect.StructField, rtype.NumField())}

	for i := 0; i < rtype.NumField(); i++ {
		res.Fields[i] = rtype.Field(i)
	}
	return
}

//func parseTag(str string) *StructMapper {
//	// parse k1:v1;k2:v2;v3 to
//}
