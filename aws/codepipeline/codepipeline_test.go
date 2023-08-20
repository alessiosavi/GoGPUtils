package codepipelineutils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	"testing"
)

func TestGetBuildStatus(t *testing.T) {
	type args struct {
		pipelineName string
		max          int
	}
	tests := []struct {
		name    string
		args    args
		want    []types.PipelineExecutionSummary
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "ok",
			args: args{
				pipelineName: "cloud-pipeline",
				max:          10,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBuildStatus(tt.args.pipelineName, tt.args.max)
			t.Log(helper.MarshalIndent(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBuildStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetBuildStatus() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
