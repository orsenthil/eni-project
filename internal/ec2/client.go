package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_ec2.go -package=mocks . EC2ClientAPI

type EC2ClientAPI interface {
	CreateNetworkInterface(ctx context.Context, input *ec2.CreateNetworkInterfaceInput, opts ...func(*ec2.Options)) (*ec2.CreateNetworkInterfaceOutput, error)
	AttachNetworkInterface(ctx context.Context, input *ec2.AttachNetworkInterfaceInput, opts ...func(*ec2.Options)) (*ec2.AttachNetworkInterfaceOutput, error)
	DeleteNetworkInterface(ctx context.Context, input *ec2.DeleteNetworkInterfaceInput, opts ...func(*ec2.Options)) (*ec2.DeleteNetworkInterfaceOutput, error)
	DetachNetworkInterface(ctx context.Context, input *ec2.DetachNetworkInterfaceInput, opts ...func(*ec2.Options)) (*ec2.DetachNetworkInterfaceOutput, error)
}
