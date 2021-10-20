package ec2

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"reflect"
	"testing"
)

func TestListEC2(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListEC2()
			t.Log(helper.MarshalIndent(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("ListEC2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListEC2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// i-0de633b504fc61f84

func TestDescribeHosts(t *testing.T) {
	type args struct {
		instanceName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{instanceName: "i-0de633b504fc61f84"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestDescribeNetwork(t *testing.T) {
	type args struct {
		instance string
	}
	tests := []struct {
		name    string
		args    args
		want    *Network
		wantErr bool
	}{

		{
			name: "",
			args: args{
				instance: "i-0de633b504fc61f84",
			},
			want: &Network{
				PrivateIPv4: "",
				PublicIPv4:  "",
				PrivateDns:  "",
				PublicDns:   "",
				KeyName:     "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DescribeNetwork(tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DescribeNetwork() got = %v, want %v", got, tt.want)
			}
		})
	}
}
