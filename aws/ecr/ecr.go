package ecr

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"sync"
)

var ecrClient *ecr.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		ecrClient = ecr.New(ecr.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}
func ListECR() ([]string, error) {

	imgs, err := ecrClient.DescribeRepositories(context.Background(), &ecr.DescribeRepositoriesInput{})
	if err != nil {
		return nil, err
	}

	var images = make([]string, len(imgs.Repositories))

	for i := range imgs.Repositories {
		images[i] = *imgs.Repositories[i].RepositoryName
	}
	for imgs.NextToken != nil {
		imgs, err = ecrClient.DescribeRepositories(context.Background(), &ecr.DescribeRepositoriesInput{
			NextToken: imgs.NextToken,
		})
		if err != nil {
			return nil, err
		}
		for i := range imgs.Repositories {
			images = append(images, *imgs.Repositories[i].RepositoryName)
		}
	}
	return images, nil
}
