package secrets

import (
	"context"
	"encoding/json"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecret(secretName string) (string, error) {
	cfg, err := awsutils.New()
	if err != nil {
		return "", err
	}
	client := secretsmanager.New(secretsmanager.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	value, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretName)})
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
