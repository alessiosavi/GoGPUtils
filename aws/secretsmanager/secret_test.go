package secretsmanager

import (
	"context"
	"errors"
	"fmt"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	smsdk "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/smithy-go"
)

type secretAPI struct {
	API
	err error
}

func (m *secretAPI) GetSecretValue(context.Context, *smsdk.GetSecretValueInput, ...func(*smsdk.Options)) (*smsdk.GetSecretValueOutput, error) {
	return nil, fmt.Errorf("request: %w", m.err)
}

func TestGetSecret_OnlyClassifiesDeletionRequestsAsDeleted(t *testing.T) {
	for _, tt := range []struct {
		name    string
		err     error
		deleted bool
	}{
		{"deleted", &types.InvalidRequestException{Message: awssdk.String("You can't perform this operation on the secret because it was marked for deletion.")}, true},
		{"other invalid request", &types.InvalidRequestException{Message: awssdk.String("You tried to enable rotation on a secret that doesn't already have a Lambda function ARN configured.")}, false},
		{"empty message", &types.InvalidRequestException{}, false},
		{"generic deleted", &smithy.GenericAPIError{Code: "InvalidRequestException", Message: "Secret is scheduled for deletion"}, true},
		{"untyped text", errors.New("InvalidRequestException: marked for deletion"), false},
	} {
		for _, binary := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/binary=%t", tt.name, binary), func(t *testing.T) {
				client := NewClientWithAPI(&secretAPI{err: tt.err})
				var err error
				if binary {
					_, err = client.GetSecretBinary(context.Background(), "secret")
				} else {
					_, err = client.GetSecretString(context.Background(), "secret")
				}
				if errors.Is(err, ErrSecretDeleted) != tt.deleted {
					t.Errorf("err = %v, deleted = %t", err, tt.deleted)
				}
				if !tt.deleted && !errors.Is(err, tt.err) {
					t.Error("original error lost")
				}
			})
		}
	}
}
