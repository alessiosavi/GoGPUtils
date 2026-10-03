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

	if len(output.Item) == 0 {
		return ErrItemNotFound
	}

	if err = attributevalue.UnmarshalMap(output.Item, dest); err != nil {
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

	if len(output.Item) == 0 {
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

// BatchWriteItems writes multiple items to a table in batches of up to 25.
// It returns unprocessed items reported by successful batches. If a later batch
// fails to marshal or execute, or an unprocessed item fails to decode, it returns
// the items decoded so far together with the error. The failed batch and unsent
// items are not included; an execution error may leave the batch's outcome unknown.
// Numbers are decoded as float64, so integers above 2^53 can lose precision.
// Use BatchWriteItemsRaw for lossless retries without decoding.
//
// Example:
//
//	items := []any{user1, user2, user3}
//	unprocessed, err := client.BatchWriteItems(ctx, "users", items)
func (c *Client) BatchWriteItems(ctx context.Context, tableName string, items []any) ([]any, error) {
	return batchItems(ctx, c, tableName, items, false, attributevalue.MarshalMap,
		func(av map[string]types.AttributeValue) (any, error) {
			var item any
			err := attributevalue.UnmarshalMap(av, &item)

			return item, err
		})
}

// BatchDeleteItems deletes multiple items from a table in batches of up to 25.
// It returns unprocessed keys reported by successful batches. If a later batch
// fails to marshal or execute, or an unprocessed key fails to decode, it returns
// the keys decoded so far together with the error. The failed batch and unsent
// keys are not included; an execution error may leave the batch's outcome unknown.
// Numbers are decoded as float64, so integers above 2^53 can lose precision.
// Use BatchDeleteItemsRaw for lossless retries without decoding.
//
// Example:
//
//	keys := []dynamodb.Key{{"pk": "user-1"}, {"pk": "user-2"}}
//	unprocessed, err := client.BatchDeleteItems(ctx, "users", keys)
func (c *Client) BatchDeleteItems(ctx context.Context, tableName string, keys []Key) ([]Key, error) {
	return batchItems(ctx, c, tableName, keys, true,
		func(key Key) (map[string]types.AttributeValue, error) {
			return attributevalue.MarshalMap(key)
		},
		func(av map[string]types.AttributeValue) (Key, error) {
			var key Key
			err := attributevalue.UnmarshalMap(av, &key)

			return key, err
		})
}

// BatchWriteItemsRaw writes raw DynamoDB items in batches of up to 25.
// It returns unprocessed items without marshaling or decoding, preserving number
// precision and all attribute types. Returned maps may alias the SDK response.
// If a later batch fails, it returns previously accumulated unprocessed items
// together with the error. The failed batch and unsent items are not included;
// an execution error may leave the failed batch's outcome unknown.
//
// Example:
//
//	items := []map[string]types.AttributeValue{
//	    {"pk": &types.AttributeValueMemberN{Value: "9007199254740993"}},
//	}
//	unprocessed, err := client.BatchWriteItemsRaw(ctx, "users", items)
func (c *Client) BatchWriteItemsRaw(ctx context.Context, tableName string, items []map[string]types.AttributeValue) ([]map[string]types.AttributeValue, error) {
	return batchItems(ctx, c, tableName, items, false, rawAttributes, rawAttributes)
}

// BatchDeleteItemsRaw deletes raw DynamoDB keys in batches of up to 25.
// It returns unprocessed keys without marshaling or decoding, preserving number
// precision and all attribute types. Returned maps may alias the SDK response.
// If a later batch fails, it returns previously accumulated unprocessed keys
// together with the error. The failed batch and unsent keys are not included;
// an execution error may leave the failed batch's outcome unknown.
//
// Example:
//
//	keys := []map[string]types.AttributeValue{
//	    {"pk": &types.AttributeValueMemberN{Value: "9007199254740993"}},
//	}
//	unprocessed, err := client.BatchDeleteItemsRaw(ctx, "users", keys)
func (c *Client) BatchDeleteItemsRaw(ctx context.Context, tableName string, keys []map[string]types.AttributeValue) ([]map[string]types.AttributeValue, error) {
	return batchItems(ctx, c, tableName, keys, true, rawAttributes, rawAttributes)
}

func rawAttributes(av map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	return av, nil
}

// batchItems shares chunking and partial-result handling across batch APIs.
func batchItems[T any](ctx context.Context, c *Client, tableName string, items []T, deleteItems bool,
	marshal func(T) (map[string]types.AttributeValue, error),
	unmarshal func(map[string]types.AttributeValue) (T, error),
) ([]T, error) {
	if tableName == "" {
		return nil, aws.ErrEmptyTable
	}

	const maxBatchSize = 25

	var unprocessed []T

	for i := 0; i < len(items); i += maxBatchSize {
		end := min(i+maxBatchSize, len(items))
		writeRequests := make([]types.WriteRequest, 0, end-i)

		for _, item := range items[i:end] {
			av, err := marshal(item)
			if err != nil {
				return unprocessed, aws.WrapError(serviceName, "BatchWriteItem", err)
			}

			req := types.WriteRequest{}
			if deleteItems {
				req.DeleteRequest = &types.DeleteRequest{Key: av}
			} else {
				req.PutRequest = &types.PutRequest{Item: av}
			}

			writeRequests = append(writeRequests, req)
		}

		output, err := c.api.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{tableName: writeRequests},
		})
		if err != nil {
			return unprocessed, aws.WrapError(serviceName, "BatchWriteItem", err)
		}

		for _, requests := range output.UnprocessedItems {
			for _, req := range requests {
				var av map[string]types.AttributeValue
				if req.PutRequest != nil {
					av = req.PutRequest.Item
				} else if req.DeleteRequest != nil {
					av = req.DeleteRequest.Key
				}

				item, err := unmarshal(av)
				if err != nil {
					return unprocessed, aws.WrapError(serviceName, "BatchWriteItem", err)
				}

				unprocessed = append(unprocessed, item)
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
