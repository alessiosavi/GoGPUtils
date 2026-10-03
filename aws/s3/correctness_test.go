package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type mockAPI struct {
	API
	copyObject   func(*s3sdk.CopyObjectInput) (*s3sdk.CopyObjectOutput, error)
	deleteObject func(*s3sdk.DeleteObjectInput) (*s3sdk.DeleteObjectOutput, error)
	headObject   func(*s3sdk.HeadObjectInput) (*s3sdk.HeadObjectOutput, error)
	getObject    func(*s3sdk.GetObjectInput) (*s3sdk.GetObjectOutput, error)
	listObjects  func(*s3sdk.ListObjectsV2Input) (*s3sdk.ListObjectsV2Output, error)
}

func (m *mockAPI) CopyObject(_ context.Context, in *s3sdk.CopyObjectInput, _ ...func(*s3sdk.Options)) (*s3sdk.CopyObjectOutput, error) {
	return m.copyObject(in)
}
func (m *mockAPI) DeleteObject(_ context.Context, in *s3sdk.DeleteObjectInput, _ ...func(*s3sdk.Options)) (*s3sdk.DeleteObjectOutput, error) {
	return m.deleteObject(in)
}
func (m *mockAPI) HeadObject(_ context.Context, in *s3sdk.HeadObjectInput, _ ...func(*s3sdk.Options)) (*s3sdk.HeadObjectOutput, error) {
	return m.headObject(in)
}
func (m *mockAPI) GetObject(_ context.Context, in *s3sdk.GetObjectInput, _ ...func(*s3sdk.Options)) (*s3sdk.GetObjectOutput, error) {
	return m.getObject(in)
}
func (m *mockAPI) ListObjectsV2(_ context.Context, in *s3sdk.ListObjectsV2Input, _ ...func(*s3sdk.Options)) (*s3sdk.ListObjectsV2Output, error) {
	return m.listObjects(in)
}

func TestCopyAndMoveObject_PreserveAndEncodeSourceKey(t *testing.T) {
	for _, tt := range []struct{ key, encoded string }{
		{"a/../b", "a/../b"}, {"dir//x", "dir//x"}, {"sp ace?#&+.txt", "sp%20ace%3F%23%26%2B.txt"},
		{"日本/é.txt", "%E6%97%A5%E6%9C%AC/%C3%A9.txt"}, {"/a%2Fb/", "/a%252Fb/"},
	} {
		for _, move := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/move=%t", tt.key, move), func(t *testing.T) {
				copied, deleted := false, false
				client := NewClientWithAPI(&mockAPI{
					copyObject: func(in *s3sdk.CopyObjectInput) (*s3sdk.CopyObjectOutput, error) {
						if got := awssdk.ToString(in.CopySource); got != "source/"+tt.encoded {
							t.Errorf("CopySource = %q, want %q", got, "source/"+tt.encoded)
						}
						copied = true
						return &s3sdk.CopyObjectOutput{}, nil
					},
					deleteObject: func(in *s3sdk.DeleteObjectInput) (*s3sdk.DeleteObjectOutput, error) {
						if !copied || awssdk.ToString(in.Key) != tt.key || awssdk.ToString(in.Bucket) != "source" {
							t.Error("deleted wrong source or deleted before copy")
						}
						deleted = true
						return &s3sdk.DeleteObjectOutput{}, nil
					},
				}, nil, nil)
				operation := client.CopyObject
				if move {
					operation = client.MoveObject
				}
				if err := operation(context.Background(), "source", tt.key, "destination", "new"); err != nil {
					t.Fatal(err)
				}
				if !copied || deleted != move {
					t.Fatal("unexpected copy/delete calls")
				}
			})
		}
	}
}

func TestObjectExists_UsesTypedNotFoundErrors(t *testing.T) {
	for _, tt := range []struct {
		name     string
		err      error
		notFound bool
	}{
		{"key", &types.NoSuchKey{}, true}, {"not found", &types.NotFound{}, true}, {"bucket", &types.NoSuchBucket{}, true},
		{"api code", &smithy.GenericAPIError{Code: "NoSuchKey"}, true},
		{"api bucket", &smithy.GenericAPIError{Code: "NoSuchBucket"}, true},
		{"api not found", &smithy.GenericAPIError{Code: "NotFound"}, true},
		{"http 404", &awshttp.ResponseError{ResponseError: &smithyhttp.ResponseError{Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 404}}, Err: errors.New("empty response")}}, true},
		{"unrelated 404", errors.New("proxy connection failed on port 4040"), false},
		{"unrelated name", errors.New("NotFound in diagnostic text"), false},
		{"access denied", &smithy.GenericAPIError{Code: "AccessDenied", Message: "404 in diagnostic text"}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := fmt.Errorf("request: %w", tt.err)
			client := NewClientWithAPI(&mockAPI{headObject: func(*s3sdk.HeadObjectInput) (*s3sdk.HeadObjectOutput, error) { return nil, wrapped }}, nil, nil)
			exists, err := client.ObjectExists(context.Background(), "bucket", "key")
			if exists || (err == nil) != tt.notFound {
				t.Errorf("exists = %t, err = %v, notFound = %t", exists, err, tt.notFound)
			}
			if !tt.notFound && !errors.Is(err, tt.err) {
				t.Error("original error was lost")
			}
		})
	}
}

func TestListObjectsCallback_LimitsTotalAfterFiltering(t *testing.T) {
	now := time.Now()
	calls, count := 0, 0
	client := NewClientWithAPI(&mockAPI{listObjects: func(in *s3sdk.ListObjectsV2Input) (*s3sdk.ListObjectsV2Output, error) {
		calls++
		if awssdk.ToInt32(in.MaxKeys) != 2 {
			t.Errorf("MaxKeys = %v", in.MaxKeys)
		}
		if calls == 1 {
			return &s3sdk.ListObjectsV2Output{Contents: []types.Object{
				{Key: awssdk.String("old"), LastModified: awssdk.Time(now.Add(-time.Hour))},
				{Key: awssdk.String("one"), LastModified: &now},
			}, IsTruncated: awssdk.Bool(true), NextContinuationToken: awssdk.String("next")}, nil
		}
		if awssdk.ToString(in.ContinuationToken) != "next" {
			t.Error("missing cursor")
		}
		return &s3sdk.ListObjectsV2Output{Contents: []types.Object{{Key: awssdk.String("two"), LastModified: &now}, {Key: awssdk.String("three"), LastModified: &now}}}, nil
	}}, nil, nil)
	err := client.ListObjectsCallback(context.Background(), "bucket", func(ObjectInfo) error { count++; return nil }, WithMaxKeys(2), WithModifiedAfter(now.Add(-time.Minute)))
	if err != nil || calls != 2 || count != 2 {
		t.Fatalf("err = %v, calls = %d, objects = %d", err, calls, count)
	}
}

type mockUploader struct{ input *s3sdk.PutObjectInput }

//lint:ignore SA1019 The mock must implement the existing UploaderAPI; transfermanager migration is deferred.
func (m *mockUploader) Upload(_ context.Context, in *s3sdk.PutObjectInput, _ ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	m.input = in
	return &manager.UploadOutput{}, nil
}

func TestPutObjectReader_ForwardsAllPutOptions(t *testing.T) {
	uploader := &mockUploader{}
	client := NewClientWithAPI(nil, uploader, nil)
	opts := []PutOption{WithContentType("text/plain"), WithStorageClass(types.StorageClassStandardIa), WithServerSideEncryption(types.ServerSideEncryptionAes256), WithMetadata(map[string]string{"author": "test"}), WithCacheControl("max-age=60")}
	if err := client.PutObject(context.Background(), "bucket", "key", []byte("data"), opts...); err != nil {
		t.Fatal(err)
	}
	want := *uploader.input
	want.Body = nil
	if err := client.PutObjectReader(context.Background(), "bucket", "key", strings.NewReader("data"), opts...); err != nil {
		t.Fatal(err)
	}
	got := *uploader.input
	got.Body = nil
	if !reflect.DeepEqual(got, want) {
		t.Errorf("reader input differs from byte upload: got cache control %v, want %v", got.CacheControl, want.CacheControl)
	}
}

func TestGetObjectReader_ReturnsAllMetadata(t *testing.T) {
	now := time.Now()
	body := io.NopCloser(strings.NewReader("data"))
	client := NewClientWithAPI(&mockAPI{getObject: func(*s3sdk.GetObjectInput) (*s3sdk.GetObjectOutput, error) {
		return &s3sdk.GetObjectOutput{Body: body, ContentType: awssdk.String("text/plain"), ContentLength: awssdk.Int64(4), ETag: awssdk.String(`"tag"`), LastModified: &now, StorageClass: types.StorageClassStandardIa, Metadata: map[string]string{"author": "test"}}, nil
	}}, nil, nil)
	r, meta, err := client.GetObjectReader(context.Background(), "bucket", "key")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	want := &ObjectMetadata{ContentType: "text/plain", ContentLength: 4, ETag: "tag", LastModified: now, StorageClass: string(types.StorageClassStandardIa), Metadata: map[string]string{"author": "test"}}
	if r != body || !reflect.DeepEqual(meta, want) {
		t.Fatalf("metadata = %+v, want %+v", meta, want)
	}
}
