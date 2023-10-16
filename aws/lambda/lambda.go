package lambdautils

import (
	"context"
	"encoding/json"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"log"
	"os"
	"strings"
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
		lambdaClient = lambda.New(lambda.Options{Credentials: cfg.Credentials, Region: cfg.Region, RetryMaxAttempts: 5, RetryMode: aws.RetryModeAdaptive})
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
	var functions []types.FunctionConfiguration
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

type ActivateLambdasProps struct {
	Prefix   []string
	Suffix   []string
	Contains []string
	Ignore   []string
}

func ActivateLambdas(conf ActivateLambdasProps) {
	lambdas, err := ListLambdas()
	if err != nil {
		panic(err)
	}

	for _, l := range lambdas {
		var b bool

		for _, v := range conf.Ignore {
			if strings.Contains(*l.FunctionName, v) {
				log.Println("Skipping function", *l.FunctionName)
				b = true
				break
			}
		}

		if b {
			continue
		}
		for i := range conf.Prefix {
			if !stringutils.IsBlank(conf.Prefix[i]) && strings.HasPrefix(*l.FunctionName, conf.Prefix[i]) {
				b = true
				break
			}
		}
		if !b {
			for i := range conf.Suffix {
				if !stringutils.IsBlank(conf.Suffix[i]) && strings.HasSuffix(*l.FunctionName, conf.Suffix[i]) {
					b = true
					break
				}
			}
		}
		if !b {
			for i := range conf.Contains {
				if !stringutils.IsBlank(conf.Contains[i]) && strings.Contains(*l.FunctionName, conf.Contains[i]) {
					b = true
					break
				}
			}
		}

		if b {
			cfg, err := lambdaClient.GetFunctionConfiguration(context.Background(), &lambda.GetFunctionConfigurationInput{
				FunctionName: l.FunctionName,
			})
			if err != nil {
				panic(err)
			}
			if cfg.State == types.StateInactive || cfg.StateReasonCode == types.StateReasonCodeIdle {
				marshal, err := json.Marshal(l)
				if err != nil {
					panic(err)
				}
				var c lambda.UpdateFunctionConfigurationInput
				if err = json.Unmarshal(marshal, &c); err != nil {
					panic(err)
				}
				if _, err = lambdaClient.UpdateFunctionConfiguration(context.Background(), &c); err != nil {
					panic(err)
				}
				log.Printf("%s - %s\n", *l.FunctionName, l.State)
			}
		}
	}
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
