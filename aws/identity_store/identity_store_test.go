package identity_store

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"testing"
)

func TestListUsers(t *testing.T) {
	type args struct {
		identityStore string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "ok",
			args: args{
				identityStore: "d-93677cd81a",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListUsers(tt.args.identityStore)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Println(helper.MarshalIndent(got))
		})
	}
}

func TestDescribeUser(t *testing.T) {
	type args struct {
		userId        string
		identityStore string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				userId:        "93677cd81a-97964c98-38be-46f8-970e-311f7b44629b",
				identityStore: "d-93677cd81a",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DescribeUser(tt.args.userId, tt.args.identityStore)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Println(helper.MarshalIndent(got))
		})
	}
}
