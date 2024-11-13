package ec2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_ec2.go -package=mocks . EC2ClientAPI

type EC2ClientAPI interface {
	CreateNetworkInterface(ctx context.Context, input *ec2.CreateNetworkInterfaceInput, opts ...func(*ec2.Options)) (*ec2.CreateNetworkInterfaceOutput, error)
	DescribeInstances(ctx context.Context, input *ec2.DescribeInstancesInput, opts ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	DescribeInstanceTypes(ctx context.Context, input *ec2.DescribeInstanceTypesInput, opts ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error)
	AttachNetworkInterface(ctx context.Context, input *ec2.AttachNetworkInterfaceInput, opts ...func(*ec2.Options)) (*ec2.AttachNetworkInterfaceOutput, error)
	DeleteNetworkInterface(ctx context.Context, input *ec2.DeleteNetworkInterfaceInput, opts ...func(*ec2.Options)) (*ec2.DeleteNetworkInterfaceOutput, error)
	DetachNetworkInterface(ctx context.Context, input *ec2.DetachNetworkInterfaceInput, opts ...func(*ec2.Options)) (*ec2.DetachNetworkInterfaceOutput, error)
	AssignPrivateIpAddresses(ctx context.Context, input *ec2.AssignPrivateIpAddressesInput, opts ...func(*ec2.Options)) (*ec2.AssignPrivateIpAddressesOutput, error)
	UnassignPrivateIpAddresses(ctx context.Context, input *ec2.UnassignPrivateIpAddressesInput, opts ...func(*ec2.Options)) (*ec2.UnassignPrivateIpAddressesOutput, error)
	AssignIpv6Addresses(ctx context.Context, input *ec2.AssignIpv6AddressesInput, opts ...func(*ec2.Options)) (*ec2.AssignIpv6AddressesOutput, error)
	UnassignIpv6Addresses(ctx context.Context, input *ec2.UnassignIpv6AddressesInput, opts ...func(*ec2.Options)) (*ec2.UnassignIpv6AddressesOutput, error)
	DescribeNetworkInterfaces(ctx context.Context, input *ec2.DescribeNetworkInterfacesInput, opts ...func(*ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error)
	ModifyNetworkInterfaceAttribute(ctx context.Context, input *ec2.ModifyNetworkInterfaceAttributeInput, opts ...func(*ec2.Options)) (*ec2.ModifyNetworkInterfaceAttributeOutput, error)
	CreateTags(ctx context.Context, input *ec2.CreateTagsInput, opts ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
	DescribeSubnets(ctx context.Context, input *ec2.DescribeSubnetsInput, opts ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
}
