package ec2

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ENIManager handles ENI operations
type ENIManager struct {
	client EC2ClientAPI
}

// NewENIManager creates a new ENI manager
func NewENIManager(client EC2ClientAPI) *ENIManager {
	return &ENIManager{client: client}
}

// CreateENI creates a new ENI in the specified subnet
func (m *ENIManager) CreateENI(ctx context.Context, subnetID string) (*ec2.CreateNetworkInterfaceOutput, error) {
	input := &ec2.CreateNetworkInterfaceInput{
		SubnetId: &subnetID,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeNetworkInterface,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("Example ENI"),
					},
					{
						Key:   aws.String("CreatedBy"),
						Value: aws.String("ENIManager"),
					},
				},
			},
		},
	}

	return m.client.CreateNetworkInterface(ctx, input)
}

// AttachENI attaches an ENI to an EC2 instance
func (m *ENIManager) AttachENI(ctx context.Context, networkInterfaceID, instanceID string) (*string, error) {
	input := &ec2.AttachNetworkInterfaceInput{
		DeviceIndex:        aws.Int32(1),
		InstanceId:         &instanceID,
		NetworkInterfaceId: &networkInterfaceID,
	}

	result, err := m.client.AttachNetworkInterface(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to attach ENI: %w", err)
	}

	return result.AttachmentId, nil
}

// DetachENI detaches an ENI from an EC2 instance
func (m *ENIManager) DetachENI(ctx context.Context, attachmentID string) error {
	input := &ec2.DetachNetworkInterfaceInput{
		AttachmentId: &attachmentID,
		Force:        aws.Bool(true),
	}

	_, err := m.client.DetachNetworkInterface(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to detach ENI: %w", err)
	}

	return nil
}

// DeleteENI deletes an ENI
func (m *ENIManager) DeleteENI(ctx context.Context, networkInterfaceID string) error {
	input := &ec2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: &networkInterfaceID,
	}

	_, err := m.client.DeleteNetworkInterface(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete ENI: %w", err)
	}

	return nil
}
