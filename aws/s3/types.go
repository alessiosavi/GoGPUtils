package s3

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ObjectMetadata contains metadata about an S3 object.
type ObjectMetadata struct {
	ContentType   string
	ContentLength int64
	ETag          string
	LastModified  time.Time
	StorageClass  string
	Metadata      map[string]string
}

// ObjectInfo contains information about an S3 object from list operations.
type ObjectInfo struct {
	Key          string
	Size         int64
	ETag         string
	LastModified time.Time
	StorageClass string
}

// DeleteError represents an error deleting an object.
type DeleteError struct {
	Key   string
	Code  string
	Error string
}

// PutOption configures PutObject operations.
type PutOption func(*putOptions)

type putOptions struct {
	contentType          string
	storageClass         types.StorageClass
	serverSideEncryption types.ServerSideEncryption
	metadata             map[string]string
	cacheControl         string
}

// WithContentType sets the content type for the upload.
//
// Example:
//
//	err := client.PutObject(ctx, bucket, key, data, s3.WithContentType("application/json"))
func WithContentType(contentType string) PutOption {
	return func(o *putOptions) {
		o.contentType = contentType
	}
}

// WithStorageClass sets the storage class for the upload.
//
// Example:
//
//	err := client.PutObject(ctx, bucket, key, data, s3.WithStorageClass(types.StorageClassStandardIa))
func WithStorageClass(storageClass types.StorageClass) PutOption {
	return func(o *putOptions) {
		o.storageClass = storageClass
	}
}

// WithServerSideEncryption sets server-side encryption for the upload.
//
// Example:
//
//	err := client.PutObject(ctx, bucket, key, data, s3.WithServerSideEncryption(types.ServerSideEncryptionAes256))
func WithServerSideEncryption(sse types.ServerSideEncryption) PutOption {
	return func(o *putOptions) {
		o.serverSideEncryption = sse
	}
}

// WithMetadata sets custom metadata for the upload.
//
// Example:
//
//	err := client.PutObject(ctx, bucket, key, data, s3.WithMetadata(map[string]string{
//	    "author": "john",
//	    "version": "1.0",
//	}))
func WithMetadata(metadata map[string]string) PutOption {
	return func(o *putOptions) {
		o.metadata = metadata
	}
}

// WithCacheControl sets the cache control header for the upload.
//
// Example:
//
//	err := client.PutObject(ctx, bucket, key, data, s3.WithCacheControl("max-age=3600"))
func WithCacheControl(cacheControl string) PutOption {
	return func(o *putOptions) {
		o.cacheControl = cacheControl
	}
}

// ListOption configures ListObjects operations.
type ListOption func(*listOptions)

type listOptions struct {
	prefix       string
	delimiter    string
	maxKeys      int32
	startAfter   string
	modifiedFrom *time.Time
	modifiedTo   *time.Time
}

// WithPrefix filters objects by key prefix.
//
// Example:
//
//	objects, err := client.ListObjects(ctx, bucket, s3.WithPrefix("logs/2024/"))
func WithPrefix(prefix string) ListOption {
	return func(o *listOptions) {
		o.prefix = prefix
	}
}

// WithDelimiter sets a delimiter for grouping keys.
//
// Example:
//
//	objects, err := client.ListObjects(ctx, bucket, s3.WithDelimiter("/"))
func WithDelimiter(delimiter string) ListOption {
	return func(o *listOptions) {
		o.delimiter = delimiter
	}
}

// WithMaxKeys limits the number of objects returned.
//
// Example:
//
//	objects, err := client.ListObjects(ctx, bucket, s3.WithMaxKeys(100))
func WithMaxKeys(maxKeys int32) ListOption {
	return func(o *listOptions) {
		o.maxKeys = maxKeys
	}
}

// WithStartAfter sets the key to start listing after (for pagination).
//
// Example:
//
//	objects, err := client.ListObjects(ctx, bucket, s3.WithStartAfter(lastKey))
func WithStartAfter(startAfter string) ListOption {
	return func(o *listOptions) {
		o.startAfter = startAfter
	}
}

// WithModifiedAfter filters objects modified after the given time.
//
// Example:
//
//	since := time.Now().Add(-24 * time.Hour)
//	objects, err := client.ListObjects(ctx, bucket, s3.WithModifiedAfter(since))
func WithModifiedAfter(t time.Time) ListOption {
	return func(o *listOptions) {
		o.modifiedFrom = &t
	}
}

// WithModifiedBefore filters objects modified before the given time.
//
// Example:
//
//	until := time.Now().Add(-7 * 24 * time.Hour)
//	objects, err := client.ListObjects(ctx, bucket, s3.WithModifiedBefore(until))
func WithModifiedBefore(t time.Time) ListOption {
	return func(o *listOptions) {
		o.modifiedTo = &t
	}
}
