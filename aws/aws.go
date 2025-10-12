package awsutils

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var cfg *aws.Config = nil
var err error = nil
var once sync.Once

func New() (*aws.Config, error) {
	once.Do(func() {
		c, e := config.LoadDefaultConfig(context.Background())
		cfg = &c
		err = e
	})
	return cfg, err
}
