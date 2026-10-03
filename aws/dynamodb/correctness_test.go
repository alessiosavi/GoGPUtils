package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	ddbsdk "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockAPI struct {
	API
	batchWrite func(*ddbsdk.BatchWriteItemInput) (*ddbsdk.BatchWriteItemOutput, error)
	query      func(*ddbsdk.QueryInput) (*ddbsdk.QueryOutput, error)
	scan       func(*ddbsdk.ScanInput) (*ddbsdk.ScanOutput, error)
}

func (m *mockAPI) BatchWriteItem(_ context.Context, in *ddbsdk.BatchWriteItemInput, _ ...func(*ddbsdk.Options)) (*ddbsdk.BatchWriteItemOutput, error) {
	return m.batchWrite(in)
}
func (m *mockAPI) Query(_ context.Context, in *ddbsdk.QueryInput, _ ...func(*ddbsdk.Options)) (*ddbsdk.QueryOutput, error) {
	return m.query(in)
}
func (m *mockAPI) Scan(_ context.Context, in *ddbsdk.ScanInput, _ ...func(*ddbsdk.Options)) (*ddbsdk.ScanOutput, error) {
	return m.scan(in)
}

type badAttribute struct{ err error }

func (b badAttribute) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return nil, b.err
}

func TestBatchItems_ReturnsEarlierUnprocessedOnLaterError(t *testing.T) {
	for _, operation := range []string{"write", "delete"} {
		for _, failure := range []string{"service", "marshal"} {
			t.Run(operation+"/"+failure, func(t *testing.T) {
				wantErr := errors.New("later chunk failed")
				calls := 0
				client := NewClientWithAPI(&mockAPI{batchWrite: func(in *ddbsdk.BatchWriteItemInput) (*ddbsdk.BatchWriteItemOutput, error) {
					calls++
					if calls == 2 {
						return nil, wantErr
					}
					if len(in.RequestItems["table"]) != 25 {
						t.Fatal("first chunk should contain 25 items")
					}
					return &ddbsdk.BatchWriteItemOutput{UnprocessedItems: map[string][]types.WriteRequest{"table": {in.RequestItems["table"][0]}}}, nil
				}})
				items, keys := make([]any, 26), make([]Key, 26)
				for i := range keys {
					keys[i] = Key{"pk": fmt.Sprint(i)}
					items[i] = keys[i]
				}
				if failure == "marshal" {
					keys[25]["pk"] = badAttribute{err: wantErr}
				}
				var err error
				var count int
				if operation == "write" {
					var got []any
					got, err = client.BatchWriteItems(context.Background(), "table", items)
					count = len(got)
				} else {
					var got []Key
					got, err = client.BatchDeleteItems(context.Background(), "table", keys)
					count = len(got)
				}
				if !errors.Is(err, wantErr) || count != 1 {
					t.Fatalf("unprocessed = %d, err = %v; want 1 and original error", count, err)
				}
				wantCalls := 2
				if failure == "marshal" {
					wantCalls = 1
				}
				if calls != wantCalls {
					t.Fatalf("calls = %d, want %d", calls, wantCalls)
				}
			})
		}
	}
}

func TestBatchItems_ReportsUnprocessedDecodeErrors(t *testing.T) {
	for _, operation := range []string{"write", "delete"} {
		t.Run(operation, func(t *testing.T) {
			client := NewClientWithAPI(&mockAPI{batchWrite: func(in *ddbsdk.BatchWriteItemInput) (*ddbsdk.BatchWriteItemOutput, error) {
				bad := map[string]types.AttributeValue{"pk": &types.AttributeValueMemberN{Value: "invalid-number"}}
				req := types.WriteRequest{PutRequest: &types.PutRequest{Item: bad}}
				if operation == "delete" {
					req = types.WriteRequest{DeleteRequest: &types.DeleteRequest{Key: bad}}
				}
				return &ddbsdk.BatchWriteItemOutput{UnprocessedItems: map[string][]types.WriteRequest{"table": {in.RequestItems["table"][0], req}}}, nil
			}})
			var err error
			var count int
			if operation == "write" {
				var got []any
				got, err = client.BatchWriteItems(context.Background(), "table", []any{Key{"pk": "good"}, Key{"pk": "bad"}})
				count = len(got)
			} else {
				var got []Key
				got, err = client.BatchDeleteItems(context.Background(), "table", []Key{{"pk": "good"}, {"pk": "bad"}})
				count = len(got)
			}
			if err == nil || count != 1 {
				t.Fatalf("unprocessed = %d, err = %v; want 1 and decoding error", count, err)
			}
		})
	}
}

// The interface assertion allows this regression to run against versions without
// the raw methods and fail explicitly, rather than fail to compile.
type rawBatchClient interface {
	BatchWriteItemsRaw(context.Context, string, []map[string]types.AttributeValue) ([]map[string]types.AttributeValue, error)
	BatchDeleteItemsRaw(context.Context, string, []map[string]types.AttributeValue) ([]map[string]types.AttributeValue, error)
}

func TestBatchItemsRaw_PreservesLargeNumbersAndPartialResults(t *testing.T) {
	for _, operation := range []string{"write", "delete"} {
		for _, failLater := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/failLater=%t", operation, failLater), func(t *testing.T) {
				calls := 0
				wantErr := errors.New("second chunk failed")
				items := make([]map[string]types.AttributeValue, 26)
				for i := range items {
					items[i] = map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberN{Value: "9007199254740993"},
						"sk": &types.AttributeValueMemberB{Value: []byte{0, 255, byte(i)}},
					}
					if operation == "write" {
						items[i]["nested"] = &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{"n": &types.AttributeValueMemberN{Value: "1.234567890123456789"}}}
					}
				}
				client := NewClientWithAPI(&mockAPI{batchWrite: func(in *ddbsdk.BatchWriteItemInput) (*ddbsdk.BatchWriteItemOutput, error) {
					calls++
					wantSize := 25
					if calls == 2 {
						wantSize = 1
					}
					if len(in.RequestItems["table"]) != wantSize {
						t.Fatal("wrong chunk size")
					}
					for i, req := range in.RequestItems["table"] {
						var got map[string]types.AttributeValue
						if operation == "write" {
							if req.PutRequest == nil || req.DeleteRequest != nil {
								t.Fatal("wrong request kind")
							}
							got = req.PutRequest.Item
						} else {
							if req.DeleteRequest == nil || req.PutRequest != nil {
								t.Fatal("wrong request kind")
							}
							got = req.DeleteRequest.Key
						}
						if !reflect.DeepEqual(got, items[(calls-1)*25+i]) {
							t.Fatal("request changed attribute values")
						}
					}
					if calls == 2 && failLater {
						return nil, wantErr
					}
					return &ddbsdk.BatchWriteItemOutput{UnprocessedItems: in.RequestItems}, nil
				}})
				raw, ok := any(client).(rawBatchClient)
				if !ok {
					t.Fatal("lossless raw batch methods are missing")
				}
				batch := raw.BatchWriteItemsRaw
				if operation == "delete" {
					batch = raw.BatchDeleteItemsRaw
				}
				got, err := batch(context.Background(), "table", items)
				want := items
				if failLater {
					want = items[:25]
					if !errors.Is(err, wantErr) {
						t.Fatalf("err = %v", err)
					}
				} else if err != nil {
					t.Fatal(err)
				}
				if calls != 2 || !reflect.DeepEqual(got, want) {
					t.Fatalf("calls = %d, unprocessed = %#v", calls, got)
				}
			})
		}
	}
}

func TestQueryAndScan_RejectEmptyTableBeforeOptions(t *testing.T) {
	for _, operation := range []string{"Query", "QueryAll", "Scan", "ScanAll", "ScanCallback"} {
		t.Run(operation, func(t *testing.T) {
			client := NewClientWithAPI(&mockAPI{
				query: func(*ddbsdk.QueryInput) (*ddbsdk.QueryOutput, error) {
					t.Fatal("Query API should not be called for an empty table")
					return nil, nil
				},
				scan: func(*ddbsdk.ScanInput) (*ddbsdk.ScanOutput, error) {
					t.Fatal("Scan API should not be called for an empty table")
					return nil, nil
				},
			})
			queryOpt := func(*queryOptions) { t.Fatal("query options applied before table validation") }
			scanOpt := func(*scanOptions) { t.Fatal("scan options applied before table validation") }
			ctx := context.Background()
			keyExpr := expression.Key("pk").Equal(expression.Value("value"))
			var dest []map[string]string
			var err error
			switch operation {
			case "Query":
				_, err = client.Query(ctx, "", keyExpr, queryOpt)
			case "QueryAll":
				err = client.QueryAll(ctx, "", keyExpr, &dest, queryOpt)
			case "Scan":
				_, err = client.Scan(ctx, "", scanOpt)
			case "ScanAll":
				err = client.ScanAll(ctx, "", &dest, scanOpt)
			case "ScanCallback":
				err = client.ScanCallback(ctx, "", func([]map[string]types.AttributeValue) error {
					t.Fatal("scan callback should not be called for an empty table")
					return nil
				}, scanOpt)
			}
			if !errors.Is(err, aws.ErrEmptyTable) {
				t.Fatalf("err = %v, want ErrEmptyTable", err)
			}
		})
	}
}

func TestPagination_AppliesOptionsOnceAndAdvancesCursor(t *testing.T) {
	for _, operation := range []string{"query", "scan", "callback"} {
		t.Run(operation, func(t *testing.T) {
			calls, applied, count := 0, 0, 0
			cursor := func(n int) map[string]types.AttributeValue {
				return map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: fmt.Sprint(n)}}
			}
			page := func(start map[string]types.AttributeValue) ([]map[string]types.AttributeValue, map[string]types.AttributeValue) {
				if !reflect.DeepEqual(start, cursor(calls)) {
					t.Errorf("cursor = %v, want %v", start, cursor(calls))
				}
				calls++
				var next map[string]types.AttributeValue
				if calls < 3 {
					next = cursor(calls)
				}
				return []map[string]types.AttributeValue{cursor(calls)}, next
			}
			client := NewClientWithAPI(&mockAPI{
				query: func(in *ddbsdk.QueryInput) (*ddbsdk.QueryOutput, error) {
					if awssdk.ToInt32(in.Limit) != 2 || awssdk.ToString(in.IndexName) != "index" {
						t.Error("lost query options")
					}
					items, next := page(in.ExclusiveStartKey)
					return &ddbsdk.QueryOutput{Items: items, LastEvaluatedKey: next}, nil
				},
				scan: func(in *ddbsdk.ScanInput) (*ddbsdk.ScanOutput, error) {
					if awssdk.ToInt32(in.Limit) != 2 || awssdk.ToInt32(in.TotalSegments) != 4 {
						t.Error("lost scan options")
					}
					items, next := page(in.ExclusiveStartKey)
					return &ddbsdk.ScanOutput{Items: items, LastEvaluatedKey: next}, nil
				},
			})
			var dest []map[string]string
			var err error
			if operation == "query" {
				err = client.QueryAll(context.Background(), "table", expression.Key("pk").Equal(expression.Value("value")), &dest, WithStartKey(cursor(0)), WithLimit(2), WithIndex("index"), func(*queryOptions) { applied++ })
				count = len(dest)
			} else {
				opts := []ScanOption{WithScanStartKey(cursor(0)), WithScanLimit(2), WithParallelScan(0, 4), func(*scanOptions) { applied++ }}
				if operation == "scan" {
					err = client.ScanAll(context.Background(), "table", &dest, opts...)
					count = len(dest)
				} else {
					err = client.ScanCallback(context.Background(), "table", func(items []map[string]types.AttributeValue) error { count += len(items); return nil }, opts...)
				}
			}
			if err != nil || calls != 3 || count != 3 || applied != 1 {
				t.Fatalf("err = %v, calls = %d, items = %d, option applications = %d", err, calls, count, applied)
			}
		})
	}
}
