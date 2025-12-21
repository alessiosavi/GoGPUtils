package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
)

// DefaultRetryMaxAttempts is the default number of retry attempts for AWS operations.
const DefaultRetryMaxAttempts = 3

// DefaultTimeout is the default timeout for AWS operations.
const DefaultTimeout = 30 * time.Second

// Config wraps aws.Config with additional helper methods.
type Config struct {
	aws.Config
}

// ConfigOption configures AWS config loading.
type ConfigOption func(*configOptions)

type configOptions struct {
	region           string
	profile          string
	retryMaxAttempts int
	retryMode        aws.RetryMode
	endpoint         string
	credentials      aws.CredentialsProvider
	loadOptions      []func(*config.LoadOptions) error
}

// WithRegion sets the AWS region.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
func WithRegion(region string) ConfigOption {
	return func(o *configOptions) {
		o.region = region
	}
}

// WithProfile sets the AWS profile to use from shared credentials.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx, aws.WithProfile("production"))
func WithProfile(profile string) ConfigOption {
	return func(o *configOptions) {
		o.profile = profile
	}
}

// WithRetryMaxAttempts sets the maximum number of retry attempts.
// Default is 3 attempts.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx, aws.WithRetryMaxAttempts(5))
func WithRetryMaxAttempts(attempts int) ConfigOption {
	return func(o *configOptions) {
		o.retryMaxAttempts = attempts
	}
}

// WithRetryMode sets the retry mode (standard or adaptive).
// Default is adaptive retry mode.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx, aws.WithRetryMode(aws.RetryModeAdaptive))
func WithRetryMode(mode aws.RetryMode) ConfigOption {
	return func(o *configOptions) {
		o.retryMode = mode
	}
}

// WithEndpoint sets a custom endpoint URL for all services.
// Useful for testing with LocalStack or other AWS-compatible services.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx, aws.WithEndpoint("http://localhost:4566"))
func WithEndpoint(endpoint string) ConfigOption {
	return func(o *configOptions) {
		o.endpoint = endpoint
	}
}

// WithCredentials sets explicit credentials provider.
//
// Example:
//
//	creds := credentials.NewStaticCredentialsProvider("key", "secret", "")
//	cfg, err := aws.LoadConfig(ctx, aws.WithCredentials(creds))
func WithCredentials(provider aws.CredentialsProvider) ConfigOption {
	return func(o *configOptions) {
		o.credentials = provider
	}
}

// WithLoadOption adds a raw AWS SDK config load option.
// Use this for advanced configuration not covered by other options.
func WithLoadOption(opt func(*config.LoadOptions) error) ConfigOption {
	return func(o *configOptions) {
		o.loadOptions = append(o.loadOptions, opt)
	}
}

// LoadConfig loads AWS configuration with the specified options.
// If no options are provided, it uses default configuration from
// environment variables, shared credentials, and IAM roles.
//
// Example:
//
//	// Load with defaults
//	cfg, err := aws.LoadConfig(ctx)
//
//	// Load with specific region and profile
//	cfg, err := aws.LoadConfig(ctx,
//	    aws.WithRegion("us-west-2"),
//	    aws.WithProfile("production"),
//	)
func LoadConfig(ctx context.Context, opts ...ConfigOption) (*Config, error) {
	options := &configOptions{
		retryMaxAttempts: DefaultRetryMaxAttempts,
		retryMode:        aws.RetryModeAdaptive,
	}

	for _, opt := range opts {
		opt(options)
	}

	var loadOpts []func(*config.LoadOptions) error

	if options.region != "" {
		loadOpts = append(loadOpts, config.WithRegion(options.region))
	}

	if options.profile != "" {
		loadOpts = append(loadOpts, config.WithSharedConfigProfile(options.profile))
	}

	if options.credentials != nil {
		loadOpts = append(loadOpts, config.WithCredentialsProvider(options.credentials))
	}

	// Configure retry behavior
	loadOpts = append(loadOpts, config.WithRetryer(func() aws.Retryer {
		return retry.NewAdaptiveMode(
			func(o *retry.AdaptiveModeOptions) {
				o.StandardOptions = append(o.StandardOptions,
					func(so *retry.StandardOptions) {
						so.MaxAttempts = options.retryMaxAttempts
					},
				)
			},
		)
	}))

	// Add any custom load options
	loadOpts = append(loadOpts, options.loadOptions...)

	cfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, &ConfigError{Err: err}
	}

	return &Config{Config: cfg}, nil
}

// AWS returns the underlying aws.Config for direct SDK usage.
func (c *Config) AWS() aws.Config {
	return c.Config
}

// Region returns the configured region.
func (c *Config) Region() string {
	return c.Config.Region
}
