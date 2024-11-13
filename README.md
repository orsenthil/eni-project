# ENI Management Program

### Define the parameters

```
	subnetID := "subnet-0a7bd03887dc3cbd5"
	instanceID := "i-04890aa7cd8cf81f3"
	securityGroupID := "sg-0f9acdf364ab834f2"
```

## Steps to create a security group for ENI testing

```
aws ec2 describe-vpcs     --query 'Vpcs[*].[VpcId,Tags[?Key==`Name`].Value|[0],CidrBlock]'     --output table
```

```
export VPC_ID=
```

```
aws ec2 create-security-group     --group-name "eni-test-sg"     --description "Security group for ENI testing"     --vpc-id "$VPC_ID"     --output text     --query 'GroupId'
```

```
export SG_ID=
```

```
curl -s https://checkip.amazonaws.com
```

```
aws ec2 authorize-security-group-ingress     --group-id "$SG_ID"     --protocol tcp     --port 22     --cidr "54.203.107.106/32"
```

```
export VPC_CIDR=
```

```
aws ec2 authorize-security-group-ingress     --group-id "$SG_ID"     --protocol all     --port -1     --cidr "$VPC_CIDR"
```

```
aws ec2 describe-security-groups     --group-ids "$SG_ID"     --query 'SecurityGroups[0].[GroupId,GroupName,Description]'     --output table
```

