package sesutils

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	mailTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"sync"
)

var sesClient *sesv2.Client = nil
var once sync.Once

type MailConf struct {
	FromName string   `json:"from_name,omitempty"`
	FromMail string   `json:"from_mail,omitempty"`
	To       string   `json:"to,omitempty"`
	CC       []string `json:"cc,omitempty"`
}

func GetClient() *sesv2.Client {
	return sesClient
}

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		sesClient = sesv2.New(sesv2.Options{Credentials: cfg.Credentials, Region: cfg.Region, RetryMaxAttempts: 5, RetryMode: aws.RetryModeAdaptive})
	})
}
func SendMail(data []byte) error {
	_, err := sesClient.SendEmail(context.Background(), &sesv2.SendEmailInput{
		Content: &mailTypes.EmailContent{
			Raw: &mailTypes.RawMessage{
				Data: data,
			},
		},
	})
	return err
}
