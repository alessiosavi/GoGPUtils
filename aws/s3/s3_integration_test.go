// Package s3 integration tests using LocalStack.
//go:build integration

package s3_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/alessiosavi/GoGPUtils/aws/internal/testutil"
	"github.com/alessiosavi/GoGPUtils/aws/s3"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
)

// testContext returns a context with timeout for tests.
func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	return ctx
}

// setupTestBucket creates a unique test bucket and returns cleanup function.
func setupTestBucket(t *testing.T, ctx context.Context, client *s3.Client, cfg *testutil.TestConfig) string {
	t.Helper()

	bucket := testutil.UniqueBucketName()
	s3Client := s3sdk.NewFromConfig(cfg.AWSConfig.AWS())

	_, err := s3Client.CreateBucket(ctx, &s3sdk.CreateBucketInput{
		Bucket: awssdk.String(bucket),
	})
	if err != nil {
		t.Fatalf("failed to create test bucket: %v", err)
	}

	t.Cleanup(func() {
		// Delete all objects first
		output, _ := s3Client.ListObjectsV2(context.Background(), &s3sdk.ListObjectsV2Input{
			Bucket: awssdk.String(bucket),
		})

		if output != nil && len(output.Contents) > 0 {
			for _, obj := range output.Contents {
				s3Client.DeleteObject(context.Background(), &s3sdk.DeleteObjectInput{
					Bucket: awssdk.String(bucket),
					Key:    obj.Key,
				})
			}
		}

		// Delete bucket
		s3Client.DeleteBucket(context.Background(), &s3sdk.DeleteBucketInput{
			Bucket: awssdk.String(bucket),
		})
	})

	return bucket
}

// testConfig holds test configuration.
type testConfig struct {
	cfg    *testutil.TestConfig
	client *s3.Client
	bucket string
}

// setupTest creates a test environment with S3 client and bucket.
func setupTest(t *testing.T) *testConfig {
	t.Helper()
	testutil.SkipIfNoLocalStack(t)

	ctx := testContext(t)

	cfg := testutil.MustLoadConfig(t)

	client, err := s3.NewClient(cfg.AWSConfig)
	if err != nil {
		t.Fatalf("failed to create S3 client: %v", err)
	}

	bucket := setupTestBucket(t, ctx, client, cfg)

	return &testConfig{
		cfg:    cfg,
		client: client,
		bucket: bucket,
	}
}

func TestPutObject(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	tests := []struct {
		name        string
		key         string
		data        []byte
		contentType string
	}{
		{
			name: "simple text file",
			key:  "test.txt",
			data: []byte("Hello, World!"),
		},
		{
			name: "json file",
			key:  "data.json",
			data: []byte(`{"name": "test", "value": 123}`),
		},
		{
			name: "nested path",
			key:  "path/to/nested/file.txt",
			data: []byte("nested content"),
		},
		{
			name: "empty file",
			key:  "empty.txt",
			data: []byte{},
		},
		{
			name: "binary data",
			key:  "binary.bin",
			data: []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tc.client.PutObject(ctx, tc.bucket, tt.key, tt.data)
			if err != nil {
				t.Fatalf("PutObject failed: %v", err)
			}

			// Verify by getting the object
			got, err := tc.client.GetObject(ctx, tc.bucket, tt.key)
			if err != nil {
				t.Fatalf("GetObject failed: %v", err)
			}

			if !bytes.Equal(got, tt.data) {
				t.Errorf("data mismatch: got %v, want %v", got, tt.data)
			}
		})
	}
}

func TestPutObjectWithOptions(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	key := "test-options.json"
	data := []byte(`{"test": true}`)

	err := tc.client.PutObject(ctx, tc.bucket, key, data,
		s3.WithContentType("application/json"),
		s3.WithMetadata(map[string]string{"custom": "value"}),
	)
	if err != nil {
		t.Fatalf("PutObject with options failed: %v", err)
	}

	// Verify metadata
	metadata, err := tc.client.HeadObject(ctx, tc.bucket, key)
	if err != nil {
		t.Fatalf("HeadObject failed: %v", err)
	}

	if metadata.ContentType != "application/json" {
		t.Errorf("content type mismatch: got %s, want application/json", metadata.ContentType)
	}
}

func TestPutObjectReader(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	key := "reader-test.txt"
	content := "content from reader"
	reader := strings.NewReader(content)

	err := tc.client.PutObjectReader(ctx, tc.bucket, key, reader,
		s3.WithContentType("text/plain"),
	)
	if err != nil {
		t.Fatalf("PutObjectReader failed: %v", err)
	}

	// Verify
	got, err := tc.client.GetObject(ctx, tc.bucket, key)
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}

	if string(got) != content {
		t.Errorf("content mismatch: got %s, want %s", string(got), content)
	}
}

func TestGetObject(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: put an object
	key := "get-test.txt"
	data := []byte("test data for get")
	if err := tc.client.PutObject(ctx, tc.bucket, key, data); err != nil {
		t.Fatalf("setup PutObject failed: %v", err)
	}

	t.Run("existing object", func(t *testing.T) {
		got, err := tc.client.GetObject(ctx, tc.bucket, key)
		if err != nil {
			t.Fatalf("GetObject failed: %v", err)
		}

		if !bytes.Equal(got, data) {
			t.Errorf("data mismatch: got %v, want %v", got, data)
		}
	})

	t.Run("non-existent object", func(t *testing.T) {
		_, err := tc.client.GetObject(ctx, tc.bucket, "non-existent-key")
		if !errors.Is(err, s3.ErrObjectNotFound) {
			t.Errorf("expected ErrObjectNotFound, got: %v", err)
		}
	})
}

func TestGetObjectReader(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	key := "reader-get-test.txt"
	data := []byte("streaming test data")
	if err := tc.client.PutObject(ctx, tc.bucket, key, data); err != nil {
		t.Fatalf("setup PutObject failed: %v", err)
	}

	reader, metadata, err := tc.client.GetObjectReader(ctx, tc.bucket, key)
	if err != nil {
		t.Fatalf("GetObjectReader failed: %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	if !bytes.Equal(got, data) {
		t.Errorf("data mismatch: got %v, want %v", got, data)
	}

	if metadata.ContentLength != int64(len(data)) {
		t.Errorf("content length mismatch: got %d, want %d", metadata.ContentLength, len(data))
	}
}

func TestDeleteObject(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	key := "delete-test.txt"
	if err := tc.client.PutObject(ctx, tc.bucket, key, []byte("to be deleted")); err != nil {
		t.Fatalf("setup PutObject failed: %v", err)
	}

	// Verify exists
	exists, err := tc.client.ObjectExists(ctx, tc.bucket, key)
	if err != nil {
		t.Fatalf("ObjectExists failed: %v", err)
	}

	if !exists {
		t.Fatal("object should exist before delete")
	}

	// Delete
	if err := tc.client.DeleteObject(ctx, tc.bucket, key); err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}

	// Verify deleted
	exists, err = tc.client.ObjectExists(ctx, tc.bucket, key)
	if err != nil {
		t.Fatalf("ObjectExists after delete failed: %v", err)
	}

	if exists {
		t.Error("object should not exist after delete")
	}
}

func TestDeleteObjects(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create multiple objects
	keys := []string{"batch/file1.txt", "batch/file2.txt", "batch/file3.txt"}
	for _, key := range keys {
		if err := tc.client.PutObject(ctx, tc.bucket, key, []byte("batch content")); err != nil {
			t.Fatalf("setup PutObject failed for %s: %v", key, err)
		}
	}

	// Delete all
	deleted, errors := tc.client.DeleteObjects(ctx, tc.bucket, keys)
	if len(errors) > 0 {
		t.Fatalf("DeleteObjects had errors: %v", errors)
	}

	if len(deleted) != len(keys) {
		t.Errorf("deleted count mismatch: got %d, want %d", len(deleted), len(keys))
	}

	// Verify all deleted
	for _, key := range keys {
		exists, err := tc.client.ObjectExists(ctx, tc.bucket, key)
		if err != nil {
			t.Fatalf("ObjectExists failed for %s: %v", key, err)
		}

		if exists {
			t.Errorf("object %s should not exist after batch delete", key)
		}
	}
}

func TestCopyObject(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	srcKey := "copy-source.txt"
	dstKey := "copy-dest.txt"
	data := []byte("copy me")

	if err := tc.client.PutObject(ctx, tc.bucket, srcKey, data); err != nil {
		t.Fatalf("setup PutObject failed: %v", err)
	}

	// Copy within same bucket
	if err := tc.client.CopyObject(ctx, tc.bucket, srcKey, tc.bucket, dstKey); err != nil {
		t.Fatalf("CopyObject failed: %v", err)
	}

	// Verify destination
	got, err := tc.client.GetObject(ctx, tc.bucket, dstKey)
	if err != nil {
		t.Fatalf("GetObject for dest failed: %v", err)
	}

	if !bytes.Equal(got, data) {
		t.Errorf("copied data mismatch: got %v, want %v", got, data)
	}

	// Verify source still exists
	exists, err := tc.client.ObjectExists(ctx, tc.bucket, srcKey)
	if err != nil {
		t.Fatalf("ObjectExists for source failed: %v", err)
	}

	if !exists {
		t.Error("source object should still exist after copy")
	}
}

func TestMoveObject(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	srcKey := "move-source.txt"
	dstKey := "move-dest.txt"
	data := []byte("move me")

	if err := tc.client.PutObject(ctx, tc.bucket, srcKey, data); err != nil {
		t.Fatalf("setup PutObject failed: %v", err)
	}

	// Move
	if err := tc.client.MoveObject(ctx, tc.bucket, srcKey, tc.bucket, dstKey); err != nil {
		t.Fatalf("MoveObject failed: %v", err)
	}

	// Verify destination exists with correct data
	got, err := tc.client.GetObject(ctx, tc.bucket, dstKey)
	if err != nil {
		t.Fatalf("GetObject for dest failed: %v", err)
	}

	if !bytes.Equal(got, data) {
		t.Errorf("moved data mismatch: got %v, want %v", got, data)
	}

	// Verify source no longer exists
	exists, err := tc.client.ObjectExists(ctx, tc.bucket, srcKey)
	if err != nil {
		t.Fatalf("ObjectExists for source failed: %v", err)
	}

	if exists {
		t.Error("source object should not exist after move")
	}
}

func TestObjectExists(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	key := "exists-test.txt"

	// Before create
	exists, err := tc.client.ObjectExists(ctx, tc.bucket, key)
	if err != nil {
		t.Fatalf("ObjectExists before create failed: %v", err)
	}

	if exists {
		t.Error("object should not exist before creation")
	}

	// Create
	if err := tc.client.PutObject(ctx, tc.bucket, key, []byte("data")); err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	// After create
	exists, err = tc.client.ObjectExists(ctx, tc.bucket, key)
	if err != nil {
		t.Fatalf("ObjectExists after create failed: %v", err)
	}

	if !exists {
		t.Error("object should exist after creation")
	}
}

func TestHeadObject(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	key := "head-test.txt"
	data := []byte("head object test data")

	if err := tc.client.PutObject(ctx, tc.bucket, key, data,
		s3.WithContentType("text/plain"),
	); err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	metadata, err := tc.client.HeadObject(ctx, tc.bucket, key)
	if err != nil {
		t.Fatalf("HeadObject failed: %v", err)
	}

	if metadata.ContentLength != int64(len(data)) {
		t.Errorf("content length mismatch: got %d, want %d", metadata.ContentLength, len(data))
	}

	if metadata.ETag == "" {
		t.Error("ETag should not be empty")
	}

	if metadata.LastModified.IsZero() {
		t.Error("LastModified should not be zero")
	}
}

func TestHeadObjectNotFound(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	_, err := tc.client.HeadObject(ctx, tc.bucket, "non-existent")
	if !errors.Is(err, s3.ErrObjectNotFound) {
		t.Errorf("expected ErrObjectNotFound, got: %v", err)
	}
}

func TestListObjects(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create multiple objects
	objects := map[string][]byte{
		"list/file1.txt":        []byte("file1"),
		"list/file2.txt":        []byte("file2"),
		"list/subdir/file3.txt": []byte("file3"),
		"other/file4.txt":       []byte("file4"),
	}

	for key, data := range objects {
		if err := tc.client.PutObject(ctx, tc.bucket, key, data); err != nil {
			t.Fatalf("setup PutObject failed for %s: %v", key, err)
		}
	}

	t.Run("list all", func(t *testing.T) {
		result, err := tc.client.ListObjects(ctx, tc.bucket)
		if err != nil {
			t.Fatalf("ListObjects failed: %v", err)
		}

		if len(result) != len(objects) {
			t.Errorf("count mismatch: got %d, want %d", len(result), len(objects))
		}
	})

	t.Run("list with prefix", func(t *testing.T) {
		result, err := tc.client.ListObjects(ctx, tc.bucket, s3.WithPrefix("list/"))
		if err != nil {
			t.Fatalf("ListObjects with prefix failed: %v", err)
		}

		if len(result) != 3 {
			t.Errorf("count mismatch for prefix 'list/': got %d, want 3", len(result))
		}
	})

	t.Run("list with max keys", func(t *testing.T) {
		result, err := tc.client.ListObjects(ctx, tc.bucket, s3.WithMaxKeys(2))
		if err != nil {
			t.Fatalf("ListObjects with max keys failed: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("count mismatch with max keys: got %d, want 2", len(result))
		}
	})
}

func TestListObjectKeys(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	keys := []string{"keys/a.txt", "keys/b.txt", "keys/c.txt"}
	for _, key := range keys {
		if err := tc.client.PutObject(ctx, tc.bucket, key, []byte("data")); err != nil {
			t.Fatalf("setup PutObject failed for %s: %v", key, err)
		}
	}

	result, err := tc.client.ListObjectKeys(ctx, tc.bucket, s3.WithPrefix("keys/"))
	if err != nil {
		t.Fatalf("ListObjectKeys failed: %v", err)
	}

	if len(result) != len(keys) {
		t.Errorf("count mismatch: got %d, want %d", len(result), len(keys))
	}
}

func TestCountObjects(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	for i := 0; i < 5; i++ {
		key := "count/file" + string(rune('0'+i)) + ".txt"
		if err := tc.client.PutObject(ctx, tc.bucket, key, []byte("data")); err != nil {
			t.Fatalf("setup PutObject failed: %v", err)
		}
	}

	count, err := tc.client.CountObjects(ctx, tc.bucket, s3.WithPrefix("count/"))
	if err != nil {
		t.Fatalf("CountObjects failed: %v", err)
	}

	if count != 5 {
		t.Errorf("count mismatch: got %d, want 5", count)
	}
}

func TestTotalSize(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create objects with known sizes
	sizes := []int{100, 200, 300}
	var expectedTotal int64

	for i, size := range sizes {
		key := "size/file" + string(rune('0'+i)) + ".bin"
		data := make([]byte, size)
		expectedTotal += int64(size)

		if err := tc.client.PutObject(ctx, tc.bucket, key, data); err != nil {
			t.Fatalf("setup PutObject failed: %v", err)
		}
	}

	total, err := tc.client.TotalSize(ctx, tc.bucket, s3.WithPrefix("size/"))
	if err != nil {
		t.Fatalf("TotalSize failed: %v", err)
	}

	if total != expectedTotal {
		t.Errorf("total size mismatch: got %d, want %d", total, expectedTotal)
	}
}

func TestListObjectsCallback(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	expected := []string{"callback/a.txt", "callback/b.txt", "callback/c.txt"}
	for _, key := range expected {
		if err := tc.client.PutObject(ctx, tc.bucket, key, []byte("data")); err != nil {
			t.Fatalf("setup PutObject failed: %v", err)
		}
	}

	var found []string
	err := tc.client.ListObjectsCallback(ctx, tc.bucket, func(obj s3.ObjectInfo) error {
		found = append(found, obj.Key)
		return nil
	}, s3.WithPrefix("callback/"))

	if err != nil {
		t.Fatalf("ListObjectsCallback failed: %v", err)
	}

	if len(found) != len(expected) {
		t.Errorf("count mismatch: got %d, want %d", len(found), len(expected))
	}
}

func TestParseS3Path(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantBucket string
		wantKey    string
		wantErr    bool
	}{
		{
			name:       "full path",
			path:       "s3://my-bucket/path/to/file.txt",
			wantBucket: "my-bucket",
			wantKey:    "path/to/file.txt",
		},
		{
			name:       "bucket only",
			path:       "s3://my-bucket",
			wantBucket: "my-bucket",
			wantKey:    "",
		},
		{
			name:       "bucket with trailing slash",
			path:       "s3://my-bucket/",
			wantBucket: "my-bucket",
			wantKey:    "",
		},
		{
			name:       "uppercase prefix",
			path:       "S3://bucket/key",
			wantBucket: "bucket",
			wantKey:    "key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucket, key, err := s3.ParseS3Path(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseS3Path() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if bucket != tt.wantBucket {
				t.Errorf("bucket = %v, want %v", bucket, tt.wantBucket)
			}

			if key != tt.wantKey {
				t.Errorf("key = %v, want %v", key, tt.wantKey)
			}
		})
	}
}

func TestFormatS3Path(t *testing.T) {
	tests := []struct {
		name   string
		bucket string
		key    string
		want   string
	}{
		{
			name:   "full path",
			bucket: "my-bucket",
			key:    "path/to/file.txt",
			want:   "s3://my-bucket/path/to/file.txt",
		},
		{
			name:   "bucket only",
			bucket: "my-bucket",
			key:    "",
			want:   "s3://my-bucket",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s3.FormatS3Path(tt.bucket, tt.key)
			if got != tt.want {
				t.Errorf("FormatS3Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	t.Run("empty bucket", func(t *testing.T) {
		_, err := tc.client.GetObject(ctx, "", "key")
		if err == nil {
			t.Error("expected error for empty bucket")
		}
	})

	t.Run("empty key", func(t *testing.T) {
		_, err := tc.client.GetObject(ctx, tc.bucket, "")
		if err == nil {
			t.Error("expected error for empty key")
		}
	})
}
