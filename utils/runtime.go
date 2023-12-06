package utils

import (
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
