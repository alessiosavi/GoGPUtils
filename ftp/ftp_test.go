package ftp

import (
	"bytes"
	"github.com/pkg/sftp"
	"reflect"
	"testing"
)

func TestSFTPClient_CreateDirectory(t *testing.T) {
	type fields struct {
		Client *sftp.Client
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SFTPClient{
				Client: tt.fields.Client,
			}
			if err := c.CreateDirectory(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CreateDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSFTPClient_DeleteDirectory(t *testing.T) {
	type fields struct {
		Client *sftp.Client
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SFTPClient{
				Client: tt.fields.Client,
			}
			if err := c.DeleteDirectory(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("DeleteDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSFTPClient_Exist(t *testing.T) {
	type fields struct {
		Client *sftp.Client
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SFTPClient{
				Client: tt.fields.Client,
			}
			got, err := c.Exist(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exist() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSFTPClient_Get(t *testing.T) {
	type fields struct {
		Client *sftp.Client
	}
	type args struct {
		remoteFile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *bytes.Buffer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SFTPClient{
				Client: tt.fields.Client,
			}
			got, err := c.Get(tt.args.remoteFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSFTPClient_IsDir(t *testing.T) {
	type fields struct {
		Client *sftp.Client
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SFTPClient{
				Client: tt.fields.Client,
			}
			got, err := c.IsDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSFTPClient_IsFile(t *testing.T) {
	type fields struct {
		Client *sftp.Client
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SFTPClient{
				Client: tt.fields.Client,
			}
			got, err := c.IsFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSFTPClient_Put(t *testing.T) {
	type fields struct {
		Client *sftp.Client
	}
	type args struct {
		data []byte
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SFTPClient{
				Client: tt.fields.Client,
			}
			if err := c.Put(tt.args.data, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSFTPConf_NewConn(t *testing.T) {
	type fields struct {
		Host     string
		User     string
		Password string
		Port     int
		Timeout  int
	}
	type args struct {
		keyExchanges []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SFTPClient
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				Host:     "test.rebex.net",
				User:     "demo",
				Password: "password",
				Port:     22,
				Timeout:  5,
			},
			args: args{
				keyExchanges: nil,
			},
			want: &SFTPClient{
				Client: &sftp.Client{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SFTPConf{
				Host:     tt.fields.Host,
				User:     tt.fields.User,
				Password: tt.fields.Password,
				Port:     tt.fields.Port,
				Timeout:  tt.fields.Timeout,
			}
			got, err := c.NewConn(tt.args.keyExchanges)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Client.Close()
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewConn() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
