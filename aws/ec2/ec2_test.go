package ec2

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"reflect"
	"testing"
)

func TestDescribeInstanceByName(t *testing.T) {
	type args struct {
		instanceName string
	}
	tests := []struct {
		name    string
		args    args
		want    *ec2.DescribeInstancesOutput
		wantErr bool
	}{
		{
			name:    "",
			args:    args{
				instanceName: "qa-bbl-server",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DescribeInstanceByName(tt.args.instanceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeInstanceByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DescribeInstanceByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}