package parser

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

func ParseTomlFromFilepath(path string, v any) error {
	if b, err := os.ReadFile(path); err != nil {
		return err
	} else {
		return toml.Unmarshal(b, v)
	}
}

func ParseYamlTomlString(data string, v any) error {
	return toml.Unmarshal([]byte(data), v)
}

func ParseYamlTomlBytes(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}
