package ec2utils

import (
	"context"
	"fmt"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/alessiosavi/GoGPUtils/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
	"sync"
)

var ec2Client *ec2.Client = nil
var once sync.Once

type InstanceDetail struct {
	InstanceName string
	InstanceID   string
}
type Network struct {
	PrivateIPv4 string
	PublicIPv4  string
	PrivateDns  string
	PublicDns   string
	KeyName     string
}

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		ec2Client = ec2.New(ec2.Options{Credentials: cfg.Credentials, Region: cfg.Region, RetryMaxAttempts: 5, RetryMode: aws.RetryModeAdaptive})
	})
}

func ListEC2() ([]ec2types.Instance, error) {
	instances, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}
	var ec2instances []ec2types.Instance
	nextToken := instances.NextToken

	for _, reservation := range instances.Reservations {
		ec2instances = append(ec2instances, reservation.Instances...)
	}
	for nextToken != nil {
		instances, err = ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
		if err != nil {
			return nil, err
		}
		nextToken = instances.NextToken
		for _, reservation := range instances.Reservations {
			ec2instances = append(ec2instances, reservation.Instances...)
		}
	}
	return ec2instances, nil
}

// GetEC2InstanceDetail return a struct that contains {InstanceName, InstanceID} for every EC2 instances
func GetEC2InstanceDetail() ([]InstanceDetail, error) {
	instances, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}
	var ec2instances []InstanceDetail
	nextToken := instances.NextToken
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					ec2instances = append(ec2instances, InstanceDetail{
						InstanceName: *tag.Value,
						InstanceID:   *instance.InstanceId,
					})
				}
			}
		}
	}
	for nextToken != nil {
		instances, err = ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
		if err != nil {
			return nil, err
		}
		nextToken = instances.NextToken
		for _, reservation := range instances.Reservations {
			for _, instance := range reservation.Instances {
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						ec2instances = append(ec2instances, InstanceDetail{
							InstanceName: *tag.Value,
							InstanceID:   *instance.InstanceId,
						})
					}
				}
			}
		}
	}
	return ec2instances, nil
}

func DescribeNetwork(instance string) (*Network, error) {
	describeInstance, err := DescribeInstanceByID(instance)
	if err != nil {
		return nil, err
	}

	data := describeInstance.Reservations[0].Instances[0]

	var network Network
	if data.PrivateIpAddress != nil {
		network.PrivateIPv4 = *data.PrivateIpAddress
	}
	if data.PrivateDnsName != nil {
		network.PrivateDns = *data.PrivateDnsName
	}
	if data.PublicDnsName != nil {
		network.PublicDns = *data.PublicDnsName
	}
	if data.PublicIpAddress != nil {
		network.PublicIPv4 = *data.PublicIpAddress
	}
	return &network, nil
}

func DescribeInstanceByID(instanceID string) (*ec2.DescribeInstancesOutput, error) {
	hosts, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{InstanceIds: []string{instanceID}})
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func DescribeInstanceByName(instanceName string) (*ec2.DescribeInstancesOutput, error) {
	listEC2, err := GetEC2InstanceDetail()
	if err != nil {
		return nil, err
	}

	for _, s := range listEC2 {
		if instanceName == s.InstanceName {
			hosts, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{InstanceIds: []string{s.InstanceID}})
			if err != nil {
				return nil, err
			}
			log.Println(helper.MarshalIndent(hosts))
			return hosts, nil
		}
	}

	return nil, fmt.Errorf("instance %s not found", instanceName)
}
