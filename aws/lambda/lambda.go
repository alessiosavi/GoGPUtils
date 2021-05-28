package lambda

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func InvokeLambda(name string, payload []byte, invocationType types.InvocationType) (*lambda.InvokeOutput, error) {
	cfg, err := awsutils.New()
	if err != nil {
		return nil, err
	}
	lambdaClient := lambda.New(lambda.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	return lambdaClient.Invoke(context.Background(), &lambda.InvokeInput{
		FunctionName:   aws.String(name),
		InvocationType: invocationType,
		Payload:        payload})
}
