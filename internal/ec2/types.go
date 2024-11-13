package ec2

// ENIConfig represents configuration for creating a network interface
type ENIConfig struct {
	SubnetID         string
	Description      string
	SecurityGroupIDs []string
	PrivateIPCount   int32
	IPv6AddressCount int32
	Tags             map[string]string
}

// ENIModifyConfig represents configuration for modifying a network interface
type ENIModifyConfig struct {
	Description      *string
	SecurityGroupIDs []string
}
