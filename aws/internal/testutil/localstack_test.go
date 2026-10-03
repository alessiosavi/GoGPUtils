package testutil

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type httpClientFunc func(*http.Request) (*http.Response, error)

func (f httpClientFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestLoadConfig_SetsLocalStackBaseEndpoint(t *testing.T) {
	t.Setenv("LOCALSTACK_ENDPOINT", "http://localhost:4566")
	cfg, err := LoadConfig(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := awssdk.ToString(cfg.BaseEndpoint); got != LocalStackEndpoint() {
		t.Fatalf("BaseEndpoint = %q", got)
	}

	cfg.HTTPClient = httpClientFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.URL.String(); got != "http://localhost:4566/test-bucket/key?x-id=GetObject" {
			t.Errorf("request URL = %q; want the LocalStack host with a bucket path", got)
		}

		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("data"))}, nil
	})
	output, err := s3.NewFromConfig(cfg.AWS()).GetObject(context.Background(), &s3.GetObjectInput{Bucket: awssdk.String("test-bucket"), Key: awssdk.String("key")})
	if err != nil {
		t.Fatal(err)
	}
	output.Body.Close()
}
