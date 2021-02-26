package S3

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
)
import "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"
import "golang.org/x/net/html/charset"

func GetObject(bucket, fileName string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
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
	cfg, err := config.LoadDefaultConfig(context.Background())
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

	object, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(filename),
		Body:            ioutil.NopCloser(bytes.NewReader(data)),
		ContentEncoding: encoding,
		ContentMD5:      aws.String(string(sum)),
		ContentType:     aws.String(contentType),
	})
	fmt.Printf("%+v\n", object)
	return err
}
