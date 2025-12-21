// Package dynamodb integration tests using LocalStack.
//go:build integration

package dynamodb_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alessiosavi/GoGPUtils/aws/dynamodb"
	"github.com/alessiosavi/GoGPUtils/aws/internal/testutil"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	dynamodbsdk "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// TestItem represents a test item for DynamoDB.
type TestItem struct {
	PK     string `dynamodbav:"pk"`
	SK     string `dynamodbav:"sk"`
	Name   string `dynamodbav:"name"`
	Email  string `dynamodbav:"email"`
	Age    int    `dynamodbav:"age"`
	Status string `dynamodbav:"status"`
}

// testContext returns a context with timeout for tests.
func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)

	return ctx
}

// testConfig holds test configuration.
type testConfig struct {
	cfg    *testutil.TestConfig
	client *dynamodb.Client
	table  string
}

// setupTestTable creates a unique test table and returns cleanup function.
func setupTestTable(t *testing.T, ctx context.Context, cfg *testutil.TestConfig) string {
	t.Helper()

	table := testutil.UniqueTableName()
	dynamoClient := dynamodbsdk.NewFromConfig(cfg.AWSConfig.AWS())

	_, err := dynamoClient.CreateTable(ctx, &dynamodbsdk.CreateTableInput{
		TableName: awssdk.String(table),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: awssdk.String("pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: awssdk.String("sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: awssdk.String("pk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: awssdk.String("sk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	// Wait for table to be active
	waiter := dynamodbsdk.NewTableExistsWaiter(dynamoClient)
	if err := waiter.Wait(ctx, &dynamodbsdk.DescribeTableInput{TableName: awssdk.String(table)}, 30*time.Second); err != nil {
		t.Fatalf("failed to wait for table: %v", err)
	}

	t.Cleanup(func() {
		dynamoClient.DeleteTable(context.Background(), &dynamodbsdk.DeleteTableInput{
			TableName: awssdk.String(table),
		})
	})

	return table
}

// setupTest creates a test environment with DynamoDB client and table.
func setupTest(t *testing.T) *testConfig {
	t.Helper()
	testutil.SkipIfNoLocalStack(t)

	ctx := testContext(t)
	cfg := testutil.MustLoadConfig(t)

	client, err := dynamodb.NewClient(cfg.AWSConfig)
	if err != nil {
		t.Fatalf("failed to create DynamoDB client: %v", err)
	}

	table := setupTestTable(t, ctx, cfg)

	return &testConfig{
		cfg:    cfg,
		client: client,
		table:  table,
	}
}

func TestPutItem(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	item := TestItem{
		PK:     "user-1",
		SK:     "profile",
		Name:   "Alice",
		Email:  "alice@example.com",
		Age:    30,
		Status: "active",
	}

	err := tc.client.PutItem(ctx, tc.table, item)
	if err != nil {
		t.Fatalf("PutItem failed: %v", err)
	}

	// Verify by getting the item
	var got TestItem
	err = tc.client.GetItem(ctx, tc.table, dynamodb.Key{"pk": "user-1", "sk": "profile"}, &got)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}

	if got.Name != item.Name {
		t.Errorf("name mismatch: got %s, want %s", got.Name, item.Name)
	}

	if got.Email != item.Email {
		t.Errorf("email mismatch: got %s, want %s", got.Email, item.Email)
	}
}

func TestPutItemIfNotExists(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	item := TestItem{
		PK:     "user-unique",
		SK:     "profile",
		Name:   "Bob",
		Email:  "bob@example.com",
		Age:    25,
		Status: "active",
	}

	// First put should succeed
	err := tc.client.PutItemIfNotExists(ctx, tc.table, item, "pk")
	if err != nil {
		t.Fatalf("first PutItemIfNotExists failed: %v", err)
	}

	// Second put should fail
	err = tc.client.PutItemIfNotExists(ctx, tc.table, item, "pk")
	if !errors.Is(err, dynamodb.ErrConditionalCheckFailed) {
		t.Errorf("expected ErrConditionalCheckFailed, got: %v", err)
	}
}

func TestGetItem(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create an item
	item := TestItem{
		PK:     "user-get",
		SK:     "profile",
		Name:   "Charlie",
		Email:  "charlie@example.com",
		Age:    35,
		Status: "active",
	}
	if err := tc.client.PutItem(ctx, tc.table, item); err != nil {
		t.Fatalf("setup PutItem failed: %v", err)
	}

	t.Run("existing item", func(t *testing.T) {
		var got TestItem
		err := tc.client.GetItem(ctx, tc.table, dynamodb.Key{"pk": "user-get", "sk": "profile"}, &got)
		if err != nil {
			t.Fatalf("GetItem failed: %v", err)
		}

		if got.Name != item.Name {
			t.Errorf("name mismatch: got %s, want %s", got.Name, item.Name)
		}
	})

	t.Run("non-existent item", func(t *testing.T) {
		var got TestItem
		err := tc.client.GetItem(ctx, tc.table, dynamodb.Key{"pk": "non-existent", "sk": "profile"}, &got)
		if !errors.Is(err, dynamodb.ErrItemNotFound) {
			t.Errorf("expected ErrItemNotFound, got: %v", err)
		}
	})
}

func TestGetItemRaw(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	item := TestItem{
		PK:     "user-raw",
		SK:     "profile",
		Name:   "David",
		Email:  "david@example.com",
		Age:    40,
		Status: "active",
	}
	if err := tc.client.PutItem(ctx, tc.table, item); err != nil {
		t.Fatalf("setup PutItem failed: %v", err)
	}

	raw, err := tc.client.GetItemRaw(ctx, tc.table, dynamodb.Key{"pk": "user-raw", "sk": "profile"})
	if err != nil {
		t.Fatalf("GetItemRaw failed: %v", err)
	}

	if raw == nil {
		t.Error("expected non-nil raw item")
	}

	nameAttr, ok := raw["name"]
	if !ok {
		t.Error("expected 'name' attribute in raw item")
	}

	if nameVal, ok := nameAttr.(*types.AttributeValueMemberS); !ok || nameVal.Value != item.Name {
		t.Errorf("name value mismatch in raw item")
	}
}

func TestDeleteItem(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	item := TestItem{
		PK:     "user-delete",
		SK:     "profile",
		Name:   "Eve",
		Email:  "eve@example.com",
		Age:    28,
		Status: "active",
	}
	if err := tc.client.PutItem(ctx, tc.table, item); err != nil {
		t.Fatalf("setup PutItem failed: %v", err)
	}

	// Verify exists
	var got TestItem
	if err := tc.client.GetItem(ctx, tc.table, dynamodb.Key{"pk": "user-delete", "sk": "profile"}, &got); err != nil {
		t.Fatalf("item should exist: %v", err)
	}

	// Delete
	if err := tc.client.DeleteItem(ctx, tc.table, dynamodb.Key{"pk": "user-delete", "sk": "profile"}); err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}

	// Verify deleted
	err := tc.client.GetItem(ctx, tc.table, dynamodb.Key{"pk": "user-delete", "sk": "profile"}, &got)
	if !errors.Is(err, dynamodb.ErrItemNotFound) {
		t.Errorf("expected ErrItemNotFound after delete, got: %v", err)
	}
}

func TestDeleteItemIfExists(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	item := TestItem{
		PK:     "user-delete-exists",
		SK:     "profile",
		Name:   "Frank",
		Email:  "frank@example.com",
		Age:    33,
		Status: "active",
	}
	if err := tc.client.PutItem(ctx, tc.table, item); err != nil {
		t.Fatalf("setup PutItem failed: %v", err)
	}

	// Delete existing item should succeed
	err := tc.client.DeleteItemIfExists(ctx, tc.table, dynamodb.Key{"pk": "user-delete-exists", "sk": "profile"}, "pk")
	if err != nil {
		t.Fatalf("DeleteItemIfExists failed: %v", err)
	}

	// Delete non-existent item should fail
	err = tc.client.DeleteItemIfExists(ctx, tc.table, dynamodb.Key{"pk": "user-delete-exists", "sk": "profile"}, "pk")
	if !errors.Is(err, dynamodb.ErrConditionalCheckFailed) {
		t.Errorf("expected ErrConditionalCheckFailed, got: %v", err)
	}
}

func TestBatchWriteItems(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	items := []any{
		TestItem{PK: "batch-1", SK: "item", Name: "Item1", Email: "item1@example.com", Age: 20, Status: "active"},
		TestItem{PK: "batch-2", SK: "item", Name: "Item2", Email: "item2@example.com", Age: 21, Status: "active"},
		TestItem{PK: "batch-3", SK: "item", Name: "Item3", Email: "item3@example.com", Age: 22, Status: "active"},
	}

	unprocessed, err := tc.client.BatchWriteItems(ctx, tc.table, items)
	if err != nil {
		t.Fatalf("BatchWriteItems failed: %v", err)
	}

	if len(unprocessed) > 0 {
		t.Errorf("expected no unprocessed items, got %d", len(unprocessed))
	}

	// Verify all items exist
	for _, item := range items {
		testItem := item.(TestItem)
		var got TestItem
		if err := tc.client.GetItem(ctx, tc.table, dynamodb.Key{"pk": testItem.PK, "sk": testItem.SK}, &got); err != nil {
			t.Errorf("item %s not found: %v", testItem.PK, err)
		}
	}
}

func TestBatchDeleteItems(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create items
	items := []any{
		TestItem{PK: "batch-del-1", SK: "item", Name: "Item1", Email: "item1@example.com", Age: 20, Status: "active"},
		TestItem{PK: "batch-del-2", SK: "item", Name: "Item2", Email: "item2@example.com", Age: 21, Status: "active"},
		TestItem{PK: "batch-del-3", SK: "item", Name: "Item3", Email: "item3@example.com", Age: 22, Status: "active"},
	}
	if _, err := tc.client.BatchWriteItems(ctx, tc.table, items); err != nil {
		t.Fatalf("setup BatchWriteItems failed: %v", err)
	}

	// Delete items
	keys := []dynamodb.Key{
		{"pk": "batch-del-1", "sk": "item"},
		{"pk": "batch-del-2", "sk": "item"},
		{"pk": "batch-del-3", "sk": "item"},
	}

	unprocessed, err := tc.client.BatchDeleteItems(ctx, tc.table, keys)
	if err != nil {
		t.Fatalf("BatchDeleteItems failed: %v", err)
	}

	if len(unprocessed) > 0 {
		t.Errorf("expected no unprocessed keys, got %d", len(unprocessed))
	}

	// Verify all items deleted
	for _, key := range keys {
		var got TestItem
		err := tc.client.GetItem(ctx, tc.table, key, &got)
		if !errors.Is(err, dynamodb.ErrItemNotFound) {
			t.Errorf("item %v should be deleted", key)
		}
	}
}

func TestQuery(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create items with same PK but different SK
	items := []any{
		TestItem{PK: "query-user", SK: "order-001", Name: "Order1", Age: 100, Status: "shipped"},
		TestItem{PK: "query-user", SK: "order-002", Name: "Order2", Age: 200, Status: "pending"},
		TestItem{PK: "query-user", SK: "order-003", Name: "Order3", Age: 300, Status: "shipped"},
		TestItem{PK: "other-user", SK: "order-001", Name: "Other", Age: 50, Status: "shipped"},
	}
	if _, err := tc.client.BatchWriteItems(ctx, tc.table, items); err != nil {
		t.Fatalf("setup BatchWriteItems failed: %v", err)
	}

	t.Run("query by partition key", func(t *testing.T) {
		keyExpr := expression.Key("pk").Equal(expression.Value("query-user"))
		result, err := tc.client.Query(ctx, tc.table, keyExpr)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}

		if result.Count != 3 {
			t.Errorf("count mismatch: got %d, want 3", result.Count)
		}
	})

	t.Run("query with sort key begins with", func(t *testing.T) {
		keyExpr := expression.KeyAnd(
			expression.Key("pk").Equal(expression.Value("query-user")),
			expression.Key("sk").BeginsWith("order-"),
		)
		result, err := tc.client.Query(ctx, tc.table, keyExpr)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}

		if result.Count != 3 {
			t.Errorf("count mismatch: got %d, want 3", result.Count)
		}
	})

	t.Run("query with filter", func(t *testing.T) {
		keyExpr := expression.Key("pk").Equal(expression.Value("query-user"))
		filter := expression.Name("status").Equal(expression.Value("shipped"))

		result, err := tc.client.Query(ctx, tc.table, keyExpr, dynamodb.WithFilter(filter))
		if err != nil {
			t.Fatalf("Query with filter failed: %v", err)
		}

		if result.Count != 2 {
			t.Errorf("count mismatch: got %d, want 2", result.Count)
		}
	})

	t.Run("query with limit", func(t *testing.T) {
		keyExpr := expression.Key("pk").Equal(expression.Value("query-user"))
		result, err := tc.client.Query(ctx, tc.table, keyExpr, dynamodb.WithLimit(2))
		if err != nil {
			t.Fatalf("Query with limit failed: %v", err)
		}

		if result.Count != 2 {
			t.Errorf("count mismatch: got %d, want 2", result.Count)
		}
	})

	t.Run("query descending", func(t *testing.T) {
		keyExpr := expression.Key("pk").Equal(expression.Value("query-user"))
		result, err := tc.client.Query(ctx, tc.table, keyExpr, dynamodb.WithScanForward(false))
		if err != nil {
			t.Fatalf("Query descending failed: %v", err)
		}

		var items []TestItem
		if err := result.Unmarshal(&items); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// First item should be order-003 (descending order)
		if items[0].SK != "order-003" {
			t.Errorf("expected first item to be order-003, got %s", items[0].SK)
		}
	})
}

func TestQueryAll(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create many items
	var items []any
	for i := 0; i < 10; i++ {
		items = append(items, TestItem{
			PK:     "queryall-user",
			SK:     "item-" + string(rune('0'+i)),
			Name:   "Item" + string(rune('0'+i)),
			Age:    i * 10,
			Status: "active",
		})
	}
	if _, err := tc.client.BatchWriteItems(ctx, tc.table, items); err != nil {
		t.Fatalf("setup BatchWriteItems failed: %v", err)
	}

	keyExpr := expression.Key("pk").Equal(expression.Value("queryall-user"))
	var results []TestItem

	err := tc.client.QueryAll(ctx, tc.table, keyExpr, &results)
	if err != nil {
		t.Fatalf("QueryAll failed: %v", err)
	}

	if len(results) != 10 {
		t.Errorf("count mismatch: got %d, want 10", len(results))
	}
}

func TestScan(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create items
	items := []any{
		TestItem{PK: "scan-1", SK: "item", Name: "Scan1", Age: 10, Status: "active"},
		TestItem{PK: "scan-2", SK: "item", Name: "Scan2", Age: 20, Status: "inactive"},
		TestItem{PK: "scan-3", SK: "item", Name: "Scan3", Age: 30, Status: "active"},
	}
	if _, err := tc.client.BatchWriteItems(ctx, tc.table, items); err != nil {
		t.Fatalf("setup BatchWriteItems failed: %v", err)
	}

	t.Run("scan all", func(t *testing.T) {
		result, err := tc.client.Scan(ctx, tc.table)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		if result.Count < 3 {
			t.Errorf("expected at least 3 items, got %d", result.Count)
		}
	})

	t.Run("scan with filter", func(t *testing.T) {
		filter := expression.Name("status").Equal(expression.Value("active"))
		result, err := tc.client.Scan(ctx, tc.table, dynamodb.WithScanFilter(filter))
		if err != nil {
			t.Fatalf("Scan with filter failed: %v", err)
		}

		// Should find at least 2 active items
		if result.Count < 2 {
			t.Errorf("expected at least 2 active items, got %d", result.Count)
		}
	})

	t.Run("scan with limit", func(t *testing.T) {
		result, err := tc.client.Scan(ctx, tc.table, dynamodb.WithScanLimit(2))
		if err != nil {
			t.Fatalf("Scan with limit failed: %v", err)
		}

		if result.Count > 2 {
			t.Errorf("expected at most 2 items, got %d", result.Count)
		}
	})
}

func TestScanAll(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create items
	var items []any
	for i := 0; i < 5; i++ {
		items = append(items, TestItem{
			PK:     "scanall-" + string(rune('0'+i)),
			SK:     "item",
			Name:   "ScanAll" + string(rune('0'+i)),
			Age:    i * 10,
			Status: "active",
		})
	}
	if _, err := tc.client.BatchWriteItems(ctx, tc.table, items); err != nil {
		t.Fatalf("setup BatchWriteItems failed: %v", err)
	}

	var results []TestItem
	err := tc.client.ScanAll(ctx, tc.table, &results)
	if err != nil {
		t.Fatalf("ScanAll failed: %v", err)
	}

	if len(results) < 5 {
		t.Errorf("expected at least 5 items, got %d", len(results))
	}
}

func TestScanCallback(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create items
	var items []any
	for i := 0; i < 5; i++ {
		items = append(items, TestItem{
			PK:     "callback-" + string(rune('0'+i)),
			SK:     "item",
			Name:   "Callback" + string(rune('0'+i)),
			Age:    i * 10,
			Status: "active",
		})
	}
	if _, err := tc.client.BatchWriteItems(ctx, tc.table, items); err != nil {
		t.Fatalf("setup BatchWriteItems failed: %v", err)
	}

	var count int
	err := tc.client.ScanCallback(ctx, tc.table, func(items []map[string]types.AttributeValue) error {
		count += len(items)
		return nil
	})
	if err != nil {
		t.Fatalf("ScanCallback failed: %v", err)
	}

	if count < 5 {
		t.Errorf("expected at least 5 items processed, got %d", count)
	}
}

func TestValidationErrors(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	t.Run("empty table name", func(t *testing.T) {
		var got TestItem
		err := tc.client.GetItem(ctx, "", dynamodb.Key{"pk": "test"}, &got)
		if err == nil {
			t.Error("expected error for empty table name")
		}
	})

	t.Run("empty key", func(t *testing.T) {
		var got TestItem
		err := tc.client.GetItem(ctx, tc.table, dynamodb.Key{}, &got)
		if err == nil {
			t.Error("expected error for empty key")
		}
	})
}

func TestTableNotFound(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	var got TestItem
	err := tc.client.GetItem(ctx, "non-existent-table", dynamodb.Key{"pk": "test", "sk": "test"}, &got)

	if !errors.Is(err, dynamodb.ErrTableNotFound) {
		t.Errorf("expected ErrTableNotFound, got: %v", err)
	}
}
