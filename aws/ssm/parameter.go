package ssm

import (
	"context"
	"time"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// ParameterInfo contains metadata about a parameter.
type ParameterInfo struct {
	Name             string
	Type             string
	Value            string
	Version          int64
	LastModifiedDate time.Time
	ARN              string
	DataType         string
}

// GetParameterOption configures GetParameter operations.
type GetParameterOption func(*getParameterOptions)

type getParameterOptions struct {
	withDecryption bool
}

// WithDecryption decrypts SecureString parameters.
// This is enabled by default.
//
// Example:
//
//	value, err := client.GetParameter(ctx, "/app/secret", ssm.WithDecryption(true))
func WithDecryption(decrypt bool) GetParameterOption {
	return func(o *getParameterOptions) {
		o.withDecryption = decrypt
	}
}

// GetParameter retrieves a parameter value.
// SecureString parameters are automatically decrypted.
//
// Example:
//
//	value, err := client.GetParameter(ctx, "/app/config/database_url")
//	if errors.Is(err, ssm.ErrParameterNotFound) {
//	    // Handle not found
//	}
func (c *Client) GetParameter(ctx context.Context, name string, opts ...GetParameterOption) (string, error) {
	if name == "" {
		return "", aws.ErrEmptyParameter
	}

	options := &getParameterOptions{
		withDecryption: true, // Default to decrypting
	}
	for _, opt := range opts {
		opt(options)
	}

	output, err := c.api.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           awssdk.String(name),
		WithDecryption: awssdk.Bool(options.withDecryption),
	})
	if err != nil {
		if isParameterNotFound(err) {
			return "", ErrParameterNotFound
		}

		return "", aws.WrapError(serviceName, "GetParameter", err)
	}

	return awssdk.ToString(output.Parameter.Value), nil
}

// GetParameterInfo retrieves a parameter with its metadata.
//
// Example:
//
//	info, err := client.GetParameterInfo(ctx, "/app/config/database_url")
//	fmt.Printf("Version: %d, Last Modified: %v\n", info.Version, info.LastModifiedDate)
func (c *Client) GetParameterInfo(ctx context.Context, name string, opts ...GetParameterOption) (*ParameterInfo, error) {
	if name == "" {
		return nil, aws.ErrEmptyParameter
	}

	options := &getParameterOptions{
		withDecryption: true,
	}
	for _, opt := range opts {
		opt(options)
	}

	output, err := c.api.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           awssdk.String(name),
		WithDecryption: awssdk.Bool(options.withDecryption),
	})
	if err != nil {
		if isParameterNotFound(err) {
			return nil, ErrParameterNotFound
		}

		return nil, aws.WrapError(serviceName, "GetParameter", err)
	}

	info := &ParameterInfo{
		Name:     awssdk.ToString(output.Parameter.Name),
		Type:     string(output.Parameter.Type),
		Value:    awssdk.ToString(output.Parameter.Value),
		Version:  output.Parameter.Version,
		ARN:      awssdk.ToString(output.Parameter.ARN),
		DataType: awssdk.ToString(output.Parameter.DataType),
	}

	if output.Parameter.LastModifiedDate != nil {
		info.LastModifiedDate = *output.Parameter.LastModifiedDate
	}

	return info, nil
}

// GetParameters retrieves multiple parameters in one call.
// Returns a map of parameter name to value.
// Invalid parameter names are returned separately.
//
// Example:
//
//	names := []string{"/app/config/a", "/app/config/b", "/app/config/c"}
//	values, invalid, err := client.GetParameters(ctx, names)
func (c *Client) GetParameters(ctx context.Context, names []string, opts ...GetParameterOption) (map[string]string, []string, error) {
	if len(names) == 0 {
		return nil, nil, nil
	}

	options := &getParameterOptions{
		withDecryption: true,
	}
	for _, opt := range opts {
		opt(options)
	}

	// AWS limits to 10 parameters per request
	const maxParams = 10

	values := make(map[string]string)

	var invalidNames []string

	for i := 0; i < len(names); i += maxParams {
		end := min(i+maxParams, len(names))

		batch := names[i:end]

		output, err := c.api.GetParameters(ctx, &ssm.GetParametersInput{
			Names:          batch,
			WithDecryption: awssdk.Bool(options.withDecryption),
		})
		if err != nil {
			return values, invalidNames, aws.WrapError(serviceName, "GetParameters", err)
		}

		for _, param := range output.Parameters {
			values[awssdk.ToString(param.Name)] = awssdk.ToString(param.Value)
		}

		invalidNames = append(invalidNames, output.InvalidParameters...)
	}

	return values, invalidNames, nil
}

// ListParametersByPathOption configures ListParametersByPath operations.
type ListParametersByPathOption func(*listByPathOptions)

type listByPathOptions struct {
	recursive      bool
	withDecryption bool
	maxResults     int32
}

// WithRecursive enables recursive listing of parameters.
//
// Example:
//
//	params, err := client.ListParametersByPath(ctx, "/app/", ssm.WithRecursive(true))
func WithRecursive(recursive bool) ListParametersByPathOption {
	return func(o *listByPathOptions) {
		o.recursive = recursive
	}
}

// WithPathDecryption decrypts SecureString parameters.
//
// Example:
//
//	params, err := client.ListParametersByPath(ctx, "/app/", ssm.WithPathDecryption(true))
func WithPathDecryption(decrypt bool) ListParametersByPathOption {
	return func(o *listByPathOptions) {
		o.withDecryption = decrypt
	}
}

// WithPathMaxResults limits the number of parameters returned.
//
// Example:
//
//	params, err := client.ListParametersByPath(ctx, "/app/", ssm.WithPathMaxResults(10))
func WithPathMaxResults(n int32) ListParametersByPathOption {
	return func(o *listByPathOptions) {
		o.maxResults = n
	}
}

// ListParametersByPath returns all parameters under a path.
//
// Example:
//
//	params, err := client.ListParametersByPath(ctx, "/app/config/")
//	for _, p := range params {
//	    fmt.Printf("%s = %s\n", p.Name, p.Value)
//	}
func (c *Client) ListParametersByPath(ctx context.Context, path string, opts ...ListParametersByPathOption) ([]ParameterInfo, error) {
	if path == "" {
		return nil, aws.ErrEmptyParameter
	}

	options := &listByPathOptions{
		withDecryption: true,
	}
	for _, opt := range opts {
		opt(options)
	}

	input := &ssm.GetParametersByPathInput{
		Path:           awssdk.String(path),
		Recursive:      awssdk.Bool(options.recursive),
		WithDecryption: awssdk.Bool(options.withDecryption),
	}

	if options.maxResults > 0 {
		input.MaxResults = awssdk.Int32(options.maxResults)
	}

	var params []ParameterInfo

	for {
		select {
		case <-ctx.Done():
			return params, ctx.Err()
		default:
		}

		output, err := c.api.GetParametersByPath(ctx, input)
		if err != nil {
			return nil, aws.WrapError(serviceName, "GetParametersByPath", err)
		}

		for _, p := range output.Parameters {
			info := ParameterInfo{
				Name:     awssdk.ToString(p.Name),
				Type:     string(p.Type),
				Value:    awssdk.ToString(p.Value),
				Version:  p.Version,
				ARN:      awssdk.ToString(p.ARN),
				DataType: awssdk.ToString(p.DataType),
			}

			if p.LastModifiedDate != nil {
				info.LastModifiedDate = *p.LastModifiedDate
			}

			params = append(params, info)

			if options.maxResults > 0 && len(params) >= int(options.maxResults) {
				return params, nil
			}
		}

		if output.NextToken == nil {
			break
		}

		input.NextToken = output.NextToken
	}

	return params, nil
}

// ListParameters returns metadata for all parameters.
//
// Example:
//
//	params, err := client.ListParameters(ctx)
func (c *Client) ListParameters(ctx context.Context) ([]string, error) {
	var names []string

	input := &ssm.DescribeParametersInput{}

	for {
		select {
		case <-ctx.Done():
			return names, ctx.Err()
		default:
		}

		output, err := c.api.DescribeParameters(ctx, input)
		if err != nil {
			return nil, aws.WrapError(serviceName, "DescribeParameters", err)
		}

		for _, p := range output.Parameters {
			names = append(names, awssdk.ToString(p.Name))
		}

		if output.NextToken == nil {
			break
		}

		input.NextToken = output.NextToken
	}

	return names, nil
}

// PutParameterOption configures PutParameter operations.
type PutParameterOption func(*putParameterOptions)

type putParameterOptions struct {
	paramType   types.ParameterType
	overwrite   bool
	description string
	keyID       string
}

// WithParameterType sets the parameter type.
// Default is String.
//
// Example:
//
//	err := client.PutParameter(ctx, "/app/secret", "value", ssm.WithParameterType(types.ParameterTypeSecureString))
func WithParameterType(t types.ParameterType) PutParameterOption {
	return func(o *putParameterOptions) {
		o.paramType = t
	}
}

// WithOverwrite allows overwriting an existing parameter.
//
// Example:
//
//	err := client.PutParameter(ctx, "/app/config", "new-value", ssm.WithOverwrite(true))
func WithOverwrite(overwrite bool) PutParameterOption {
	return func(o *putParameterOptions) {
		o.overwrite = overwrite
	}
}

// WithDescription sets the parameter description.
//
// Example:
//
//	err := client.PutParameter(ctx, "/app/config", "value", ssm.WithDescription("Database URL"))
func WithDescription(desc string) PutParameterOption {
	return func(o *putParameterOptions) {
		o.description = desc
	}
}

// WithKMSKeyID sets the KMS key for SecureString encryption.
//
// Example:
//
//	err := client.PutParameter(ctx, "/app/secret", "value",
//	    ssm.WithParameterType(types.ParameterTypeSecureString),
//	    ssm.WithKMSKeyID("alias/my-key"),
//	)
func WithKMSKeyID(keyID string) PutParameterOption {
	return func(o *putParameterOptions) {
		o.keyID = keyID
	}
}

// PutParameter creates or updates a parameter.
//
// Example:
//
//	err := client.PutParameter(ctx, "/app/config/database_url", "postgres://...")
func (c *Client) PutParameter(ctx context.Context, name, value string, opts ...PutParameterOption) error {
	if name == "" {
		return aws.ErrEmptyParameter
	}

	options := &putParameterOptions{
		paramType: types.ParameterTypeString,
	}
	for _, opt := range opts {
		opt(options)
	}

	input := &ssm.PutParameterInput{
		Name:      awssdk.String(name),
		Value:     awssdk.String(value),
		Type:      options.paramType,
		Overwrite: awssdk.Bool(options.overwrite),
	}

	if options.description != "" {
		input.Description = awssdk.String(options.description)
	}

	if options.keyID != "" {
		input.KeyId = awssdk.String(options.keyID)
	}

	_, err := c.api.PutParameter(ctx, input)
	if err != nil {
		return aws.WrapError(serviceName, "PutParameter", err)
	}

	return nil
}

// DeleteParameter deletes a parameter.
//
// Example:
//
//	err := client.DeleteParameter(ctx, "/app/config/old-setting")
func (c *Client) DeleteParameter(ctx context.Context, name string) error {
	if name == "" {
		return aws.ErrEmptyParameter
	}

	_, err := c.api.DeleteParameter(ctx, &ssm.DeleteParameterInput{
		Name: awssdk.String(name),
	})
	if err != nil {
		if isParameterNotFound(err) {
			return ErrParameterNotFound
		}

		return aws.WrapError(serviceName, "DeleteParameter", err)
	}

	return nil
}
