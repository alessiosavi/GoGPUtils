package S3

import (
	"bytes"
	"context"
	"crypto/md5"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

func GetObject(bucket, fileName string) ([]byte, error) {
	cfg, err := awsutils.New()
	if err != nil {
		return nil, err
	}

	S3Client := s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	s3CsvConf, err := S3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName)})

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(s3CsvConf.Body)
	return data, err
}

func PutObject(bucket, filename string, data []byte) error {
	cfg, err := awsutils.New()
	if err != nil {
		return err
	}

	S3Client := s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	h := md5.New()
	var sum []byte = nil
	if _, err := h.Write(data); err != nil {
		sum = h.Sum(nil)
	}

	contentType := http.DetectContentType(data)
	_, enc, ok := charset.DetermineEncoding(data, contentType)
	var encoding *string = nil
	if ok {
		encoding = aws.String(enc)
	}

	uploader := manager.NewUploader(S3Client)

	_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(filename),
		Body:            ioutil.NopCloser(bytes.NewReader(data)),
		ContentEncoding: encoding,
		ContentMD5:      aws.String(string(sum)),
		ContentType:     aws.String(contentType),
	})
	return err
}

func DeleteObject(bucket, key string) error {
	cfg, err := awsutils.New()
	if err != nil {
		return err
	}
	S3Client := s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	_, err = S3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

func PutObjectStream(bucket, filename string, stream io.ReadCloser, contentType, encoding, md5 *string) error {
	cfg, err := awsutils.New()
	if err != nil {
		return err
	}

	S3Client := s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	uploader := manager.NewUploader(S3Client)

	_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(filename),
		Body:            stream,
		ContentEncoding: encoding,
		ContentMD5:      md5,
		ContentType:     contentType,
	})
	return err
}

func ListBucketObject(bucket string) ([]string, error) {
	cfg, err := awsutils.New()
	if err != nil {
		return nil, err
	}
	S3Client := s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region})

	objects, err := S3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}
	continuationToken := objects.NextContinuationToken

	var buckets = make([]string, len(objects.Contents))

	for i, bucketName := range objects.Contents {
		buckets[i] = *bucketName.Key
	}

	for continuationToken != nil {
		newObjects, err := S3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			ContinuationToken: continuationToken,
			Bucket:            aws.String(bucket),
		})
		if err != nil {
			return nil, err
		}
		continuationToken = newObjects.NextContinuationToken
		for _, bucketName := range objects.Contents {
			buckets = append(buckets, *bucketName.Key)
		}
	}

	return buckets, nil
}

func CopyObject(bucketSource, bucketTarget, keySource, keyTarget string) error {
	cfg, err := awsutils.New()
	if err != nil {
		return err
	}
	S3Client := s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region})

	if _, err = S3Client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     aws.String(bucketTarget),
		CopySource: aws.String(path.Join(bucketSource, keySource)),
		Key:        aws.String(keyTarget),
	}); err != nil {
		return err
	}
	return nil
}

func ObjectExists(bucket, key string) bool {
	cfg, err := awsutils.New()
	if err != nil {
		return false
	}
	S3Client := s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	if _, err = S3Client.HeadObject(context.Background(), &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}); err != nil {
		return false
	}
	return true
}
func SyncBucket(bucket string, bucketsTarget ...string) ([]string, error) {
	var fileNotSynced []string
	var err error
	objects, err := ListBucketObject(bucket)
	if err != nil {
		return nil, err
	}
	for _, object := range objects {
		for _, bucketTarget := range bucketsTarget {
			exist := ObjectExists(bucketTarget, object)
			if !exist || IsDifferent(bucket, bucketTarget, object, object) {
				log.Printf("Copying %s\n", path.Join(bucket, object))
				if err = CopyObject(bucket, bucketTarget, object, object); err != nil {
					fileNotSynced = append(fileNotSynced, object)
				}
			} else {
				log.Printf("File %s skipped cause it already exists\n", path.Join(bucket, object))
			}
		}
	}

	return fileNotSynced, nil
}

func IsDifferent(bucket_base, bucket_target, key_base, key_target string) bool {
	cfg, err := awsutils.New()
	if err != nil {
		return true
	}
	S3Client := s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	head_base, err := S3Client.HeadObject(context.Background(), &s3.HeadObjectInput{Bucket: aws.String(bucket_base), Key: aws.String(key_base)})
	if err != nil {
		return true
	}
	head_target, err := S3Client.HeadObject(context.Background(), &s3.HeadObjectInput{Bucket: aws.String(bucket_target), Key: aws.String(key_target)})
	if err != nil {
		return true
	}

	return *head_base.ETag != *head_target.ETag || head_base.ContentLength != head_target.ContentLength
}
