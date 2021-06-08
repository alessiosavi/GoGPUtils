package glue

import (
	"context"
	"errors"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
)

func StartWorkflow(workflowName string, params map[string]string) error {
	if stringutils.IsBlank(workflowName) {
		return errors.New("workflow is empty")
	}
	cfg, err := awsutils.New()
	if err != nil {
		return nil
	}

	glueClient := glue.New(glue.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	workflow, err := glueClient.GetWorkflow(context.Background(), &glue.GetWorkflowInput{Name: aws.String(workflowName), IncludeGraph: aws.Bool(false)})
	if err != nil {
		return err
	}

	if len(params) > 0 {
		for k, v := range params {
			workflow.Workflow.DefaultRunProperties[k] = v
		}
		if _, err = glueClient.UpdateWorkflow(context.Background(), &glue.UpdateWorkflowInput{Name: aws.String(workflowName), DefaultRunProperties: workflow.Workflow.DefaultRunProperties}); err != nil {
			return err
		}
	}

	_, err = glueClient.StartWorkflowRun(context.Background(), &glue.StartWorkflowRunInput{Name: aws.String(workflowName)})
	return err
}
