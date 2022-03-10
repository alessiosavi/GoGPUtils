package lambdautils

import (
	S3utils "github.com/alessiosavi/GoGPUtils/aws/S3"
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"strings"
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

func TestDeployAllLambda(t *testing.T) {
	env := "prod"
	objects, err := S3utils.ListBucketObjects(env+"-lambda-asset", "go-")
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(objects))

	lambdas, err := ListLambdas()
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(lambdas))

	for _, object := range objects {
		lambdaName := strings.TrimSuffix(object, ".zip")
		for _, lambda := range lambdas {
			if lambda == env+"-"+lambdaName || lambda == lambdaName+"-"+env {
				log.Println("Uploading lambda", lambda)
				if err = DeployLambdaFromS3(lambda, env+"-lambda-asset", object); err != nil {
					panic(err)
				}
			}
		}
	}
}

//func TestDeleteAllLambda(t *testing.T) {
//	lambdas, err := ListLambdas()
//	if err != nil {
//		panic(err)
//	}
//	for _, lambda := range lambdas {
//		if _, err = DeleteLambda(lambda); err != nil {
//			panic(err)
//		}
//
//	}
//}

func TestDescribeLambda(t *testing.T) {
	lambdas, err := ListLambdas()
	if err != nil {
		panic(err)
	}

	var envs = make(map[string]string)
	for _, lambda := range lambdas {
		if strings.Contains(lambda, "prod") {
			d, err := DescribeLambda(lambda)
			if err != nil {
				panic(err)
			}

			if d.Configuration.Environment != nil {
				for k, v := range d.Configuration.Environment.Variables {
					envs[k] = v
				}
			}
		}
	}

	t.Log(helper.MarshalIndent(envs))
}
