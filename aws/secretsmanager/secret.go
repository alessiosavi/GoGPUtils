package secretsmanager

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

// SecretInfo contains metadata about a secret.
type SecretInfo struct {
	Name             string
	ARN              string
	Description      string
	CreatedDate      time.Time
	LastChangedDate  time.Time
	LastAccessedDate time.Time
	Tags             map[string]string
}

// GetSecretOption configures GetSecret operations.
type GetSecretOption func(*getSecretOptions)

type getSecretOptions struct {
	versionID    string
	versionStage string
}

// WithVersionID retrieves a specific version of the secret.
//
// Example:
//
//	value, err := client.GetSecretString(ctx, "my-secret", secretsmanager.WithVersionID("abc123"))
func WithVersionID(versionID string) GetSecretOption {
	return func(o *getSecretOptions) {
		o.versionID = versionID
	}
}

// WithVersionStage retrieves a specific version stage of the secret.
// Common values: AWSCURRENT (default), AWSPREVIOUS, AWSPENDING.
//
// Example:
//
//	value, err := client.GetSecretString(ctx, "my-secret", secretsmanager.WithVersionStage("AWSPREVIOUS"))
func WithVersionStage(stage string) GetSecretOption {
	return func(o *getSecretOptions) {
		o.versionStage = stage
	}
}

// GetSecretString retrieves a secret value as a string.
//
// Example:
//
//	apiKey, err := client.GetSecretString(ctx, "api-key")
//	if errors.Is(err, secretsmanager.ErrSecretNotFound) {
//	    // Handle not found
//	}
func (c *Client) GetSecretString(ctx context.Context, secretName string, opts ...GetSecretOption) (string, error) {
	if secretName == "" {
		return "", aws.ErrEmptySecret
	}

	options := &getSecretOptions{}
	for _, opt := range opts {
		opt(options)
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId: awssdk.String(secretName),
	}

	if options.versionID != "" {
		input.VersionId = awssdk.String(options.versionID)
	}

	if options.versionStage != "" {
		input.VersionStage = awssdk.String(options.versionStage)
	}

	output, err := c.api.GetSecretValue(ctx, input)
	if err != nil {
		if isResourceNotFound(err) {
			return "", ErrSecretNotFound
		}

		if isInvalidRequest(err) {
			return "", ErrSecretDeleted
		}

		return "", aws.WrapError(serviceName, "GetSecretValue", err)
	}

	return awssdk.ToString(output.SecretString), nil
}

// GetSecretBinary retrieves a secret value as binary data.
//
// Example:
//
//	certData, err := client.GetSecretBinary(ctx, "tls-certificate")
func (c *Client) GetSecretBinary(ctx context.Context, secretName string, opts ...GetSecretOption) ([]byte, error) {
	if secretName == "" {
		return nil, aws.ErrEmptySecret
	}

	options := &getSecretOptions{}
	for _, opt := range opts {
		opt(options)
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId: awssdk.String(secretName),
	}

	if options.versionID != "" {
		input.VersionId = awssdk.String(options.versionID)
	}

	if options.versionStage != "" {
		input.VersionStage = awssdk.String(options.versionStage)
	}

	output, err := c.api.GetSecretValue(ctx, input)
	if err != nil {
		if isResourceNotFound(err) {
			return nil, ErrSecretNotFound
		}

		if isInvalidRequest(err) {
			return nil, ErrSecretDeleted
		}

		return nil, aws.WrapError(serviceName, "GetSecretValue", err)
	}

	return output.SecretBinary, nil
}

// GetSecretJSON retrieves a secret and unmarshals it into the provided destination.
// The secret value must be valid JSON.
//
// Example:
//
//	type DBConfig struct {
//	    Host     string `json:"host"`
//	    Port     int    `json:"port"`
//	    Username string `json:"username"`
//	    Password string `json:"password"`
//	}
//
//	var config DBConfig
//	err := client.GetSecretJSON(ctx, "db-credentials", &config)
func (c *Client) GetSecretJSON(ctx context.Context, secretName string, dest any, opts ...GetSecretOption) error {
	secretString, err := c.GetSecretString(ctx, secretName, opts...)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(secretString), dest); err != nil {
		return aws.WrapError(serviceName, "GetSecretJSON", err)
	}

	return nil
}

// ListSecretsOption configures ListSecrets operations.
type ListSecretsOption func(*listSecretsOptions)

type listSecretsOptions struct {
	maxResults int32
	filters    []types.Filter
}

// WithMaxResults limits the number of secrets returned.
//
// Example:
//
//	secrets, err := client.ListSecrets(ctx, secretsmanager.WithMaxResults(10))
func WithMaxResults(n int32) ListSecretsOption {
	return func(o *listSecretsOptions) {
		o.maxResults = n
	}
}

// WithNameFilter filters secrets by name prefix.
//
// Example:
//
//	secrets, err := client.ListSecrets(ctx, secretsmanager.WithNameFilter("prod/"))
func WithNameFilter(prefix string) ListSecretsOption {
	return func(o *listSecretsOptions) {
		o.filters = append(o.filters, types.Filter{
			Key:    types.FilterNameStringTypeName,
			Values: []string{prefix},
		})
	}
}

// WithTagFilter filters secrets by tag.
//
// Example:
//
//	secrets, err := client.ListSecrets(ctx, secretsmanager.WithTagFilter("env", "production"))
func WithTagFilter(key, value string) ListSecretsOption {
	return func(o *listSecretsOptions) {
		o.filters = append(o.filters, types.Filter{
			Key:    types.FilterNameStringTypeTagKey,
			Values: []string{key},
		})
		o.filters = append(o.filters, types.Filter{
			Key:    types.FilterNameStringTypeTagValue,
			Values: []string{value},
		})
	}
}

// ListSecrets returns all secrets in the account.
//
// Example:
//
//	secrets, err := client.ListSecrets(ctx)
//	for _, secret := range secrets {
//	    fmt.Println(secret.Name)
//	}
func (c *Client) ListSecrets(ctx context.Context, opts ...ListSecretsOption) ([]SecretInfo, error) {
	options := &listSecretsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	input := &secretsmanager.ListSecretsInput{
		SortOrder: types.SortOrderTypeAsc,
	}

	if options.maxResults > 0 {
		input.MaxResults = awssdk.Int32(options.maxResults)
	}

	if len(options.filters) > 0 {
		input.Filters = options.filters
	}

	var secrets []SecretInfo

	for {
		select {
		case <-ctx.Done():
			return secrets, ctx.Err()
		default:
		}

		output, err := c.api.ListSecrets(ctx, input)
		if err != nil {
			return nil, aws.WrapError(serviceName, "ListSecrets", err)
		}

		for _, s := range output.SecretList {
			tags := make(map[string]string)
			for _, tag := range s.Tags {
				tags[awssdk.ToString(tag.Key)] = awssdk.ToString(tag.Value)
			}

			info := SecretInfo{
				Name:        awssdk.ToString(s.Name),
				ARN:         awssdk.ToString(s.ARN),
				Description: awssdk.ToString(s.Description),
				Tags:        tags,
			}

			if s.CreatedDate != nil {
				info.CreatedDate = *s.CreatedDate
			}

			if s.LastChangedDate != nil {
				info.LastChangedDate = *s.LastChangedDate
			}

			if s.LastAccessedDate != nil {
				info.LastAccessedDate = *s.LastAccessedDate
			}

			secrets = append(secrets, info)

			// Check max results limit
			if options.maxResults > 0 && len(secrets) >= int(options.maxResults) {
				return secrets, nil
			}
		}

		if output.NextToken == nil {
			break
		}

		input.NextToken = output.NextToken
	}

	return secrets, nil
}

// ListSecretNames returns just the names of all secrets.
//
// Example:
//
//	names, err := client.ListSecretNames(ctx)
func (c *Client) ListSecretNames(ctx context.Context, opts ...ListSecretsOption) ([]string, error) {
	secrets, err := c.ListSecrets(ctx, opts...)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(secrets))
	for i, s := range secrets {
		names[i] = s.Name
	}

	return names, nil
}

// DescribeSecret returns metadata about a secret without retrieving its value.
//
// Example:
//
//	info, err := client.DescribeSecret(ctx, "my-secret")
//	fmt.Printf("Last changed: %v\n", info.LastChangedDate)
func (c *Client) DescribeSecret(ctx context.Context, secretName string) (*SecretInfo, error) {
	if secretName == "" {
		return nil, aws.ErrEmptySecret
	}

	output, err := c.api.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
		SecretId: awssdk.String(secretName),
	})
	if err != nil {
		if isResourceNotFound(err) {
			return nil, ErrSecretNotFound
		}

		return nil, aws.WrapError(serviceName, "DescribeSecret", err)
	}

	tags := make(map[string]string)
	for _, tag := range output.Tags {
		tags[awssdk.ToString(tag.Key)] = awssdk.ToString(tag.Value)
	}

	info := &SecretInfo{
		Name:        awssdk.ToString(output.Name),
		ARN:         awssdk.ToString(output.ARN),
		Description: awssdk.ToString(output.Description),
		Tags:        tags,
	}

	if output.CreatedDate != nil {
		info.CreatedDate = *output.CreatedDate
	}

	if output.LastChangedDate != nil {
		info.LastChangedDate = *output.LastChangedDate
	}

	if output.LastAccessedDate != nil {
		info.LastAccessedDate = *output.LastAccessedDate
	}

	return info, nil
}
