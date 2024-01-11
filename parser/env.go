package parser

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func ParseEnv(v any) error {
	return envconfig.Process(context.Background(), v)
}

func ParseEnvFromDotENV(v any) error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return envconfig.Process(context.Background(), v)
}
func ParseEnvFromPaths(v any, path ...string) error {
	if err := godotenv.Load(path...); err != nil {
		return err
	}
	return envconfig.Process(context.Background(), v)
}
