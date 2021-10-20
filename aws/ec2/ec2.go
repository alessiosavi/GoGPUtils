package ec2

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
		ec2Client = ec2.New(ec2.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}

func ListEC2() ([]InstanceDetail, error) {
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
	describeInstance, err := DescribeInstance(instance)
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

func DescribeInstance(instance string) (*ec2.DescribeInstancesOutput, error) {
	hosts, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{InstanceIds: []string{instance}})
	if err != nil {
		return nil, err
	}
	return hosts, nil
}
