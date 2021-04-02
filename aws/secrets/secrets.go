package secrets

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecrets(secretARN string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}
	client := secretsmanager.New(secretsmanager.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	value, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretARN)})
	if err != nil {
		panic(err)
	}
	return *value.SecretString, err
}
