// Package s3 provides helpers for Amazon S3 operations.
//
// # Features
//
//   - Object operations: Get, Put, Delete, Copy, Move
//   - Bucket operations: List objects with filtering
//   - Streaming uploads and downloads
//   - Automatic content type detection
//   - Parallel operations for batch processing
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
//	client := s3.NewClient(cfg)
//
// # Basic Operations
//
//	// Upload an object
//	err := client.PutObject(ctx, "my-bucket", "path/to/file.txt", data)
//
//	// Download an object
//	data, err := client.GetObject(ctx, "my-bucket", "path/to/file.txt")
//
//	// Delete an object
//	err := client.DeleteObject(ctx, "my-bucket", "path/to/file.txt")
//
//	// List objects
//	objects, err := client.ListObjects(ctx, "my-bucket", s3.WithPrefix("path/"))
//
// # Testing
//
// For testing, use the interface-based client:
//
//	type mockS3API struct {
//	    s3.API
//	    getObjectFunc func(...) (...)
//	}
//
//	mock := &mockS3API{
//	    getObjectFunc: func(...) (...) { ... },
//	}
//	client := s3.NewClientWithAPI(mock)
package s3
