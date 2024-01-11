package parser

import (
	"gopkg.in/yaml.v3"
	"os"
)

func ParseYamlFromFilepath(path string, v any) error {
	if b, err := os.ReadFile(path); err != nil {
		return err
	} else {
		return yaml.Unmarshal(b, v)
	}
}

func ParseYamlFromString(data string, v any) error {
	return yaml.Unmarshal([]byte(data), v)
}

func ParseYamlFromBytes(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
