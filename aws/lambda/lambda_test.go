package lambdautils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"testing"
)

func TestListLambda(t *testing.T) {
	lambdas, err := ListLambdas()
	if err != nil {
		panic(err)
	}
	t.Log(helper.MarshalIndent(lambdas))
}

func TestDeployLambdaFromS3(t *testing.T) {
	bucket := "qa-lambda-asset"
	key := "go-centric-parser.zip"
	functionName := "qa-go-wac-parser"
	err := DeployLambdaFromS3(functionName, bucket, key)
	if err != nil {
		t.Error(err)
	}
}

func TestDeployLambdaFromZIP(t *testing.T) {
	zipFile := "/opt/Workspace/Go/Lavoro/thom-browne-lambdas/go-wac-parser/go-wac-parser.zip"
	functionName := "qa-go-wac-parser"
	err := DeployLambdaFromZIP(functionName, zipFile)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteAllLambda(t *testing.T) {
	lambdas, err := ListLambdas()
	if err != nil {
		panic(err)
	}
	for _, lambda := range lambdas {
		if _, err = DeleteLambda(lambda); err != nil {
			panic(err)
		}

	}
}
