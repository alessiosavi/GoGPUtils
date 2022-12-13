package identity_store

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"sort"
	"sync"
)

var identityClient *identitystore.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		identityClient = identitystore.New(identitystore.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}

func ListUsers(identityStore string) ([]types.User, error) {
	res, err := identityClient.ListUsers(context.Background(), &identitystore.ListUsersInput{IdentityStoreId: aws.String(identityStore)})
	if err != nil {
		return nil, err
	}
	var users []types.User
	users = append(users, res.Users...)
	continuationToken := res.NextToken

	for continuationToken != nil {
		res, err = identityClient.ListUsers(context.Background(), &identitystore.ListUsersInput{NextToken: continuationToken})
		if err != nil {
			return nil, err
		}

		continuationToken = res.NextToken
		users = append(users, res.Users...)
	}

	sort.Slice(users, func(i, j int) bool {
		return *users[i].UserName < *users[j].UserName
	})
	return users, nil
}

func DescribeUser(userId, identityStore string) (*identitystore.DescribeUserOutput, error) {
	return identityClient.DescribeUser(context.Background(), &identitystore.DescribeUserInput{
		IdentityStoreId: aws.String(identityStore),
		UserId:          aws.String(userId),
	})
}
