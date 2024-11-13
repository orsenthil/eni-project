package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ENIManager struct {
	client EC2ClientAPI
}

func NewENIManager(client EC2ClientAPI) *ENIManager {
	return &ENIManager{client: client}
}

func (m *ENIManager) CreateENI(ctx context.Context, config ENIConfig) (*ec2.CreateNetworkInterfaceOutput, error) {
	var tags []types.TagSpecification
	if len(config.Tags) > 0 {
		var tagList []types.Tag
		for k, v := range config.Tags {
			tagList = append(tagList, types.Tag{
				Key:   aws.String(k),
				Value: aws.String(v),
			})
		}
		tags = append(tags, types.TagSpecification{
			ResourceType: types.ResourceTypeNetworkInterface,
			Tags:         tagList,
		})
	}

	input := &ec2.CreateNetworkInterfaceInput{
		SubnetId:          aws.String(config.SubnetID),
		Description:       aws.String(config.Description),
		Groups:            config.SecurityGroupIDs,
		TagSpecifications: tags,
	}

	if config.PrivateIPCount > 0 {
		input.SecondaryPrivateIpAddressCount = aws.Int32(config.PrivateIPCount)
	}

	if config.IPv6AddressCount > 0 {
		input.Ipv6AddressCount = aws.Int32(config.IPv6AddressCount)
	}

	return m.client.CreateNetworkInterface(ctx, input)
}

func (m *ENIManager) AttachENI(ctx context.Context, networkInterfaceID, instanceID string, deviceIndex int32) (*string, error) {
	input := &ec2.AttachNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(networkInterfaceID),
		InstanceId:         aws.String(instanceID),
		DeviceIndex:        aws.Int32(deviceIndex),
	}

	result, err := m.client.AttachNetworkInterface(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to attach ENI: %w", err)
	}

	return result.AttachmentId, nil
}

func (m *ENIManager) DetachENI(ctx context.Context, attachmentID string, force bool) error {
	input := &ec2.DetachNetworkInterfaceInput{
		AttachmentId: aws.String(attachmentID),
		Force:        aws.Bool(force),
	}

	_, err := m.client.DetachNetworkInterface(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to detach ENI: %w", err)
	}

	return nil
}

func (m *ENIManager) DeleteENI(ctx context.Context, networkInterfaceID string) error {
	input := &ec2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(networkInterfaceID),
	}

	_, err := m.client.DeleteNetworkInterface(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete ENI: %w", err)
	}

	return nil
}

func (m *ENIManager) ModifyENIAttribute(ctx context.Context, networkInterfaceID string, config ENIModifyConfig) error {
	input := &ec2.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(networkInterfaceID),
	}

	if config.Description != nil {
		input.Description = &types.AttributeValue{Value: config.Description}
	}

	if len(config.SecurityGroupIDs) > 0 {
		input.Groups = config.SecurityGroupIDs
	}

	_, err := m.client.ModifyNetworkInterfaceAttribute(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to modify ENI attribute: %w", err)
	}

	return nil
}

func (m *ENIManager) AssignPrivateIPs(ctx context.Context, networkInterfaceID string, count int32, specificIPs []string) error {
	input := &ec2.AssignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String(networkInterfaceID),
	}

	if count > 0 {
		input.SecondaryPrivateIpAddressCount = aws.Int32(count)
	}

	if len(specificIPs) > 0 {
		input.PrivateIpAddresses = specificIPs
	}

	_, err := m.client.AssignPrivateIpAddresses(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to assign private IPs: %w", err)
	}

	return nil
}

func (m *ENIManager) UnassignPrivateIPs(ctx context.Context, networkInterfaceID string, ips []string) error {
	input := &ec2.UnassignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String(networkInterfaceID),
		PrivateIpAddresses: ips,
	}

	_, err := m.client.UnassignPrivateIpAddresses(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to unassign private IPs: %w", err)
	}

	return nil
}

func (m *ENIManager) AssignIPv6Addresses(ctx context.Context, networkInterfaceID string, addresses []string, count *int32) error {
	input := &ec2.AssignIpv6AddressesInput{
		NetworkInterfaceId: aws.String(networkInterfaceID),
	}

	if len(addresses) > 0 {
		input.Ipv6Addresses = addresses
	}

	if count != nil {
		input.Ipv6AddressCount = count
	}

	_, err := m.client.AssignIpv6Addresses(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to assign IPv6 addresses: %w", err)
	}

	return nil
}

func (m *ENIManager) UnassignIPv6Addresses(ctx context.Context, networkInterfaceID string, addresses []string) error {
	input := &ec2.UnassignIpv6AddressesInput{
		NetworkInterfaceId: aws.String(networkInterfaceID),
		Ipv6Addresses:      addresses,
	}

	_, err := m.client.UnassignIpv6Addresses(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to unassign IPv6 addresses: %w", err)
	}

	return nil
}

func (m *ENIManager) DescribeENIs(ctx context.Context, filters []types.Filter) (*ec2.DescribeNetworkInterfacesOutput, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: filters,
	}

	return m.client.DescribeNetworkInterfaces(ctx, input)
}

func (m *ENIManager) DescribeSubnet(ctx context.Context, subnetID string) (*ec2.DescribeSubnetsOutput, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetID},
	}

	return m.client.DescribeSubnets(ctx, input)
}
