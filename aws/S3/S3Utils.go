package S3utils

import (
	"bytes"
	"context"
	"crypto/md5"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"github.com/alessiosavi/GoGPUtils/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html/charset"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

var S3Client *s3.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		S3Client = s3.New(s3.Options{Credentials: cfg.Credentials, Region: cfg.Region, RetryMaxAttempts: 5, RetryMode: aws.RetryModeAdaptive})
	})
}

func GetObject(bucket, fileName string) ([]byte, error) {
	object, err := S3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName)})

	if err != nil {
		return nil, err
	}

	return io.ReadAll(object.Body)
}

func Move(bucket, filename, targetName string) error {
	if err := CopyObject(bucket, bucket, filename, targetName); err != nil {
		return err
	}
	return DeleteObject(bucket, filename)
}

func PutObject(bucket, filename string, data []byte) error {
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

	_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(filename),
		Body:            io.NopCloser(bytes.NewReader(data)),
		ContentEncoding: encoding,
		ContentMD5:      aws.String(string(sum)),
		ContentType:     aws.String(contentType),
	})
	return err
}

func DeleteObject(bucket, key string) error {
	_, err := S3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

// DeleteObjects is delegated to remove the given data (value -> []string) from the given bucket (key -> string)
func DeleteObjects(data map[string][]string) error {
	// TODO: Mange more than 1000 keys and use Thread
	for k := range data {
		var toDelete = make([][]string, 0)
		maxIteration := len(data[k]) / 1000
		for i := 0; i < maxIteration; i++ {
			toDelete = append(toDelete, data[k][1000*i:1000*(i+1)])
		}
		toDelete = append(toDelete, data[k][maxIteration*1000:])

		for i := range toDelete {
			var del = make([]types.ObjectIdentifier, len(toDelete[i]), len(toDelete[i]))
			for j, v := range toDelete[i] {
				del[j] = types.ObjectIdentifier{Key: aws.String(v)}
			}
			objects, err := S3Client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
				Bucket: aws.String(k),
				Delete: &types.Delete{Objects: del},
			})
			if err != nil {
				panic(err)
			}
			if len(objects.Errors) > 0 {
				log.Println("Error for the following data:", helper.MarshalIndent(objects.Errors))
			}
			if len(objects.Deleted) > 0 {
				log.Println("Deleted:", helper.MarshalIndent(objects.Deleted))
			}
		}
	}
	return nil
}

func PutObjectStream(bucket, filename string, stream io.ReadCloser, contentType, encoding, md5 *string) error {
	defer stream.Close()
	uploader := manager.NewUploader(S3Client)
	uploader.Concurrency = 10

	_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(filename),
		Body:            stream,
		ContentEncoding: encoding,
		ContentMD5:      md5,
		ContentType:     contentType,
	})
	return err
}

// ListBucketObjectsDetails is delegated to list all the objects (details) in the given bucket. Prefix is optional. The result is return ordered by the last modified
func ListBucketObjectsDetails(bucket, prefix string) ([]types.Object, error) {
	objects, err := S3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	var buckets []types.Object
	buckets = append(buckets, objects.Contents...)

	continuationToken := objects.NextContinuationToken
	truncated := objects.IsTruncated
	for *truncated {
		newObjects, err := S3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(prefix),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		continuationToken = newObjects.NextContinuationToken
		buckets = append(buckets, newObjects.Contents...)

		truncated = newObjects.IsTruncated
	}

	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].LastModified.Before(*buckets[i].LastModified)
	})
	return buckets, nil
}

// ListBucketObjects is delegated to list all the objects (name only) in the given bucket. Prefix is optional. The result is return ordered by the last modified
func ListBucketObjects(bucket, prefix string) ([]string, error) {
	objects, err := ListBucketObjectsDetails(bucket, prefix)
	if err != nil {
		return nil, err
	}
	var data = make([]string, len(objects))
	for i := range objects {
		data[i] = *objects[i].Key
	}
	return data, nil

}

func CopyObject(bucketSource, bucketTarget, keySource, keyTarget string) error {
	if _, err := S3Client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     aws.String(bucketTarget),
		CopySource: aws.String(path.Join(bucketSource, keySource)),
		Key:        aws.String(keyTarget),
	}); err != nil {
		return err
	}
	return nil
}

func ObjectExists(bucket, key string) bool {
	if _, err := S3Client.HeadObject(context.Background(), &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}); err != nil {
		return false
	}
	return true
}

func SyncBucket(bucket, prefix string, bucketsTarget ...string) ([]string, error) {
	var fileNotSynced []string

	objects, err := ListBucketObjects(bucket, prefix)
	if err != nil {
		return nil, err
	}
	// Cache the name of the files for the given bucket, by this way we avoid unnecessary copy/head operation
	// NOTE: The file will not be copied if it has the same name but with different content
	var cache = make(map[string]map[string]struct{})

	// Iterate the buckets and save the file name in a map
	for _, b := range bucketsTarget {
		bucketObjects, err := ListBucketObjects(b, "")
		if err != nil {
			return nil, err
		}
		cache[b] = make(map[string]struct{})
		for i := range bucketObjects {
			cache[b][bucketObjects[i]] = struct{}{}
		}
	}

	for _, object := range objects {
		for _, bucketTarget := range bucketsTarget {
			//if IsDifferent(bucket, bucketTarget, object, object) {
			if _, ok := cache[bucketTarget][object]; !ok { // File not exists
				log.Printf("Copying %s\n", path.Join(bucket, object))
				if err = CopyObject(bucket, bucketTarget, object, object); err != nil {
					fileNotSynced = append(fileNotSynced, object)
				}
			}
			//} else {
			//	log.Printf("File %s skipped cause it already exists\n", path.Join(bucket, object))
			//}
		}
	}

	return fileNotSynced, nil
}

func HeadObject(bucket, key string) (*s3.HeadObjectOutput, error) {
	return S3Client.HeadObject(context.Background(), &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
}

func IsDifferent(bucket_base, bucket_target, key_base, key_target string) bool {
	var (
		head_base, head_target *s3.HeadObjectOutput
		err1, err2             error
	)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		head_base, err1 = HeadObject(bucket_base, key_base)
		wg.Done()
	}()
	go func() {
		head_target, err2 = HeadObject(bucket_target, key_target)
		wg.Done()
	}()
	wg.Wait()

	if err1 != nil || err2 != nil {
		if err1 != nil {
			log.Println("INPUT: ", bucket_base, key_base)
			log.Println(err1)
		} else {
			log.Println("INPUT: ", bucket_target, key_target)
			log.Println(err2)
		}
		return true
	}
	return head_base.ContentLength != head_target.ContentLength || *head_base.ETag != *head_target.ETag
}
func IsDifferentLegacy(bucket_base, bucket_target, key_base, key_target string) bool {
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

// FIXME: Use a list of string as prefix
func GetAfterDate(bucket, prefix string, date time.Time) ([]types.Object, error) {
	details, err := ListBucketObjectsDetails(bucket, prefix)
	if err != nil {
		return nil, err
	}
	var res []types.Object

	sort.Slice(details, func(i, j int) bool {
		return details[i].LastModified.Before(*details[j].LastModified)
	})
	for _, detail := range details {
		if detail.LastModified.After(date) {
			res = append(res, detail)
		}
	}
	return res, nil
}

func GetBetweenDate(bucket, prefix string, start, stop time.Time) ([]string, error) {
	details, err := ListBucketObjectsDetails(bucket, prefix)
	if err != nil {
		return nil, err
	}
	var res []string

	for _, detail := range details {
		if detail.LastModified.After(start) && detail.LastModified.Before(stop) {
			res = append(res, *detail.Key)
		}
	}
	return res, nil
}

func SyncAfterDate(bucket, prefix, localPath string, date time.Time) error {
	details, err := ListBucketObjectsDetails(bucket, strings.TrimLeft(prefix, "/"))
	if err != nil {
		return err
	}

	var detailsAfter []types.Object
	for _, detail := range details {
		if detail.LastModified.After(date) {
			detailsAfter = append(detailsAfter, detail)
		}
	}
	details = nil
	bar := progressbar.Default(int64(len(detailsAfter)))
	defer bar.Close()
	for _, detail := range detailsAfter {
		bar.Describe(*detail.Key)
		object, err := GetObject(bucket, *detail.Key)
		if err != nil {
			return err
		}
		basepath := path.Join(localPath, path.Dir(*detail.Key))
		if !fileutils.IsDir(basepath) {
			if err = os.MkdirAll(basepath, 0755); err != nil {
				return err
			}
		}
		f := path.Join(basepath, path.Base(*detail.Key))
		if err = os.WriteFile(f, object, 0755); err != nil {
			return err
		}
		if err = os.Chtimes(f, *detail.LastModified, *detail.LastModified); err != nil {
			log.Println("Unable to set time for file,", f, "ERR:", err.Error())
		}
		bar.Add(1)
	}
	return nil
}

// ParseS3Path is delegated to return bucket name and filename of a given s3 path
func ParseS3Path(p string) (string, string) {
	p = strings.TrimLeft(p, "s3://")
	split := strings.Split(p, "/")
	return split[0], path.Clean(strings.Join(split[1:], "/"))
}
