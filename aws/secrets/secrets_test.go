package secretsutils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"reflect"
	"testing"
)

//func TestGetSecrets(t *testing.T) {
//	secret, err := GetSecrets("")
//	if err != nil {
//		t.Error(err)
//	}
//	t.Log(secret)
//}

func TestGetSecret(t *testing.T) {
	type args struct {
		secretName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				secretName: "qa-go-am-parser-mail-validator",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSecret(tt.args.secretName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type MailConf struct {
	FromName string   `json:"from_name,omitempty"`
	FromMail string   `json:"from_mail,omitempty"`
	To       string   `json:"to,omitempty"`
	CC       []string `json:"cc,omitempty"`
}

func TestUnmarshalSecret(t *testing.T) {
	type args struct {
		secretName string
		dest       MailConf
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				secretName: "qa-go-am-parser-mail-validator",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnmarshalSecret(tt.args.secretName, &tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(helper.MarshalIndent(tt.args.dest))
		})
	}
}

func TestListSecret(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListSecrets()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}
