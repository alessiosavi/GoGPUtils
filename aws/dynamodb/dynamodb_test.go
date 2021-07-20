package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
	"time"
)

type Report struct {
	Finished   bool      `json:"finished"`
	Filename   string    `json:"filename"`
	StartTime  time.Time `json:"start_time"`
	FinishTime time.Time `json:"finish_time"`
}

func TestCreateTable(t *testing.T) {
	type args struct {
		definition *dynamodb.CreateTableInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_ok",
			args: args{
				definition: &dynamodb.CreateTableInput{
					AttributeDefinitions: []types.AttributeDefinition{
						{
							AttributeName: aws.String("Filename"),
							AttributeType: types.ScalarAttributeTypeS,
						},
					},
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("Filename"),
							KeyType:       types.KeyTypeHash,
						},
					},
					TableName:              aws.String("table_test"),
					BillingMode:            types.BillingModeProvisioned,
					GlobalSecondaryIndexes: nil,
					LocalSecondaryIndexes:  nil,
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
					SSESpecification:    nil,
					StreamSpecification: nil,
					Tags:                nil,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateTable(tt.args.definition); (err != nil) != tt.wantErr {
				t.Errorf("CreateTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteItem(t *testing.T) {
	type args struct {
		tableName string
		item      interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_write",
			args: args{
				tableName: "table_test",
				item: Report{
					Finished:  false,
					Filename:  "test1",
					StartTime: time.Now(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteItem(tt.args.tableName, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("WriteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteTable(t *testing.T) {
	type args struct {
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "remove_test",
			args:    args{tableName: "table_test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteTable(tt.args.tableName); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}