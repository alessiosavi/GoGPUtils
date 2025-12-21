package s3

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// GetObject downloads an object from S3 and returns its contents.
//
// Example:
//
//	data, err := client.GetObject(ctx, "my-bucket", "path/to/file.txt")
//	if errors.Is(err, s3.ErrObjectNotFound) {
//	    // Handle not found
//	}
func (c *Client) GetObject(ctx context.Context, bucket, key string) ([]byte, error) {
	if err := validateBucketKey(bucket, key); err != nil {
		return nil, err
	}

	output, err := c.api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: awssdk.String(bucket),
		Key:    awssdk.String(key),
	})
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrObjectNotFound
		}

		return nil, aws.WrapError(serviceName, "GetObject", err)
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, aws.WrapError(serviceName, "GetObject", err)
	}

	return data, nil
}

// GetObjectReader returns a reader for streaming object content.
// The caller is responsible for closing the reader.
//
// Example:
//
//	reader, metadata, err := client.GetObjectReader(ctx, "my-bucket", "large-file.zip")
//	if err != nil {
//	    return err
//	}
//	defer reader.Close()
//	io.Copy(dst, reader)
func (c *Client) GetObjectReader(ctx context.Context, bucket, key string) (io.ReadCloser, *ObjectMetadata, error) {
	if err := validateBucketKey(bucket, key); err != nil {
		return nil, nil, err
	}

	output, err := c.api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: awssdk.String(bucket),
		Key:    awssdk.String(key),
	})
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil, ErrObjectNotFound
		}

		return nil, nil, aws.WrapError(serviceName, "GetObject", err)
	}

	metadata := &ObjectMetadata{
		ContentType:   awssdk.ToString(output.ContentType),
		ContentLength: awssdk.ToInt64(output.ContentLength),
		ETag:          strings.Trim(awssdk.ToString(output.ETag), "\""),
		LastModified:  awssdk.ToTime(output.LastModified),
	}

	return output.Body, metadata, nil
}

// PutObject uploads data to S3.
// Content type is automatically detected if not specified.
//
// Example:
//
//	err := client.PutObject(ctx, "my-bucket", "path/to/file.txt", []byte("content"))
//
//	// With options
//	err := client.PutObject(ctx, "my-bucket", "path/to/file.json", data,
//	    s3.WithContentType("application/json"),
//	    s3.WithStorageClass(types.StorageClassStandardIa),
//	)
func (c *Client) PutObject(ctx context.Context, bucket, key string, data []byte, opts ...PutOption) error {
	if err := validateBucketKey(bucket, key); err != nil {
		return err
	}

	options := &putOptions{}
	for _, opt := range opts {
		opt(options)
	}

	contentType := options.contentType
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	input := &s3.PutObjectInput{
		Bucket:      awssdk.String(bucket),
		Key:         awssdk.String(key),
		Body:        bytes.NewReader(data),
		ContentType: awssdk.String(contentType),
	}

	if options.storageClass != "" {
		input.StorageClass = options.storageClass
	}

	if options.serverSideEncryption != "" {
		input.ServerSideEncryption = options.serverSideEncryption
	}

	if options.metadata != nil {
		input.Metadata = options.metadata
	}

	if options.cacheControl != "" {
		input.CacheControl = awssdk.String(options.cacheControl)
	}

	_, err := c.uploader.Upload(ctx, input)
	if err != nil {
		return aws.WrapError(serviceName, "PutObject", err)
	}

	return nil
}

// PutObjectReader uploads data from a reader to S3.
// Content length should be provided for optimal performance.
//
// Example:
//
//	file, _ := os.Open("large-file.zip")
//	defer file.Close()
//	err := client.PutObjectReader(ctx, "my-bucket", "uploads/file.zip", file,
//	    s3.WithContentType("application/zip"),
//	)
func (c *Client) PutObjectReader(ctx context.Context, bucket, key string, reader io.Reader, opts ...PutOption) error {
	if err := validateBucketKey(bucket, key); err != nil {
		return err
	}

	options := &putOptions{}
	for _, opt := range opts {
		opt(options)
	}

	input := &s3.PutObjectInput{
		Bucket: awssdk.String(bucket),
		Key:    awssdk.String(key),
		Body:   reader,
	}

	if options.contentType != "" {
		input.ContentType = awssdk.String(options.contentType)
	}

	if options.storageClass != "" {
		input.StorageClass = options.storageClass
	}

	if options.serverSideEncryption != "" {
		input.ServerSideEncryption = options.serverSideEncryption
	}

	if options.metadata != nil {
		input.Metadata = options.metadata
	}

	_, err := c.uploader.Upload(ctx, input)
	if err != nil {
		return aws.WrapError(serviceName, "PutObject", err)
	}

	return nil
}

// DeleteObject deletes an object from S3.
//
// Example:
//
//	err := client.DeleteObject(ctx, "my-bucket", "path/to/file.txt")
func (c *Client) DeleteObject(ctx context.Context, bucket, key string) error {
	if err := validateBucketKey(bucket, key); err != nil {
		return err
	}

	_, err := c.api.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: awssdk.String(bucket),
		Key:    awssdk.String(key),
	})
	if err != nil {
		return aws.WrapError(serviceName, "DeleteObject", err)
	}

	return nil
}

// DeleteObjects deletes multiple objects from S3.
// Maximum of 1000 objects per call (AWS limit).
//
// Example:
//
//	keys := []string{"file1.txt", "file2.txt", "file3.txt"}
//	deleted, errors := client.DeleteObjects(ctx, "my-bucket", keys)
func (c *Client) DeleteObjects(ctx context.Context, bucket string, keys []string) ([]string, []DeleteError) {
	if bucket == "" {
		return nil, []DeleteError{{Key: "", Error: aws.ErrEmptyBucket.Error()}}
	}

	if len(keys) == 0 {
		return nil, nil
	}

	// AWS limits to 1000 objects per request
	const maxKeys = 1000

	var (
		allDeleted []string
		allErrors  []DeleteError
	)

	for i := 0; i < len(keys); i += maxKeys {
		end := min(i+maxKeys, len(keys))

		batch := keys[i:end]
		objects := make([]types.ObjectIdentifier, len(batch))

		for j, key := range batch {
			objects[j] = types.ObjectIdentifier{Key: awssdk.String(key)}
		}

		output, err := c.api.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: awssdk.String(bucket),
			Delete: &types.Delete{
				Objects: objects,
				Quiet:   awssdk.Bool(false),
			},
		})
		if err != nil {
			// Add all keys as errors
			for _, key := range batch {
				allErrors = append(allErrors, DeleteError{Key: key, Error: err.Error()})
			}

			continue
		}

		for _, deleted := range output.Deleted {
			allDeleted = append(allDeleted, awssdk.ToString(deleted.Key))
		}

		for _, errObj := range output.Errors {
			allErrors = append(allErrors, DeleteError{
				Key:   awssdk.ToString(errObj.Key),
				Error: awssdk.ToString(errObj.Message),
				Code:  awssdk.ToString(errObj.Code),
			})
		}
	}

	return allDeleted, allErrors
}

// CopyObject copies an object within or between buckets.
//
// Example:
//
//	// Copy within same bucket
//	err := client.CopyObject(ctx, "my-bucket", "source.txt", "my-bucket", "dest.txt")
//
//	// Copy to different bucket
//	err := client.CopyObject(ctx, "source-bucket", "file.txt", "dest-bucket", "file.txt")
func (c *Client) CopyObject(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string) error {
	if srcBucket == "" || dstBucket == "" {
		return aws.ErrEmptyBucket
	}

	if srcKey == "" || dstKey == "" {
		return aws.ErrEmptyKey
	}

	copySource := path.Join(srcBucket, srcKey)

	_, err := c.api.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     awssdk.String(dstBucket),
		CopySource: awssdk.String(copySource),
		Key:        awssdk.String(dstKey),
	})
	if err != nil {
		if isNotFoundError(err) {
			return ErrObjectNotFound
		}

		return aws.WrapError(serviceName, "CopyObject", err)
	}

	return nil
}

// MoveObject moves an object by copying then deleting the source.
//
// Example:
//
//	err := client.MoveObject(ctx, "my-bucket", "old-path/file.txt", "my-bucket", "new-path/file.txt")
func (c *Client) MoveObject(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string) error {
	err := c.CopyObject(ctx, srcBucket, srcKey, dstBucket, dstKey)
	if err != nil {
		return err
	}

	return c.DeleteObject(ctx, srcBucket, srcKey)
}

// ObjectExists checks if an object exists in the bucket.
//
// Example:
//
//	exists, err := client.ObjectExists(ctx, "my-bucket", "path/to/file.txt")
func (c *Client) ObjectExists(ctx context.Context, bucket, key string) (bool, error) {
	if err := validateBucketKey(bucket, key); err != nil {
		return false, err
	}

	_, err := c.api.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: awssdk.String(bucket),
		Key:    awssdk.String(key),
	})
	if err != nil {
		if isNotFoundError(err) {
			return false, nil
		}

		return false, aws.WrapError(serviceName, "HeadObject", err)
	}

	return true, nil
}

// HeadObject retrieves object metadata without downloading the content.
//
// Example:
//
//	metadata, err := client.HeadObject(ctx, "my-bucket", "path/to/file.txt")
//	fmt.Printf("Size: %d bytes\n", metadata.ContentLength)
func (c *Client) HeadObject(ctx context.Context, bucket, key string) (*ObjectMetadata, error) {
	if err := validateBucketKey(bucket, key); err != nil {
		return nil, err
	}

	output, err := c.api.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: awssdk.String(bucket),
		Key:    awssdk.String(key),
	})
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrObjectNotFound
		}

		return nil, aws.WrapError(serviceName, "HeadObject", err)
	}

	return &ObjectMetadata{
		ContentType:   awssdk.ToString(output.ContentType),
		ContentLength: awssdk.ToInt64(output.ContentLength),
		ETag:          strings.Trim(awssdk.ToString(output.ETag), "\""),
		LastModified:  awssdk.ToTime(output.LastModified),
		StorageClass:  string(output.StorageClass),
		Metadata:      output.Metadata,
	}, nil
}

// validateBucketKey validates bucket and key are not empty.
func validateBucketKey(bucket, key string) error {
	if bucket == "" {
		return aws.ErrEmptyBucket
	}

	if key == "" {
		return aws.ErrEmptyKey
	}

	return nil
}
