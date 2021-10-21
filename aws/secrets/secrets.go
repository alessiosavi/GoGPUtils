package secrets

import (
	"context"
	"encoding/json"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secretTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"sync"
)

var secretClient *secretsmanager.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		secretClient = secretsmanager.New(secretsmanager.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}
func GetSecret(secretName string) (string, error) {
	value, err := secretClient.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretName)})
	if err != nil {
		return "", err
	}
	return *value.SecretString, err
}

func ListSecret() ([]string, error) {
	secrets, err := secretClient.ListSecrets(context.Background(), &secretsmanager.ListSecretsInput{SortOrder: secretTypes.SortOrderTypeAsc})
	if err != nil {
		return nil, err
	}
	var s []string = make([]string, len(secrets.SecretList))
	for i, secret := range secrets.SecretList {
		s[i] = *secret.Name
	}
	return s, nil
}
func UnmarshalSecret(secretName string, dest interface{}) error {
	secret, err := GetSecret(secretName)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(secret), dest)
}
