package lambdautils

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"io/ioutil"
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

func ListLambda() ([]string, error) {
	f, err := lambdaClient.ListFunctions(context.Background(), &lambda.ListFunctionsInput{})

	if err != nil {
		return nil, err
	}

	var functions = make([]string, len(f.Functions))
	for i, functionName := range f.Functions {
		functions[i] = *functionName.FunctionName
	}

	continuationToken := f.NextMarker
	for continuationToken != nil {
		f, err = lambdaClient.ListFunctions(context.Background(), &lambda.ListFunctionsInput{Marker: continuationToken})
		if err != nil {
			return nil, err
		}
		continuationToken = f.NextMarker
		for _, functionName := range f.Functions {
			functions = append(functions, *functionName.FunctionName)
		}
	}
	return functions, nil
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
	//log.Println(helper.MarshalIndent(code))
	return nil
}

func DeployLambdaFromZIP(functionName, zipPath string) error {
	function, err := lambdaClient.GetFunction(context.Background(), &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(zipPath)
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
