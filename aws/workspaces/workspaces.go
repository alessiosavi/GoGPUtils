package workspaces

import (
	"context"
	"fmt"
	"sync"

	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
)

var workspaceClient *workspaces.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		workspaceClient = workspaces.New(workspaces.Options{Credentials: cfg.Credentials, Region: cfg.Region, RetryMaxAttempts: 5, RetryMode: aws.RetryModeAdaptive})
	})
}

func ListWorkspaces() (map[string][]string, error) {
	wsList, err := workspaceClient.DescribeWorkspaces(context.Background(), &workspaces.DescribeWorkspacesInput{})
	if err != nil {
		return nil, err
	}

	var workspacesList = make(map[string][]string)

	for _, workspace := range wsList.Workspaces {
		workspacesList[*workspace.UserName] = append(workspacesList[*workspace.UserName], *workspace.WorkspaceId)
	}

	continuationToken := wsList.NextToken

	for continuationToken != nil {
		wsList, err = workspaceClient.DescribeWorkspaces(context.Background(), &workspaces.DescribeWorkspacesInput{NextToken: continuationToken})
		if err != nil {
			return nil, err
		}
		for _, workspace := range wsList.Workspaces {
			workspacesList[*workspace.UserName] = append(workspacesList[*workspace.UserName], *workspace.WorkspaceId)
		}
		continuationToken = wsList.NextToken
	}
	return workspacesList, nil

}

func GetWorkspaces(username string) (*workspaces.DescribeWorkspacesOutput, error) {
	listWorkspaces, err := ListWorkspaces()
	if err != nil {
		return nil, nil
	}

	var v []string
	var ok bool
	if v, ok = listWorkspaces[username]; !ok {
		return nil, fmt.Errorf("username %s not found", username)
	}
	return workspaceClient.DescribeWorkspaces(context.Background(), &workspaces.DescribeWorkspacesInput{
		WorkspaceIds: v,
	})

}

//func DescribeImage(username string) error {
//	workspaceClient.DescribeClientProperties()
//}

func StartWorkspaces() {

}

func RemoveWorkspaces() {

}
