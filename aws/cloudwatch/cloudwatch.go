package cloudwatchutils

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"sync"
	"time"
)

var cloudwatchClient *cloudwatchlogs.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		cloudwatchClient = cloudwatchlogs.New(cloudwatchlogs.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}

func GetLogGroups() (map[string]types.LogGroup, error) {
	var res = make(map[string]types.LogGroup)

	groups, err := cloudwatchClient.DescribeLogGroups(context.Background(), &cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		return nil, err
	}

	for _, group := range groups.LogGroups {
		res[*group.LogGroupName] = group
	}

	for groups.NextToken != nil {
		groups, err = cloudwatchClient.DescribeLogGroups(context.Background(), &cloudwatchlogs.DescribeLogGroupsInput{
			NextToken: groups.NextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, group := range groups.LogGroups {
			res[*group.LogGroupName] = group
		}
	}

	return res, nil
}

func ExportLog(bucket, logGroupName, destinationPrefix string, start, stop time.Time) (*cloudwatchlogs.CreateExportTaskOutput, error) {
	return cloudwatchClient.CreateExportTask(context.Background(), &cloudwatchlogs.CreateExportTaskInput{
		Destination:       aws.String(bucket),
		From:              aws.Int64(start.UTC().UnixMilli()),
		LogGroupName:      aws.String(logGroupName),
		To:                aws.Int64(stop.UTC().UnixMilli()),
		DestinationPrefix: aws.String(destinationPrefix),
		TaskName:          aws.String("github.com/alessiosavi/GoGPUtils"),
	})
}

func DescribeExportTask(taskId string) ([]types.ExportTask, error) {
	tasks, err := cloudwatchClient.DescribeExportTasks(context.Background(), &cloudwatchlogs.DescribeExportTasksInput{
		StatusCode: "",
		TaskId:     aws.String(taskId),
	})
	var res []types.ExportTask = make([]types.ExportTask, len(tasks.ExportTasks))
	if err != nil {
		return res, err
	}
	for i := range tasks.ExportTasks {
		res[i] = tasks.ExportTasks[i]
	}
	continuationToken := tasks.NextToken
	for continuationToken != nil {
		tasks, err = cloudwatchClient.DescribeExportTasks(context.Background(), &cloudwatchlogs.DescribeExportTasksInput{
			NextToken:  continuationToken,
			StatusCode: "",
			TaskId:     aws.String(taskId),
		})
		res = append(res, tasks.ExportTasks...)
	}
	return res, nil
}
