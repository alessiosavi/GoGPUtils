package lambda

import (
	"context"
	"encoding/json"
	"os"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// InvokeResult contains the result of a Lambda invocation.
type InvokeResult struct {
	StatusCode      int32
	Payload         []byte
	FunctionError   string
	LogResult       string
	ExecutedVersion string
}

// Unmarshal unmarshals the payload into the provided destination.
//
// Example:
//
//	result, err := client.Invoke(ctx, "my-function", input)
//	var response MyResponse
//	err = result.Unmarshal(&response)
func (r *InvokeResult) Unmarshal(dest any) error {
	return json.Unmarshal(r.Payload, dest)
}

// HasError returns true if the function returned an error.
func (r *InvokeResult) HasError() bool {
	return r.FunctionError != ""
}

// InvokeOption configures Invoke operations.
type InvokeOption func(*invokeOptions)

type invokeOptions struct {
	invocationType types.InvocationType
	logType        types.LogType
	qualifier      string
}

// WithInvocationType sets the invocation type.
// Default is RequestResponse (synchronous).
//
// Example:
//
//	result, err := client.Invoke(ctx, "my-function", payload, lambda.WithInvocationType(types.InvocationTypeEvent))
func WithInvocationType(t types.InvocationType) InvokeOption {
	return func(o *invokeOptions) {
		o.invocationType = t
	}
}

// WithLogType sets the log type.
// Use LogTypeTail to include the last 4KB of logs in the response.
//
// Example:
//
//	result, err := client.Invoke(ctx, "my-function", payload, lambda.WithLogType(types.LogTypeTail))
//	fmt.Println("Logs:", result.LogResult)
func WithLogType(t types.LogType) InvokeOption {
	return func(o *invokeOptions) {
		o.logType = t
	}
}

// WithQualifier sets the function version or alias.
//
// Example:
//
//	result, err := client.Invoke(ctx, "my-function", payload, lambda.WithQualifier("prod"))
func WithQualifier(qualifier string) InvokeOption {
	return func(o *invokeOptions) {
		o.qualifier = qualifier
	}
}

// Invoke calls a Lambda function synchronously and returns the result.
//
// Example:
//
//	payload := []byte(`{"name": "John"}`)
//	result, err := client.Invoke(ctx, "my-function", payload)
//	if err != nil {
//	    return err
//	}
//	if result.HasError() {
//	    return fmt.Errorf("function error: %s", result.FunctionError)
//	}
//	fmt.Println(string(result.Payload))
func (c *Client) Invoke(ctx context.Context, functionName string, payload []byte, opts ...InvokeOption) (*InvokeResult, error) {
	if functionName == "" {
		return nil, aws.ErrEmptyFunction
	}

	options := &invokeOptions{
		invocationType: types.InvocationTypeRequestResponse,
	}
	for _, opt := range opts {
		opt(options)
	}

	input := &lambda.InvokeInput{
		FunctionName:   awssdk.String(functionName),
		InvocationType: options.invocationType,
		Payload:        payload,
	}

	if options.logType != "" {
		input.LogType = options.logType
	}

	if options.qualifier != "" {
		input.Qualifier = awssdk.String(options.qualifier)
	}

	output, err := c.api.Invoke(ctx, input)
	if err != nil {
		if isResourceNotFound(err) {
			return nil, ErrFunctionNotFound
		}

		return nil, aws.WrapError(serviceName, "Invoke", err)
	}

	result := &InvokeResult{
		StatusCode:      output.StatusCode,
		Payload:         output.Payload,
		ExecutedVersion: awssdk.ToString(output.ExecutedVersion),
	}

	if output.FunctionError != nil {
		result.FunctionError = *output.FunctionError
	}

	if output.LogResult != nil {
		result.LogResult = *output.LogResult
	}

	return result, nil
}

// InvokeJSON calls a Lambda function with a JSON payload.
// The input is automatically marshaled to JSON.
//
// Example:
//
//	input := MyInput{Name: "John", Age: 30}
//	result, err := client.InvokeJSON(ctx, "my-function", input)
func (c *Client) InvokeJSON(ctx context.Context, functionName string, input any, opts ...InvokeOption) (*InvokeResult, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, aws.WrapError(serviceName, "InvokeJSON", err)
	}

	return c.Invoke(ctx, functionName, payload, opts...)
}

// InvokeAsync calls a Lambda function asynchronously (fire and forget).
//
// Example:
//
//	err := client.InvokeAsync(ctx, "my-function", payload)
func (c *Client) InvokeAsync(ctx context.Context, functionName string, payload []byte) error {
	_, err := c.Invoke(ctx, functionName, payload, WithInvocationType(types.InvocationTypeEvent))

	return err
}

// FunctionInfo contains information about a Lambda function.
type FunctionInfo struct {
	Name          string
	ARN           string
	Description   string
	Runtime       string
	Handler       string
	MemorySize    int32
	Timeout       int32
	LastModified  string
	CodeSize      int64
	State         string
	StateReason   string
	Architectures []string
}

// ListFunctions returns all Lambda functions in the account.
//
// Example:
//
//	functions, err := client.ListFunctions(ctx)
//	for _, fn := range functions {
//	    fmt.Printf("%s (%s)\n", fn.Name, fn.Runtime)
//	}
func (c *Client) ListFunctions(ctx context.Context) ([]FunctionInfo, error) {
	var functions []FunctionInfo

	input := &lambda.ListFunctionsInput{}

	for {
		select {
		case <-ctx.Done():
			return functions, ctx.Err()
		default:
		}

		output, err := c.api.ListFunctions(ctx, input)
		if err != nil {
			return nil, aws.WrapError(serviceName, "ListFunctions", err)
		}

		for _, fn := range output.Functions {
			archs := make([]string, len(fn.Architectures))
			for i, a := range fn.Architectures {
				archs[i] = string(a)
			}

			functions = append(functions, FunctionInfo{
				Name:          awssdk.ToString(fn.FunctionName),
				ARN:           awssdk.ToString(fn.FunctionArn),
				Description:   awssdk.ToString(fn.Description),
				Runtime:       string(fn.Runtime),
				Handler:       awssdk.ToString(fn.Handler),
				MemorySize:    awssdk.ToInt32(fn.MemorySize),
				Timeout:       awssdk.ToInt32(fn.Timeout),
				LastModified:  awssdk.ToString(fn.LastModified),
				CodeSize:      fn.CodeSize,
				State:         string(fn.State),
				StateReason:   awssdk.ToString(fn.StateReason),
				Architectures: archs,
			})
		}

		if output.NextMarker == nil {
			break
		}

		input.Marker = output.NextMarker
	}

	return functions, nil
}

// ListFunctionNames returns just the names of all Lambda functions.
//
// Example:
//
//	names, err := client.ListFunctionNames(ctx)
func (c *Client) ListFunctionNames(ctx context.Context) ([]string, error) {
	functions, err := c.ListFunctions(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(functions))
	for i, fn := range functions {
		names[i] = fn.Name
	}

	return names, nil
}

// GetFunction retrieves information about a Lambda function.
//
// Example:
//
//	info, err := client.GetFunction(ctx, "my-function")
//	fmt.Printf("Runtime: %s, Memory: %dMB\n", info.Runtime, info.MemorySize)
func (c *Client) GetFunction(ctx context.Context, functionName string) (*FunctionInfo, error) {
	if functionName == "" {
		return nil, aws.ErrEmptyFunction
	}

	output, err := c.api.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: awssdk.String(functionName),
	})
	if err != nil {
		if isResourceNotFound(err) {
			return nil, ErrFunctionNotFound
		}

		return nil, aws.WrapError(serviceName, "GetFunction", err)
	}

	fn := output.Configuration

	archs := make([]string, len(fn.Architectures))
	for i, a := range fn.Architectures {
		archs[i] = string(a)
	}

	return &FunctionInfo{
		Name:          awssdk.ToString(fn.FunctionName),
		ARN:           awssdk.ToString(fn.FunctionArn),
		Description:   awssdk.ToString(fn.Description),
		Runtime:       string(fn.Runtime),
		Handler:       awssdk.ToString(fn.Handler),
		MemorySize:    awssdk.ToInt32(fn.MemorySize),
		Timeout:       awssdk.ToInt32(fn.Timeout),
		LastModified:  awssdk.ToString(fn.LastModified),
		CodeSize:      fn.CodeSize,
		State:         string(fn.State),
		StateReason:   awssdk.ToString(fn.StateReason),
		Architectures: archs,
	}, nil
}

// DeployFromS3 updates a function's code from an S3 object.
//
// Example:
//
//	err := client.DeployFromS3(ctx, "my-function", "my-bucket", "deployments/code.zip")
func (c *Client) DeployFromS3(ctx context.Context, functionName, bucket, key string) error {
	if functionName == "" {
		return aws.ErrEmptyFunction
	}

	if bucket == "" {
		return aws.ErrEmptyBucket
	}

	if key == "" {
		return aws.ErrEmptyKey
	}

	_, err := c.api.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		FunctionName: awssdk.String(functionName),
		S3Bucket:     awssdk.String(bucket),
		S3Key:        awssdk.String(key),
		Publish:      true,
	})
	if err != nil {
		if isResourceNotFound(err) {
			return ErrFunctionNotFound
		}

		return aws.WrapError(serviceName, "UpdateFunctionCode", err)
	}

	return nil
}

// DeployFromZip updates a function's code from a ZIP file.
//
// Example:
//
//	err := client.DeployFromZip(ctx, "my-function", "function.zip")
func (c *Client) DeployFromZip(ctx context.Context, functionName, zipPath string) error {
	if functionName == "" {
		return aws.ErrEmptyFunction
	}

	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		return aws.WrapError(serviceName, "DeployFromZip", err)
	}

	_, err = c.api.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		FunctionName: awssdk.String(functionName),
		ZipFile:      zipData,
		Publish:      true,
	})
	if err != nil {
		if isResourceNotFound(err) {
			return ErrFunctionNotFound
		}

		return aws.WrapError(serviceName, "UpdateFunctionCode", err)
	}

	return nil
}

// DeployFromBytes updates a function's code from a ZIP byte slice.
//
// Example:
//
//	zipData, _ := os.ReadFile("function.zip")
//	err := client.DeployFromBytes(ctx, "my-function", zipData)
func (c *Client) DeployFromBytes(ctx context.Context, functionName string, zipData []byte) error {
	if functionName == "" {
		return aws.ErrEmptyFunction
	}

	_, err := c.api.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		FunctionName: awssdk.String(functionName),
		ZipFile:      zipData,
		Publish:      true,
	})
	if err != nil {
		if isResourceNotFound(err) {
			return ErrFunctionNotFound
		}

		return aws.WrapError(serviceName, "UpdateFunctionCode", err)
	}

	return nil
}

// DeleteFunction deletes a Lambda function.
//
// Example:
//
//	err := client.DeleteFunction(ctx, "my-function")
func (c *Client) DeleteFunction(ctx context.Context, functionName string) error {
	if functionName == "" {
		return aws.ErrEmptyFunction
	}

	_, err := c.api.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
		FunctionName: awssdk.String(functionName),
	})
	if err != nil {
		if isResourceNotFound(err) {
			return ErrFunctionNotFound
		}

		return aws.WrapError(serviceName, "DeleteFunction", err)
	}

	return nil
}

// GetTags retrieves the tags for a Lambda function.
//
// Example:
//
//	tags, err := client.GetTags(ctx, "arn:aws:lambda:us-west-2:123456789:function:my-function")
func (c *Client) GetTags(ctx context.Context, functionARN string) (map[string]string, error) {
	if functionARN == "" {
		return nil, aws.ErrEmptyFunction
	}

	output, err := c.api.ListTags(ctx, &lambda.ListTagsInput{
		Resource: awssdk.String(functionARN),
	})
	if err != nil {
		if isResourceNotFound(err) {
			return nil, ErrFunctionNotFound
		}

		return nil, aws.WrapError(serviceName, "ListTags", err)
	}

	return output.Tags, nil
}
