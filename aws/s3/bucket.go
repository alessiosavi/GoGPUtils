package s3

import (
	"context"
	"sort"
	"strings"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ListObjects returns all objects in a bucket matching the options.
// Results are sorted by LastModified ascending.
//
// Example:
//
//	// List all objects
//	objects, err := client.ListObjects(ctx, "my-bucket")
//
//	// List with prefix
//	objects, err := client.ListObjects(ctx, "my-bucket", s3.WithPrefix("logs/"))
//
//	// List with time filter
//	since := time.Now().Add(-24 * time.Hour)
//	objects, err := client.ListObjects(ctx, "my-bucket", s3.WithModifiedAfter(since))
func (c *Client) ListObjects(ctx context.Context, bucket string, opts ...ListOption) ([]ObjectInfo, error) {
	if bucket == "" {
		return nil, aws.ErrEmptyBucket
	}

	options := &listOptions{}
	for _, opt := range opts {
		opt(options)
	}

	input := &s3.ListObjectsV2Input{
		Bucket: awssdk.String(bucket),
	}

	if options.prefix != "" {
		input.Prefix = awssdk.String(options.prefix)
	}

	if options.delimiter != "" {
		input.Delimiter = awssdk.String(options.delimiter)
	}

	if options.maxKeys > 0 {
		input.MaxKeys = awssdk.Int32(options.maxKeys)
	}

	if options.startAfter != "" {
		input.StartAfter = awssdk.String(options.startAfter)
	}

	var objects []ObjectInfo

	for {
		select {
		case <-ctx.Done():
			return objects, ctx.Err()
		default:
		}

		output, err := c.api.ListObjectsV2(ctx, input)
		if err != nil {
			return nil, aws.WrapError(serviceName, "ListObjectsV2", err)
		}

		for _, obj := range output.Contents {
			// Apply time filters if specified
			lastModified := awssdk.ToTime(obj.LastModified)

			if options.modifiedFrom != nil && lastModified.Before(*options.modifiedFrom) {
				continue
			}

			if options.modifiedTo != nil && lastModified.After(*options.modifiedTo) {
				continue
			}

			objects = append(objects, ObjectInfo{
				Key:          awssdk.ToString(obj.Key),
				Size:         awssdk.ToInt64(obj.Size),
				ETag:         strings.Trim(awssdk.ToString(obj.ETag), "\""),
				LastModified: lastModified,
				StorageClass: string(obj.StorageClass),
			})

			// Check if we've hit the max keys limit
			if options.maxKeys > 0 && len(objects) >= int(options.maxKeys) {
				break
			}
		}

		// Check if we need to continue pagination
		if !awssdk.ToBool(output.IsTruncated) {
			break
		}

		// Check max keys limit
		if options.maxKeys > 0 && len(objects) >= int(options.maxKeys) {
			break
		}

		input.ContinuationToken = output.NextContinuationToken
	}

	// Sort by LastModified ascending
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].LastModified.Before(objects[j].LastModified)
	})

	return objects, nil
}

// ListObjectKeys returns just the keys of objects in a bucket.
// More efficient than ListObjects when you only need keys.
//
// Example:
//
//	keys, err := client.ListObjectKeys(ctx, "my-bucket", s3.WithPrefix("data/"))
func (c *Client) ListObjectKeys(ctx context.Context, bucket string, opts ...ListOption) ([]string, error) {
	objects, err := c.ListObjects(ctx, bucket, opts...)
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(objects))
	for i, obj := range objects {
		keys[i] = obj.Key
	}

	return keys, nil
}

// ListObjectsCallback iterates through objects and calls a callback for each.
// Useful for processing large numbers of objects without loading all into memory.
// Return an error from the callback to stop iteration.
//
// Example:
//
//	err := client.ListObjectsCallback(ctx, "my-bucket", func(obj s3.ObjectInfo) error {
//	    fmt.Printf("Found: %s (%d bytes)\n", obj.Key, obj.Size)
//	    return nil
//	}, s3.WithPrefix("logs/"))
func (c *Client) ListObjectsCallback(ctx context.Context, bucket string, callback func(ObjectInfo) error, opts ...ListOption) error {
	if bucket == "" {
		return aws.ErrEmptyBucket
	}

	options := &listOptions{}
	for _, opt := range opts {
		opt(options)
	}

	input := &s3.ListObjectsV2Input{
		Bucket: awssdk.String(bucket),
	}

	if options.prefix != "" {
		input.Prefix = awssdk.String(options.prefix)
	}

	if options.delimiter != "" {
		input.Delimiter = awssdk.String(options.delimiter)
	}

	if options.startAfter != "" {
		input.StartAfter = awssdk.String(options.startAfter)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		output, err := c.api.ListObjectsV2(ctx, input)
		if err != nil {
			return aws.WrapError(serviceName, "ListObjectsV2", err)
		}

		for _, obj := range output.Contents {
			lastModified := awssdk.ToTime(obj.LastModified)

			// Apply time filters
			if options.modifiedFrom != nil && lastModified.Before(*options.modifiedFrom) {
				continue
			}

			if options.modifiedTo != nil && lastModified.After(*options.modifiedTo) {
				continue
			}

			info := ObjectInfo{
				Key:          awssdk.ToString(obj.Key),
				Size:         awssdk.ToInt64(obj.Size),
				ETag:         strings.Trim(awssdk.ToString(obj.ETag), "\""),
				LastModified: lastModified,
				StorageClass: string(obj.StorageClass),
			}

			err := callback(info)
			if err != nil {
				return err
			}
		}

		if !awssdk.ToBool(output.IsTruncated) {
			break
		}

		input.ContinuationToken = output.NextContinuationToken
	}

	return nil
}

// CountObjects returns the count of objects matching the options.
//
// Example:
//
//	count, err := client.CountObjects(ctx, "my-bucket", s3.WithPrefix("logs/"))
func (c *Client) CountObjects(ctx context.Context, bucket string, opts ...ListOption) (int64, error) {
	var count int64

	err := c.ListObjectsCallback(ctx, bucket, func(obj ObjectInfo) error {
		count++

		return nil
	}, opts...)

	return count, err
}

// TotalSize returns the total size in bytes of objects matching the options.
//
// Example:
//
//	size, err := client.TotalSize(ctx, "my-bucket", s3.WithPrefix("logs/"))
//	fmt.Printf("Total size: %.2f MB\n", float64(size)/1024/1024)
func (c *Client) TotalSize(ctx context.Context, bucket string, opts ...ListOption) (int64, error) {
	var total int64

	err := c.ListObjectsCallback(ctx, bucket, func(obj ObjectInfo) error {
		total += obj.Size

		return nil
	}, opts...)

	return total, err
}

// ParseS3Path parses an S3 URI (s3://bucket/key) into bucket and key components.
//
// Example:
//
//	bucket, key, err := s3.ParseS3Path("s3://my-bucket/path/to/file.txt")
//	// bucket = "my-bucket", key = "path/to/file.txt"
func ParseS3Path(s3Path string) (bucket, key string, err error) {
	// Remove s3:// prefix
	path := strings.TrimPrefix(s3Path, "s3://")
	path = strings.TrimPrefix(path, "S3://")

	// Find first slash
	idx := strings.Index(path, "/")
	if idx == -1 {
		// No key, just bucket
		return path, "", nil
	}

	bucket = path[:idx]
	key = strings.TrimPrefix(path[idx:], "/")

	if bucket == "" {
		return "", "", aws.NewValidationError("s3Path", "invalid S3 path: missing bucket")
	}

	return bucket, key, nil
}

// FormatS3Path formats a bucket and key into an S3 URI.
//
// Example:
//
//	path := s3.FormatS3Path("my-bucket", "path/to/file.txt")
//	// path = "s3://my-bucket/path/to/file.txt"
func FormatS3Path(bucket, key string) string {
	if key == "" {
		return "s3://" + bucket
	}

	return "s3://" + bucket + "/" + key
}
