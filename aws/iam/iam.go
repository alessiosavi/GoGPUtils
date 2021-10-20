package iam

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"sync"
)

type IamUser struct {
	Username        string
	Password        string
	AccessKeyId     string
	SecretAccessKey string
}

var iamClient *iam.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		iamClient = iam.New(iam.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}

func CreateIamUser(username, password string, passwordResetRequired, accessKeys bool) (*IamUser, error) {
	// Create User
	if _, err := iamClient.CreateUser(context.Background(), &iam.CreateUserInput{
		UserName: aws.String(username),
	}); err != nil {
		return nil, err
	}
	// Setting password policy
	if _, err := iamClient.CreateLoginProfile(context.Background(), &iam.CreateLoginProfileInput{
		UserName:              aws.String(username),
		Password:              aws.String(password),
		PasswordResetRequired: passwordResetRequired,
	}); err != nil {
		return nil, err
	}

	// Setting access key if necessary
	accessKeyId := ""
	secretAccessKey := ""
	if accessKeys {
		key, err := iamClient.CreateAccessKey(context.Background(), &iam.CreateAccessKeyInput{
			UserName: aws.String(username),
		})
		if err != nil {
			return nil, err
		}
		accessKeyId = *key.AccessKey.AccessKeyId
		secretAccessKey = *key.AccessKey.SecretAccessKey
	}

	return &IamUser{Username: username, Password: password, AccessKeyId: accessKeyId, SecretAccessKey: secretAccessKey}, nil
}
