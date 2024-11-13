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

	subnetID := "subnet-12345678"
	expectedENIID := "eni-12345678"

	expectedInput := &ec2.CreateNetworkInterfaceInput{
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

	expectedOutput := &ec2.CreateNetworkInterfaceOutput{
		NetworkInterface: &types.NetworkInterface{
			NetworkInterfaceId: &expectedENIID,
		},
	}

	mockClient.EXPECT().
		CreateNetworkInterface(gomock.Any(), gomock.Eq(expectedInput)).
		Return(expectedOutput, nil)

	result, err := manager.CreateENI(context.Background(), subnetID)

	assert.NoError(t, err)
	assert.Equal(t, expectedENIID, *result.NetworkInterface.NetworkInterfaceId)
}

func TestENIManager_AttachENI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	networkInterfaceID := "eni-12345678"
	instanceID := "i-1234567890abcd"
	expectedAttachmentID := "eni-attach-12345678"

	expectedInput := &ec2.AttachNetworkInterfaceInput{
		DeviceIndex:        aws.Int32(1),
		InstanceId:         &instanceID,
		NetworkInterfaceId: &networkInterfaceID,
	}

	mockClient.EXPECT().
		AttachNetworkInterface(gomock.Any(), gomock.Eq(expectedInput)).
		Return(&ec2.AttachNetworkInterfaceOutput{AttachmentId: &expectedAttachmentID}, nil)

	attachmentID, err := manager.AttachENI(context.Background(), networkInterfaceID, instanceID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAttachmentID, *attachmentID)
}

func TestENIManager_DetachENI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	attachmentID := "eni-attach-12345678"

	expectedInput := &ec2.DetachNetworkInterfaceInput{
		AttachmentId: &attachmentID,
		Force:        aws.Bool(true),
	}

	mockClient.EXPECT().
		DetachNetworkInterface(gomock.Any(), gomock.Eq(expectedInput)).
		Return(&ec2.DetachNetworkInterfaceOutput{}, nil)

	err := manager.DetachENI(context.Background(), attachmentID)
	assert.NoError(t, err)
}

func TestENIManager_DeleteENI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockEC2ClientAPI(ctrl)
	manager := NewENIManager(mockClient)

	networkInterfaceID := "eni-12345678"

	expectedInput := &ec2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: &networkInterfaceID,
	}

	mockClient.EXPECT().
		DeleteNetworkInterface(gomock.Any(), gomock.Eq(expectedInput)).
		Return(&ec2.DeleteNetworkInterfaceOutput{}, nil)

	err := manager.DeleteENI(context.Background(), networkInterfaceID)
	assert.NoError(t, err)
}
