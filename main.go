package main

import (
	"context"
	"log"
	"time"

	"eni-project/internal/ec2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// Create EC2 client
	client := awsec2.NewFromConfig(cfg)

	// Create ENI manager
	eniManager := ec2.NewENIManager(client)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Example parameters
	subnetID := "subnet-0a7bd03887dc3cbd5"
	instanceID := "i-04890aa7cd8cf81f3"
	securityGroupID := "sg-0f9acdf364ab834f2"

	// Example: Create an ENI
	eniConfig := ec2.ENIConfig{
		SubnetID:         subnetID,
		Description:      "Example ENI",
		SecurityGroupIDs: []string{securityGroupID},
		PrivateIPCount:   2,
		IPv6AddressCount: 0,
		Tags: map[string]string{
			"Name":        "example-eni",
			"Environment": "development",
			"ManagedBy":   "eni-manager",
		},
	}

	log.Println("Creating ENI...")
	eni, err := eniManager.CreateENI(ctx, eniConfig)
	if err != nil {
		log.Fatalf("Failed to create ENI: %v", err)
	}
	log.Printf("Created ENI: %s\n", *eni.NetworkInterface.NetworkInterfaceId)

	// Example: Attach ENI to an instance
	log.Printf("Attaching ENI to instance %s...\n", instanceID)
	attachID, err := eniManager.AttachENI(ctx, *eni.NetworkInterface.NetworkInterfaceId, instanceID, 1)
	if err != nil {
		log.Fatalf("Failed to attach ENI: %v", err)
	}
	log.Printf("Attached ENI with attachment ID: %s\n", *attachID)

	// Example: Assign additional private IPs
	log.Println("Assigning additional private IPs...")
	err = eniManager.AssignPrivateIPs(ctx, *eni.NetworkInterface.NetworkInterfaceId, 2, nil)
	if err != nil {
		log.Printf("Failed to assign private IPs: %v", err)
	}

	// Example: Describe ENIs in the subnet
	log.Println("Describing ENIs...")
	filters := []types.Filter{
		{
			Name:   aws.String("subnet-id"),
			Values: []string{eniConfig.SubnetID},
		},
	}

	enis, err := eniManager.DescribeENIs(ctx, filters)
	if err != nil {
		log.Printf("Failed to describe ENIs: %v", err)
	} else {
		for _, eni := range enis.NetworkInterfaces {
			log.Printf("Found ENI: %s, Status: %s\n", *eni.NetworkInterfaceId, eni.Status)
		}
	}

	// Example: Modify ENI attributes
	log.Println("Modifying ENI attributes...")
	modifyConfig := ec2.ENIModifyConfig{
		Description: aws.String("Updated description"),
	}
	err = eniManager.ModifyENIAttribute(ctx, *eni.NetworkInterface.NetworkInterfaceId, modifyConfig)
	if err != nil {
		log.Printf("Failed to modify ENI: %v", err)
	}

	// Wait before cleanup
	log.Println("Waiting for 5 seconds before cleanup...")
	time.Sleep(5 * time.Second)

	// Example: Detach ENI
	log.Println("Detaching ENI...")
	err = eniManager.DetachENI(ctx, *attachID, true)
	if err != nil {
		log.Fatalf("Failed to detach ENI: %v", err)
	}

	// Wait for detachment to complete
	log.Println("Waiting for 5 seconds after detachment...")
	time.Sleep(5 * time.Second)

	// Example: Delete ENI
	log.Println("Deleting ENI...")
	err = eniManager.DeleteENI(ctx, *eni.NetworkInterface.NetworkInterfaceId)
	if err != nil {
		log.Fatalf("Failed to delete ENI: %v", err)
	}

	log.Println("ENI management operations completed successfully")
}
