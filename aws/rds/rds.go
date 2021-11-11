package rds

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/alessiosavi/GoGPUtils/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"log"
	"sync"
)

var rdsClient *rds.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		rdsClient = rds.New(rds.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}

func ListRDS() ([]string, error) {
	clustersList, err := rdsClient.DescribeDBInstances(context.Background(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, err
	}

	var clusters = make([]string, len(clustersList.DBInstances))
	for i, clusterName := range clustersList.DBInstances {
		clusters[i] = *clusterName.DBInstanceIdentifier
	}

	continuationToken := clustersList.Marker
	for continuationToken != nil {
		clustersList, err = rdsClient.DescribeDBInstances(context.Background(), &rds.DescribeDBInstancesInput{Marker: continuationToken})
		if err != nil {
			return nil, err
		}
		continuationToken = clustersList.Marker
		for _, clusterName := range clustersList.DBInstances {
			clusters = append(clusters, *clusterName.DBInstanceIdentifier)
		}
	}
	log.Println(clusters)
	return clusters, nil
}

func DescribeInstanceByID(instanceID string) *rds.DescribeDBInstancesOutput {
	instances, err := rdsClient.DescribeDBInstances(context.Background(), &rds.DescribeDBInstancesInput{DBInstanceIdentifier: aws.String(instanceID)})
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(instances))
	return instances
}
