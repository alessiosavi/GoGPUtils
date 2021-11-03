package lambdautils

import (
	"reflect"
	"testing"
)

func TestListLambda(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{{
		name:    "ok",
		want:    nil,
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListLambda()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListLambda() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListLambda() got = %v, want %v", got, tt.want)
			}
		})
	}
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
