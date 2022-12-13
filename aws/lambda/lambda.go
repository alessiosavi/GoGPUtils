package lambdautils

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"os"
	"sync"
)

var lambdaClient *lambda.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		lambdaClient = lambda.New(lambda.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}
func InvokeLambda(name string, payload []byte, invocationType types.InvocationType) (*lambda.InvokeOutput, error) {
	return lambdaClient.Invoke(context.Background(), &lambda.InvokeInput{
		FunctionName:   aws.String(name),
		InvocationType: invocationType,
		Payload:        payload})
}

func ListLambdas() ([]types.FunctionConfiguration, error) {
	f, err := lambdaClient.ListFunctions(context.Background(), &lambda.ListFunctionsInput{})

	if err != nil {
		return nil, err
	}
	var functions = make([]types.FunctionConfiguration, len(f.Functions))
	functions = append(functions, f.Functions...)

	continuationToken := f.NextMarker
	for continuationToken != nil {
		f, err = lambdaClient.ListFunctions(context.Background(), &lambda.ListFunctionsInput{Marker: continuationToken})
		if err != nil {
			return nil, err
		}
		continuationToken = f.NextMarker
		functions = append(functions, f.Functions...)
	}
	return functions, nil
}

func ListLambdaNames() ([]string, error) {
	lambdas, err := ListLambdas()
	if err != nil {
		return nil, err
	}
	var lambdaNames = make([]string, len(lambdas))
	for i := range lambdas {
		lambdaNames[i] = *lambdas[i].FunctionName
	}
	return lambdaNames, nil
}

func DeleteLambda(lambdaName string) (*lambda.DeleteFunctionOutput, error) {
	return lambdaClient.DeleteFunction(context.Background(), &lambda.DeleteFunctionInput{FunctionName: aws.String(lambdaName)})
}

func DescribeLambda(lambdaName string) (*lambda.GetFunctionOutput, error) {
	function, err := lambdaClient.GetFunction(context.Background(), &lambda.GetFunctionInput{FunctionName: aws.String(lambdaName)})
	if err != nil {
		return nil, err
	}
	return function, nil
}

func DeployLambdaFromS3(functionName, bucket, key string) error {
	function, err := lambdaClient.GetFunction(context.Background(), &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return err
	}

	if _, err = lambdaClient.UpdateFunctionCode(context.Background(), &lambda.UpdateFunctionCodeInput{
		FunctionName:  aws.String(functionName),
		Architectures: function.Configuration.Architectures,
		DryRun:        false,
		Publish:       true,
		RevisionId:    function.Configuration.RevisionId,
		S3Bucket:      aws.String(bucket),
		S3Key:         aws.String(key),
	}); err != nil {
		return err
	}
	return nil
}

func DeployLambdaFromZIP(functionName, zipPath string) error {
	function, err := lambdaClient.GetFunction(context.Background(), &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return err
	}

	file, err := os.ReadFile(zipPath)
	if err != nil {
		return err
	}
	if _, err = lambdaClient.UpdateFunctionCode(context.Background(), &lambda.UpdateFunctionCodeInput{
		FunctionName:  aws.String(functionName),
		Architectures: function.Configuration.Architectures,
		DryRun:        false,
		Publish:       true,
		RevisionId:    function.Configuration.RevisionId,
		ZipFile:       file,
	}); err != nil {
		return err
	}
	//log.Println(helper.MarshalIndent(code))
	return nil
}

func ListTags(lambdaARN string) (*lambda.ListTagsOutput, error) {
	return lambdaClient.ListTags(context.Background(), &lambda.ListTagsInput{
		Resource: aws.String(lambdaARN),
	})
}
