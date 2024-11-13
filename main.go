package main

import (
	"context"
	"log"
	"time"

	"eni-project/internal/ec2"
	"github.com/aws/aws-sdk-go-v2/config"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
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

	subnetID := "subnet-0a7bd03887dc3cbd5"
	instanceID := "i-04890aa7cd8cf81f3"

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create ENI
	eni, err := eniManager.CreateENI(ctx, subnetID)
	if err != nil {
		log.Fatalf("Failed to create ENI: %v", err)
	}
	log.Printf("Created ENI: %s\n", *eni.NetworkInterface.NetworkInterfaceId)

	// Attach ENI
	attachID, err := eniManager.AttachENI(ctx, *eni.NetworkInterface.NetworkInterfaceId, instanceID)
	if err != nil {
		log.Fatalf("Failed to attach ENI: %v", err)
	}
	log.Printf("Attached ENI with attachment ID: %s\n", *attachID)

	// Wait a bit before detaching (in a real application, you might wait until you're done using it)
	time.Sleep(5 * time.Second)

	// Detach ENI
	if err := eniManager.DetachENI(ctx, *attachID); err != nil {
		log.Fatalf("Failed to detach ENI: %v", err)
	}
	log.Println("Detached ENI")

	// Wait for detachment to complete
	time.Sleep(5 * time.Second)

	// Delete ENI
	if err := eniManager.DeleteENI(ctx, *eni.NetworkInterface.NetworkInterfaceId); err != nil {
		log.Fatalf("Failed to delete ENI: %v", err)
	}
	log.Println("Deleted ENI")
}
