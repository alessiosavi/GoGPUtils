package codepipelineutils

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	"sync"
)

var codepClient *codepipeline.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		codepClient = codepipeline.New(codepipeline.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}

func GetBuildStatus(pipelineName string, max int) ([]types.PipelineExecutionSummary, error) {
	var executions []types.PipelineExecutionSummary
	tmp, err := codepClient.ListPipelineExecutions(context.Background(), &codepipeline.ListPipelineExecutionsInput{PipelineName: aws.String(pipelineName), MaxResults: aws.Int32(int32(max))})
	if err != nil {
		return nil, err
	}
	executions = append(executions, tmp.PipelineExecutionSummaries...)
	if len(executions) == max {
		return executions, nil
	}

	for tmp.NextToken != nil || len(executions) >= max {
		tmp, err = codepClient.ListPipelineExecutions(context.Background(), &codepipeline.ListPipelineExecutionsInput{PipelineName: aws.String(pipelineName), MaxResults: aws.Int32(int32(max - len(executions)))})
		if err != nil {
			return nil, err
		}
		executions = append(executions, tmp.PipelineExecutionSummaries...)
	}

	return executions, err
}
