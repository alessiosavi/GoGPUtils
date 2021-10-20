package secrets

import (
	"context"
	"encoding/json"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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
	cfg, err := awsutils.New()
	if err != nil {
		return "", err
	}
	secretClient := secretsmanager.New(secretsmanager.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	value, err := secretClient.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretName)})
	if err != nil {
		return "", err
	}
	return *value.SecretString, err
}

func UnmarshalSecret(secretName string, dest interface{}) error {
	secret, err := GetSecret(secretName)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(secret), dest)
}
