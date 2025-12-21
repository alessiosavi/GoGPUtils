package dynamodb

import (
	"context"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Key represents a DynamoDB item key.
// Keys should use string, number, or binary attribute values.
//
// Example:
//
//	// Simple partition key
//	:= dynamodb.Key{"pk": "user-123"}
//
//	// Composite key (partition + sort)
//	key := dynamodb.Key{"pk": "user-123", "sk": "profile"}
type Key map[string]any

// GetItem retrieves an item from a DynamoDB table.
// The result is unmarshaled into the provided destination.
//
// Example:
//
//	var user User
//	err := client.GetItem(ctx, "users", dynamodb.Key{"pk": "user-123"}, &user)
//	if errors.Is(err, dynamodb.ErrItemNotFound) {
//	    // Handle not found
//	}
func (c *Client) GetItem(ctx context.Context, tableName string, key Key, dest any) error {
	if tableName == "" {
		return aws.ErrEmptyTable
	}

	if len(key) == 0 {
		return aws.ErrEmptyKey
	}

	keyAV, err := marshalKey(key)
	if err != nil {
		return err
	}

	output, err := c.api.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: awssdk.String(tableName),
		Key:       keyAV,
	})
	if err != nil {
		if isResourceNotFound(err) {
			return ErrTableNotFound
		}

		return aws.WrapError(serviceName, "GetItem", err)
	}

	if output.Item == nil || len(output.Item) == 0 {
		return ErrItemNotFound
	}

	if err := attributevalue.UnmarshalMap(output.Item, dest); err != nil {
		return aws.WrapError(serviceName, "GetItem", err)
	}

	return nil
}

// GetItemRaw retrieves an item and returns the raw attribute value map.
// Useful when you don't want to unmarshal into a struct.
//
// Example:
//
//	item, err := client.GetItemRaw(ctx, "users", dynamodb.Key{"pk": "user-123"})
func (c *Client) GetItemRaw(ctx context.Context, tableName string, key Key) (map[string]types.AttributeValue, error) {
	if tableName == "" {
		return nil, aws.ErrEmptyTable
	}

	if len(key) == 0 {
		return nil, aws.ErrEmptyKey
	}

	keyAV, err := marshalKey(key)
	if err != nil {
		return nil, err
	}

	output, err := c.api.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: awssdk.String(tableName),
		Key:       keyAV,
	})
	if err != nil {
		if isResourceNotFound(err) {
			return nil, ErrTableNotFound
		}

		return nil, aws.WrapError(serviceName, "GetItem", err)
	}

	if output.Item == nil || len(output.Item) == 0 {
		return nil, ErrItemNotFound
	}

	return output.Item, nil
}

// PutItem writes an item to a DynamoDB table.
// The item is automatically marshaled from the provided struct.
//
// Example:
//
//	user := User{ID: "user-123", Email: "alice@example.com", Name: "Alice"}
//	err := client.PutItem(ctx, "users", user)
func (c *Client) PutItem(ctx context.Context, tableName string, item any) error {
	if tableName == "" {
		return aws.ErrEmptyTable
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return aws.WrapError(serviceName, "PutItem", err)
	}

	_, err = c.api.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: awssdk.String(tableName),
		Item:      av,
	})
	if err != nil {
		if isResourceNotFound(err) {
			return ErrTableNotFound
		}

		return aws.WrapError(serviceName, "PutItem", err)
	}

	return nil
}

// PutItemIfNotExists writes an item only if it doesn't already exist.
// Returns ErrConditionalCheckFailed if the item exists.
//
// Example:
//
//	user := User{ID: "user-123", Email: "alice@example.com"}
//	err := client.PutItemIfNotExists(ctx, "users", user, "pk")
//	if errors.Is(err, dynamodb.ErrConditionalCheckFailed) {
//	    // Item already exists
//	}
func (c *Client) PutItemIfNotExists(ctx context.Context, tableName string, item any, pkAttribute string) error {
	if tableName == "" {
		return aws.ErrEmptyTable
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return aws.WrapError(serviceName, "PutItem", err)
	}

	_, err = c.api.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           awssdk.String(tableName),
		Item:                av,
		ConditionExpression: awssdk.String("attribute_not_exists(" + pkAttribute + ")"),
	})
	if err != nil {
		if isConditionalCheckFailed(err) {
			return ErrConditionalCheckFailed
		}

		if isResourceNotFound(err) {
			return ErrTableNotFound
		}

		return aws.WrapError(serviceName, "PutItem", err)
	}

	return nil
}

// DeleteItem deletes an item from a DynamoDB table.
//
// Example:
//
//	err := client.DeleteItem(ctx, "users", dynamodb.Key{"pk": "user-123"})
func (c *Client) DeleteItem(ctx context.Context, tableName string, key Key) error {
	if tableName == "" {
		return aws.ErrEmptyTable
	}

	if len(key) == 0 {
		return aws.ErrEmptyKey
	}

	keyAV, err := marshalKey(key)
	if err != nil {
		return err
	}

	_, err = c.api.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: awssdk.String(tableName),
		Key:       keyAV,
	})
	if err != nil {
		if isResourceNotFound(err) {
			return ErrTableNotFound
		}

		return aws.WrapError(serviceName, "DeleteItem", err)
	}

	return nil
}

// DeleteItemIfExists deletes an item only if it exists.
// Returns ErrConditionalCheckFailed if the item doesn't exist.
//
// Example:
//
//	err := client.DeleteItemIfExists(ctx, "users", dynamodb.Key{"pk": "user-123"}, "pk")
func (c *Client) DeleteItemIfExists(ctx context.Context, tableName string, key Key, pkAttribute string) error {
	if tableName == "" {
		return aws.ErrEmptyTable
	}

	if len(key) == 0 {
		return aws.ErrEmptyKey
	}

	keyAV, err := marshalKey(key)
	if err != nil {
		return err
	}

	_, err = c.api.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:           awssdk.String(tableName),
		Key:                 keyAV,
		ConditionExpression: awssdk.String("attribute_exists(" + pkAttribute + ")"),
	})
	if err != nil {
		if isConditionalCheckFailed(err) {
			return ErrConditionalCheckFailed
		}

		if isResourceNotFound(err) {
			return ErrTableNotFound
		}

		return aws.WrapError(serviceName, "DeleteItem", err)
	}

	return nil
}

// BatchWriteItems writes multiple items to one or more tables.
// Automatically handles the 25-item limit per batch.
//
// Example:
//
//	items := []User{user1, user2, user3}
//	unprocessed, err := client.BatchWriteItems(ctx, "users", items)
func (c *Client) BatchWriteItems(ctx context.Context, tableName string, items []any) ([]any, error) {
	if tableName == "" {
		return nil, aws.ErrEmptyTable
	}

	if len(items) == 0 {
		return nil, nil
	}

	const maxBatchSize = 25

	var unprocessed []any

	for i := 0; i < len(items); i += maxBatchSize {
		end := min(i+maxBatchSize, len(items))

		batch := items[i:end]
		writeRequests := make([]types.WriteRequest, 0, len(batch))

		for _, item := range batch {
			av, err := attributevalue.MarshalMap(item)
			if err != nil {
				return nil, aws.WrapError(serviceName, "BatchWriteItem", err)
			}

			writeRequests = append(writeRequests, types.WriteRequest{
				PutRequest: &types.PutRequest{Item: av},
			})
		}

		requestItems := map[string][]types.WriteRequest{tableName: writeRequests}

		output, err := c.api.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: requestItems,
		})
		if err != nil {
			return nil, aws.WrapError(serviceName, "BatchWriteItem", err)
		}

		// Handle unprocessed items
		if len(output.UnprocessedItems) > 0 {
			for _, reqs := range output.UnprocessedItems {
				for _, req := range reqs {
					if req.PutRequest != nil {
						var item any
						err := attributevalue.UnmarshalMap(req.PutRequest.Item, &item)
						if err == nil {
							unprocessed = append(unprocessed, item)
						}
					}
				}
			}
		}
	}

	return unprocessed, nil
}

// BatchDeleteItems deletes multiple items from a table.
// Automatically handles the 25-item limit per batch.
//
// Example:
//
//	keys := []dynamodb.Key{{"pk": "user-1"}, {"pk": "user-2"}}
//	unprocessed, err := client.BatchDeleteItems(ctx, "users", keys)
func (c *Client) BatchDeleteItems(ctx context.Context, tableName string, keys []Key) ([]Key, error) {
	if tableName == "" {
		return nil, aws.ErrEmptyTable
	}

	if len(keys) == 0 {
		return nil, nil
	}

	const maxBatchSize = 25

	var unprocessed []Key

	for i := 0; i < len(keys); i += maxBatchSize {
		end := min(i+maxBatchSize, len(keys))

		batch := keys[i:end]
		writeRequests := make([]types.WriteRequest, 0, len(batch))

		for _, key := range batch {
			keyAV, err := marshalKey(key)
			if err != nil {
				return nil, err
			}

			writeRequests = append(writeRequests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{Key: keyAV},
			})
		}

		requestItems := map[string][]types.WriteRequest{tableName: writeRequests}

		output, err := c.api.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: requestItems,
		})
		if err != nil {
			return nil, aws.WrapError(serviceName, "BatchWriteItem", err)
		}

		// Handle unprocessed items
		if len(output.UnprocessedItems) > 0 {
			for _, reqs := range output.UnprocessedItems {
				for _, req := range reqs {
					if req.DeleteRequest != nil {
						key := make(Key)
						err := attributevalue.UnmarshalMap(req.DeleteRequest.Key, &key)
						if err == nil {
							unprocessed = append(unprocessed, key)
						}
					}
				}
			}
		}
	}

	return unprocessed, nil
}

// marshalKey marshals a Key to DynamoDB attribute values.
func marshalKey(key Key) (map[string]types.AttributeValue, error) {
	result, err := attributevalue.MarshalMap(key)
	if err != nil {
		return nil, aws.WrapError(serviceName, "marshalKey", err)
	}

	return result, nil
}
