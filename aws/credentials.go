package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewAwsConfig() *aws.Config {
	if cfg, err := config.LoadDefaultConfig(context.Background()); err != nil {
		panic(err)
	} else {
		return &cfg
	}
}
func NewAwsConfigWUsingProfile(profile string) *aws.Config {
	if cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithSharedConfigProfile(profile),
	); err != nil {
		panic(err)
	} else {
		return &cfg
	}
}
