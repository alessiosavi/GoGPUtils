package ssmutils

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
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
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *parameter.Parameter.Value, nil
}

func List() ([]string, error) {
	parameters, err := ssmClient.DescribeParameters(context.Background(), &ssm.DescribeParametersInput{})
	if err != nil {
		return nil, err
	}

	var params []string
	for _, p := range parameters.Parameters {
		params = append(params, *p.Name)
	}
	continuationToken := parameters.NextToken
	for continuationToken != nil {
		parameters, err = ssmClient.DescribeParameters(context.Background(), &ssm.DescribeParametersInput{})
		if err != nil {
			return nil, err
		}
		for _, p := range parameters.Parameters {
			params = append(params, *p.Name)
		}
	}

	return params, nil
}

func Describe(paramName string) (*ssm.GetParameterOutput, error) {
	return ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	})
}
