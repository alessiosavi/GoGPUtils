package ssm

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"sync"
)

var ssmClient *ssm.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		ssmClient = ssm.New(ssm.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}

func Get(paramName string) (string, error) {
	parameter, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           nil,
		WithDecryption: true,
	})
	if err != nil {
		return "", err
	}
	return *parameter.Parameter.Value, nil
}
