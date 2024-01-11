package parser

import (
	"encoding/json"
	"os"
)

func ParseJsonFromFilepath(path string, v any) error {
	if b, err := os.ReadFile(path); err != nil {
		return err
	} else {
		return json.Unmarshal(b, v)
	}
}

func ParseJsonFromString(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}

func ParseJsonFromBytes(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
