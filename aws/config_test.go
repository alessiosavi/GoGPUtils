package aws

import (
	"context"
	"os"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func TestLoadConfig_AppliesEndpointAndRetryOptions(t *testing.T) {
	t.Setenv("AWS_CONFIG_FILE", os.DevNull)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", os.DevNull)
	t.Setenv("AWS_PROFILE", "")
	for _, tt := range []struct {
		name     string
		mode     awssdk.RetryMode
		attempts int
	}{
		{"standard", awssdk.RetryModeStandard, 7},
		{"adaptive", awssdk.RetryModeAdaptive, 5},
		{"default", "", DefaultRetryMaxAttempts},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts := []ConfigOption{WithRegion("us-east-1"), WithEndpoint("http://localhost:4566"),
				WithCredentials(credentials.NewStaticCredentialsProvider("test", "test", ""))}
			if tt.mode != "" {
				opts = append(opts, WithRetryMode(tt.mode), WithRetryMaxAttempts(tt.attempts))
			}
			cfg, err := LoadConfig(context.Background(), opts...)
			if err != nil {
				t.Fatal(err)
			}
			if got := awssdk.ToString(cfg.BaseEndpoint); got != "http://localhost:4566" {
				t.Errorf("BaseEndpoint = %q", got)
			}
			r := cfg.Retryer()
			if tt.mode == awssdk.RetryModeStandard {
				if _, ok := r.(*retry.Standard); !ok {
					t.Errorf("retryer = %T, want *retry.Standard", r)
				}
			} else if _, ok := r.(*retry.AdaptiveMode); !ok {
				t.Errorf("retryer = %T, want *retry.AdaptiveMode", r)
			}
			if r.MaxAttempts() != tt.attempts {
				t.Errorf("attempts = %d, want %d", r.MaxAttempts(), tt.attempts)
			}
		})
	}
}
