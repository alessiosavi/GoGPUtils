package iamutils

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/alessiosavi/GoGPUtils/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
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
func ListRoles(prefix *string) ([]iamTypes.Role, error) {
	var res []iamTypes.Role
	roles, err := iamClient.ListRoles(context.Background(), &iam.ListRolesInput{PathPrefix: prefix})
	if err != nil {
		return nil, err
	}

	var marker *string = roles.Marker
	for roles.IsTruncated {
		roles, err = iamClient.ListRoles(context.Background(), &iam.ListRolesInput{
			Marker:     marker,
			PathPrefix: prefix,
		})
		if err != nil {
			return nil, err
		}

		res = append(res, roles.Roles...)
		marker = roles.Marker
	}

	return res, nil
}

func Info(r iamTypes.Role) string {
	return helper.MarshalIndent(r)
}

func GetRole(name string) (*iamTypes.Role, error) {
	role, err := iamClient.GetRole(context.Background(), &iam.GetRoleInput{RoleName: aws.String(name)})
	if err != nil {
		return nil, err
	}
	return role.Role, nil
}
