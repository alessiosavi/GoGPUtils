package s3

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const serviceName = "s3"

// API defines the S3 operations used by this package.
// This interface enables testing with mocks.
type API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)
	CopyObject(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// UploaderAPI defines the S3 upload manager operations.
type UploaderAPI interface {
	Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error)
}

// DownloaderAPI defines the S3 download manager operations.
type DownloaderAPI interface {
	Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error)
}

// Client provides S3 operations.
type Client struct {
	api        API
	uploader   UploaderAPI
	downloader DownloaderAPI
	cfg        *aws.Config
}

// ClientOption configures a Client.
type ClientOption func(*clientOptions)

type clientOptions struct {
	uploaderConcurrency   int
	downloaderConcurrency int
	partSize              int64
}

// WithUploaderConcurrency sets the number of concurrent upload parts.
// Default is 5.
func WithUploaderConcurrency(n int) ClientOption {
	return func(o *clientOptions) {
		o.uploaderConcurrency = n
	}
}

// WithDownloaderConcurrency sets the number of concurrent download parts.
// Default is 5.
func WithDownloaderConcurrency(n int) ClientOption {
	return func(o *clientOptions) {
		o.downloaderConcurrency = n
	}
}

// WithPartSize sets the part size for multipart uploads in bytes.
// Default is 5 MB.
func WithPartSize(size int64) ClientOption {
	return func(o *clientOptions) {
		o.partSize = size
	}
}

// NewClient creates a new S3 client with the given configuration.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx)
//	if err != nil {
//	    return err
//	}
//	client := s3.NewClient(cfg)
func NewClient(cfg *aws.Config, opts ...ClientOption) (*Client, error) {
	if cfg == nil {
		return nil, aws.ErrNilConfig
	}

	options := &clientOptions{
		uploaderConcurrency:   5,
		downloaderConcurrency: 5,
		partSize:              5 * 1024 * 1024, // 5 MB
	}

	for _, opt := range opts {
		opt(options)
	}

	s3Client := s3.NewFromConfig(cfg.AWS())

	uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.Concurrency = options.uploaderConcurrency
		u.PartSize = options.partSize
	})

	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.Concurrency = options.downloaderConcurrency
		d.PartSize = options.partSize
	})

	return &Client{
		api:        s3Client,
		uploader:   uploader,
		downloader: downloader,
		cfg:        cfg,
	}, nil
}

// NewClientWithAPI creates a client with a custom API implementation.
// Useful for testing with mocks.
//
// Example:
//
//	mock := &MockS3API{}
//	client := s3.NewClientWithAPI(mock, nil, nil)
func NewClientWithAPI(api API, uploader UploaderAPI, downloader DownloaderAPI) *Client {
	return &Client{
		api:        api,
		uploader:   uploader,
		downloader: downloader,
	}
}

// API returns the underlying S3 API for direct SDK access.
// Returns nil if the client was created with NewClientWithAPI.
func (c *Client) API() API {
	return c.api
}

// Common S3 errors.
var (
	// ErrObjectNotFound is returned when an object does not exist.
	ErrObjectNotFound = errors.New("s3: object not found")

	// ErrBucketNotFound is returned when a bucket does not exist.
	ErrBucketNotFound = errors.New("s3: bucket not found")

	// ErrAccessDenied is returned when access to a resource is denied.
	ErrAccessDenied = errors.New("s3: access denied")
)

// isNotFoundError checks if the error indicates a not found condition.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var nsk *types.NoSuchKey
	if errors.As(err, &nsk) {
		return true
	}

	var nf *types.NotFound
	if errors.As(err, &nf) {
		return true
	}

	// Check for "NoSuchKey" or "NotFound" in error message as fallback
	errStr := err.Error()

	return strings.Contains(errStr, "NoSuchKey") || strings.Contains(errStr, "NotFound") || strings.Contains(errStr, "404")
}
