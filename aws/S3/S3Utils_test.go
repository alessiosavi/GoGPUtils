package S3

import (
	"reflect"
	"testing"
)

// FIXME: create test method that create the S3 bucket for test the methods

func TestListBucketObject(t *testing.T) {
	type args struct {
		bucket string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				bucket: "my-bucket-test-s3",
			},
			want:    []string{"covid.csv", "empty.txt"},
			wantErr: false,
		},
		{
			name: "ko",
			args: args{
				bucket: "this-bucket-does-not-exists",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListBucketObject(tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListBucketObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListBucketObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetObject(t *testing.T) {
	type args struct {
		bucket   string
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{

		{
			name: "ok",
			args: args{
				bucket:   "my-bucket-test-s3",
				fileName: "empty.txt",
			},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetObject(tt.args.bucket, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyObject(t *testing.T) {
	type args struct {
		bucket       string
		bucketTarget string
		key          string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				bucket:       "my-bucket-test-s3",
				bucketTarget: "my-bucket-test-s3-copy",
				key:          "covid.csv",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyObject(tt.args.bucket, tt.args.bucketTarget, tt.args.key, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("CopyObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectExists(t *testing.T) {
	type args struct {
		bucket string
		key    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				bucket: "aws-glue-scripts-796325849317-eu-west-1",
				key:    "FabricaLabAdmin/prova",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ko1",
			args: args{
				bucket: "my-bucket-test-s3",
				key:    "this-file-does-not-exists",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "ko2",
			args: args{
				bucket: "this-bucket-does-not-exists",
				key:    "this-file-does-not-exists",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ObjectExists(tt.args.bucket, tt.args.key)
			if got != tt.want {
				t.Errorf("ObjectExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}
