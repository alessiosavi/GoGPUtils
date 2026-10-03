// Package s3 provides helpers for Amazon S3 operations.
//
// # Features
//
//   - Object operations: Get, Put, Delete, Copy, Move
//   - Bucket operations: List objects with filtering
//   - Streaming uploads and downloads
//   - Automatic content type detection
//   - Concurrent multipart uploads and downloads
//
// # Client Creation
//
// Create a client using AWS configuration:
//
//	cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
//	if err != nil {
//	    return err
//	}
//
//	client, err := s3.NewClient(cfg)
//	if err != nil {
//	    return err
//	}
//
// # Basic Operations
//
//	// Upload an object
//	err = client.PutObject(ctx, "my-bucket", "path/to/file.txt", data)
//
//	// Download an object
//	data, err := client.GetObject(ctx, "my-bucket", "path/to/file.txt")
//
//	// Delete an object
//	err = client.DeleteObject(ctx, "my-bucket", "path/to/file.txt")
//
//	// List objects
//	objects, err := client.ListObjects(ctx, "my-bucket", s3.WithPrefix("path/"))
//
// # Testing
//
// For testing, use the interface-based client:
//
// Implement s3.API with a mock and pass it to the constructor. For example,
// this mock supports GetObject; it embeds s3.API for the unused operations.
// Import bytes, context, and io, and alias the SDK service package as s3sdk.
//
//	type mockS3API struct {
//	    s3.API
//	    data []byte
//	}
//
//	func (m *mockS3API) GetObject(ctx context.Context, input *s3sdk.GetObjectInput, opts ...func(*s3sdk.Options)) (*s3sdk.GetObjectOutput, error) {
//	    return &s3sdk.GetObjectOutput{Body: io.NopCloser(bytes.NewReader(m.data))}, nil
//	}
//
// In a test function:
//
//	mock := &mockS3API{data: []byte("content")}
//	client := s3.NewClientWithAPI(mock, nil, nil)
package s3
