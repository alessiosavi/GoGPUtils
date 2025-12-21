package dynamodb

import (
	"context"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// QueryOption configures a Query operation.
type QueryOption func(*queryOptions)

type queryOptions struct {
	indexName            string
	limit                int32
	scanForward          *bool
	filterExpression     expression.ConditionBuilder
	projectionExpression []string
	consistentRead       bool
	exclusiveStartKey    map[string]types.AttributeValue
}

// WithIndex specifies a secondary index to query.
//
// Example:
//
//	results, err := client.Query(ctx, "users", keyExpr, dynamodb.WithIndex("email-index"))
func WithIndex(indexName string) QueryOption {
	return func(o *queryOptions) {
		o.indexName = indexName
	}
}

// WithLimit limits the number of items returned.
//
// Example:
//
//	results, err := client.Query(ctx, "users", keyExpr, dynamodb.WithLimit(10))
func WithLimit(limit int32) QueryOption {
	return func(o *queryOptions) {
		o.limit = limit
	}
}

// WithScanForward sets the scan direction.
// true = ascending, false = descending.
//
// Example:
//
//	results, err := client.Query(ctx, "users", keyExpr, dynamodb.WithScanForward(false))
func WithScanForward(forward bool) QueryOption {
	return func(o *queryOptions) {
		o.scanForward = &forward
	}
}

// WithFilter adds a filter expression to the query.
//
// Example:
//
//	filter := expression.Name("status").Equal(expression.Value("active"))
//	results, err := client.Query(ctx, "users", keyExpr, dynamodb.WithFilter(filter))
func WithFilter(filter expression.ConditionBuilder) QueryOption {
	return func(o *queryOptions) {
		o.filterExpression = filter
	}
}

// WithProjection specifies which attributes to return.
//
// Example:
//
//	results, err := client.Query(ctx, "users", keyExpr, dynamodb.WithProjection("id", "name", "email"))
func WithProjection(attrs ...string) QueryOption {
	return func(o *queryOptions) {
		o.projectionExpression = attrs
	}
}

// WithConsistentRead enables consistent reads.
//
// Example:
//
//	results, err := client.Query(ctx, "users", keyExpr, dynamodb.WithConsistentRead())
func WithConsistentRead() QueryOption {
	return func(o *queryOptions) {
		o.consistentRead = true
	}
}

// WithStartKey sets the exclusive start key for pagination.
//
// Example:
//
//	results, err := client.Query(ctx, "users", keyExpr, dynamodb.WithStartKey(lastKey))
func WithStartKey(key map[string]types.AttributeValue) QueryOption {
	return func(o *queryOptions) {
		o.exclusiveStartKey = key
	}
}

// QueryResult contains the results of a Query operation.
type QueryResult struct {
	Items            []map[string]types.AttributeValue
	Count            int32
	ScannedCount     int32
	LastEvaluatedKey map[string]types.AttributeValue
}

// Unmarshal unmarshals the query results into a slice of the provided type.
//
// Example:
//
//	result, err := client.Query(ctx, "users", keyExpr)
//	var users []User
//	err = result.Unmarshal(&users)
func (r *QueryResult) Unmarshal(dest any) error {
	return attributevalue.UnmarshalListOfMaps(r.Items, dest)
}

// HasMorePages returns true if there are more results to fetch.
func (r *QueryResult) HasMorePages() bool {
	return len(r.LastEvaluatedKey) > 0
}

// Query performs a query operation on a DynamoDB table.
//
// Example:
//
//	// Query by partition key
//	keyExpr := expression.Key("pk").Equal(expression.Value("user-123"))
//	result, err := client.Query(ctx, "users", keyExpr)
//
//	// Query with sort key condition
//	keyExpr := expression.KeyAnd(
//	    expression.Key("pk").Equal(expression.Value("user-123")),
//	    expression.Key("sk").BeginsWith("order-"),
//	)
//	result, err := client.Query(ctx, "orders", keyExpr)
func (c *Client) Query(ctx context.Context, tableName string, keyCondition expression.KeyConditionBuilder, opts ...QueryOption) (*QueryResult, error) {
	if tableName == "" {
		return nil, aws.ErrEmptyTable
	}

	options := &queryOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Build expression
	builder := expression.NewBuilder().WithKeyCondition(keyCondition)

	if options.filterExpression.IsSet() {
		builder = builder.WithFilter(options.filterExpression)
	}

	if len(options.projectionExpression) > 0 {
		var names []expression.NameBuilder
		for _, attr := range options.projectionExpression {
			names = append(names, expression.Name(attr))
		}

		proj := expression.ProjectionBuilder{}
		for _, name := range names {
			proj = proj.AddNames(name)
		}

		builder = builder.WithProjection(proj)
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, aws.WrapError(serviceName, "Query", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 awssdk.String(tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConsistentRead:            awssdk.Bool(options.consistentRead),
	}

	if expr.Filter() != nil {
		input.FilterExpression = expr.Filter()
	}

	if expr.Projection() != nil {
		input.ProjectionExpression = expr.Projection()
	}

	if options.indexName != "" {
		input.IndexName = awssdk.String(options.indexName)
	}

	if options.limit > 0 {
		input.Limit = awssdk.Int32(options.limit)
	}

	if options.scanForward != nil {
		input.ScanIndexForward = options.scanForward
	}

	if options.exclusiveStartKey != nil {
		input.ExclusiveStartKey = options.exclusiveStartKey
	}

	output, err := c.api.Query(ctx, input)
	if err != nil {
		if isResourceNotFound(err) {
			return nil, ErrTableNotFound
		}

		return nil, aws.WrapError(serviceName, "Query", err)
	}

	return &QueryResult{
		Items:            output.Items,
		Count:            output.Count,
		ScannedCount:     output.ScannedCount,
		LastEvaluatedKey: output.LastEvaluatedKey,
	}, nil
}

// QueryAll performs a query and automatically paginates through all results.
// Use with caution on large datasets.
//
// Example:
//
//	keyExpr := expression.Key("pk").Equal(expression.Value("user-123"))
//	var users []User
//	err := client.QueryAll(ctx, "users", keyExpr, &users)
func (c *Client) QueryAll(ctx context.Context, tableName string, keyCondition expression.KeyConditionBuilder, dest any, opts ...QueryOption) error {
	var allItems []map[string]types.AttributeValue

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		result, err := c.Query(ctx, tableName, keyCondition, opts...)
		if err != nil {
			return err
		}

		allItems = append(allItems, result.Items...)

		if !result.HasMorePages() {
			break
		}

		// Add start key option for next iteration
		opts = append(opts, WithStartKey(result.LastEvaluatedKey))
	}

	return attributevalue.UnmarshalListOfMaps(allItems, dest)
}

// ScanResult contains the results of a Scan operation.
type ScanResult struct {
	Items            []map[string]types.AttributeValue
	Count            int32
	ScannedCount     int32
	LastEvaluatedKey map[string]types.AttributeValue
}

// Unmarshal unmarshals the scan results into a slice of the provided type.
func (r *ScanResult) Unmarshal(dest any) error {
	return attributevalue.UnmarshalListOfMaps(r.Items, dest)
}

// HasMorePages returns true if there are more results to fetch.
func (r *ScanResult) HasMorePages() bool {
	return len(r.LastEvaluatedKey) > 0
}

// ScanOption configures a Scan operation.
type ScanOption func(*scanOptions)

type scanOptions struct {
	limit                int32
	filterExpression     expression.ConditionBuilder
	projectionExpression []string
	consistentRead       bool
	exclusiveStartKey    map[string]types.AttributeValue
	segment              *int32
	totalSegments        *int32
}

// WithScanLimit limits the number of items scanned.
func WithScanLimit(limit int32) ScanOption {
	return func(o *scanOptions) {
		o.limit = limit
	}
}

// WithScanFilter adds a filter expression to the scan.
func WithScanFilter(filter expression.ConditionBuilder) ScanOption {
	return func(o *scanOptions) {
		o.filterExpression = filter
	}
}

// WithScanProjection specifies which attributes to return.
func WithScanProjection(attrs ...string) ScanOption {
	return func(o *scanOptions) {
		o.projectionExpression = attrs
	}
}

// WithScanStartKey sets the exclusive start key for pagination.
func WithScanStartKey(key map[string]types.AttributeValue) ScanOption {
	return func(o *scanOptions) {
		o.exclusiveStartKey = key
	}
}

// WithParallelScan configures parallel scanning.
//
// Example:
//
//	// Scan segment 0 of 4 total segments
//	result, err := client.Scan(ctx, "users", dynamodb.WithParallelScan(0, 4))
func WithParallelScan(segment, totalSegments int32) ScanOption {
	return func(o *scanOptions) {
		o.segment = &segment
		o.totalSegments = &totalSegments
	}
}

// Scan performs a scan operation on a DynamoDB table.
// Scans read every item in the table - use Query when possible.
//
// Example:
//
//	result, err := client.Scan(ctx, "users")
//	var users []User
//	err = result.Unmarshal(&users)
func (c *Client) Scan(ctx context.Context, tableName string, opts ...ScanOption) (*ScanResult, error) {
	if tableName == "" {
		return nil, aws.ErrEmptyTable
	}

	options := &scanOptions{}
	for _, opt := range opts {
		opt(options)
	}

	input := &dynamodb.ScanInput{
		TableName:      awssdk.String(tableName),
		ConsistentRead: awssdk.Bool(options.consistentRead),
	}

	// Build expression if we have filter or projection
	if options.filterExpression.IsSet() || len(options.projectionExpression) > 0 {
		builder := expression.NewBuilder()

		if options.filterExpression.IsSet() {
			builder = builder.WithFilter(options.filterExpression)
		}

		if len(options.projectionExpression) > 0 {
			var names []expression.NameBuilder
			for _, attr := range options.projectionExpression {
				names = append(names, expression.Name(attr))
			}

			proj := expression.ProjectionBuilder{}
			for _, name := range names {
				proj = proj.AddNames(name)
			}

			builder = builder.WithProjection(proj)
		}

		expr, err := builder.Build()
		if err != nil {
			return nil, aws.WrapError(serviceName, "Scan", err)
		}

		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()

		if expr.Filter() != nil {
			input.FilterExpression = expr.Filter()
		}

		if expr.Projection() != nil {
			input.ProjectionExpression = expr.Projection()
		}
	}

	if options.limit > 0 {
		input.Limit = awssdk.Int32(options.limit)
	}

	if options.exclusiveStartKey != nil {
		input.ExclusiveStartKey = options.exclusiveStartKey
	}

	if options.segment != nil && options.totalSegments != nil {
		input.Segment = options.segment
		input.TotalSegments = options.totalSegments
	}

	output, err := c.api.Scan(ctx, input)
	if err != nil {
		if isResourceNotFound(err) {
			return nil, ErrTableNotFound
		}

		return nil, aws.WrapError(serviceName, "Scan", err)
	}

	return &ScanResult{
		Items:            output.Items,
		Count:            output.Count,
		ScannedCount:     output.ScannedCount,
		LastEvaluatedKey: output.LastEvaluatedKey,
	}, nil
}

// ScanAll performs a scan and automatically paginates through all results.
// Use with extreme caution on large tables.
//
// Example:
//
//	var users []User
//	err := client.ScanAll(ctx, "users", &users)
func (c *Client) ScanAll(ctx context.Context, tableName string, dest any, opts ...ScanOption) error {
	var allItems []map[string]types.AttributeValue

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		result, err := c.Scan(ctx, tableName, opts...)
		if err != nil {
			return err
		}

		allItems = append(allItems, result.Items...)

		if !result.HasMorePages() {
			break
		}

		// Add start key option for next iteration
		opts = append(opts, WithScanStartKey(result.LastEvaluatedKey))
	}

	return attributevalue.UnmarshalListOfMaps(allItems, dest)
}

// ScanCallback iterates through all items in a table and calls a callback for each batch.
// Return an error from the callback to stop scanning.
//
// Example:
//
//	err := client.ScanCallback(ctx, "users", func(items []map[string]types.AttributeValue) error {
//	    var users []User
//	    attributevalue.UnmarshalListOfMaps(items, &users)
//	    for _, user := range users {
//	        process(user)
//	    }
//	    return nil
//	})
func (c *Client) ScanCallback(ctx context.Context, tableName string, callback func([]map[string]types.AttributeValue) error, opts ...ScanOption) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		result, err := c.Scan(ctx, tableName, opts...)
		if err != nil {
			return err
		}

		if len(result.Items) > 0 {
			err := callback(result.Items)
			if err != nil {
				return err
			}
		}

		if !result.HasMorePages() {
			break
		}

		opts = append(opts, WithScanStartKey(result.LastEvaluatedKey))
	}

	return nil
}
