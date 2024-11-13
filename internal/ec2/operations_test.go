// internal/ec2/operations_test.go
package ec2

import (
	"context"
	"testing"

	"eni-project/internal/ec2/mocks"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestENIManager_CreateENI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	config := ENIConfig{
		SubnetID:         "subnet-12345678",
		Description:      "Test ENI",
		SecurityGroupIDs: []string{"sg-12345678"},
		PrivateIPCount:   2,
		IPv6AddressCount: 1,
		Tags: map[string]string{
			"Name": "TestENI",
			"Env":  "Test",
		},
	}

	expectedInput := &ec2.CreateNetworkInterfaceInput{
		SubnetId:                       aws.String(config.SubnetID),
		Description:                    aws.String(config.Description),
		Groups:                         config.SecurityGroupIDs,
		SecondaryPrivateIpAddressCount: aws.Int32(config.PrivateIPCount),
		Ipv6AddressCount:               aws.Int32(config.IPv6AddressCount),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeNetworkInterface,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("TestENI"),
					},
					{
						Key:   aws.String("Env"),
						Value: aws.String("Test"),
					},
				},
			},
		},
	}

	expectedOutput := &ec2.CreateNetworkInterfaceOutput{
		NetworkInterface: &types.NetworkInterface{
			NetworkInterfaceId: aws.String("eni-12345678"),
			SubnetId:           aws.String(config.SubnetID),
			Description:        aws.String(config.Description),
		},
	}

	mockClient.EXPECT().
		CreateNetworkInterface(gomock.Any(), gomock.Eq(expectedInput)).
		Return(expectedOutput, nil)

	result, err := manager.CreateENI(context.Background(), config)
	assert.NoError(t, err)
	assert.Equal(t, "eni-12345678", *result.NetworkInterface.NetworkInterfaceId)
}

func TestENIManager_AttachENI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	expectedInput := &ec2.AttachNetworkInterfaceInput{
		NetworkInterfaceId: aws.String("eni-12345678"),
		InstanceId:         aws.String("i-12345678"),
		DeviceIndex:        aws.Int32(1),
	}

	expectedOutput := &ec2.AttachNetworkInterfaceOutput{
		AttachmentId: aws.String("eni-attach-12345678"),
	}

	mockClient.EXPECT().
		AttachNetworkInterface(gomock.Any(), gomock.Eq(expectedInput)).
		Return(expectedOutput, nil)

	attachmentID, err := manager.AttachENI(context.Background(), "eni-12345678", "i-12345678", 1)
	assert.NoError(t, err)
	assert.Equal(t, "eni-attach-12345678", *attachmentID)
}

func TestENIManager_DetachENI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	expectedInput := &ec2.DetachNetworkInterfaceInput{
		AttachmentId: aws.String("eni-attach-12345678"),
		Force:        aws.Bool(true),
	}

	mockClient.EXPECT().
		DetachNetworkInterface(gomock.Any(), gomock.Eq(expectedInput)).
		Return(&ec2.DetachNetworkInterfaceOutput{}, nil)

	err := manager.DetachENI(context.Background(), "eni-attach-12345678", true)
	assert.NoError(t, err)
}

func TestENIManager_AssignPrivateIPs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	expectedInput := &ec2.AssignPrivateIpAddressesInput{
		NetworkInterfaceId:             aws.String("eni-12345678"),
		SecondaryPrivateIpAddressCount: aws.Int32(2),
	}

	mockClient.EXPECT().
		AssignPrivateIpAddresses(gomock.Any(), gomock.Eq(expectedInput)).
		Return(&ec2.AssignPrivateIpAddressesOutput{}, nil)

	err := manager.AssignPrivateIPs(context.Background(), "eni-12345678", 2, nil)
	assert.NoError(t, err)
}

func TestENIManager_AssignSpecificPrivateIPs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	specificIPs := []string{"10.0.0.10", "10.0.0.11"}
	expectedInput := &ec2.AssignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String("eni-12345678"),
		PrivateIpAddresses: specificIPs,
	}

	mockClient.EXPECT().
		AssignPrivateIpAddresses(gomock.Any(), gomock.Eq(expectedInput)).
		Return(&ec2.AssignPrivateIpAddressesOutput{}, nil)

	err := manager.AssignPrivateIPs(context.Background(), "eni-12345678", 0, specificIPs)
	assert.NoError(t, err)
}

func TestENIManager_DescribeENIs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	filters := []types.Filter{
		{
			Name:   aws.String("subnet-id"),
			Values: []string{"subnet-12345678"},
		},
	}

	expectedInput := &ec2.DescribeNetworkInterfacesInput{
		Filters: filters,
	}

	expectedOutput := &ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{
			{
				NetworkInterfaceId: aws.String("eni-12345678"),
				SubnetId:           aws.String("subnet-12345678"),
			},
		},
	}

	mockClient.EXPECT().
		DescribeNetworkInterfaces(gomock.Any(), gomock.Eq(expectedInput)).
		Return(expectedOutput, nil)

	result, err := manager.DescribeENIs(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, result.NetworkInterfaces, 1)
	assert.Equal(t, "eni-12345678", *result.NetworkInterfaces[0].NetworkInterfaceId)
}
