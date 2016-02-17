package cfn

import (
	"fmt"

	. "github.com/awslabs/aws-cfn-go-template"
)

type SecurityGroupIngress struct {
	GroupId               interface{} `json:",omitempty"`
	SourceSecurityGroupId interface{} `json:",omitempty"`
	CidrIp                interface{} `json:",omitempty"`
	IpProtocol            string      `json:",omitempty"`
	FromPort              string      `json:",omitempty"`
	ToPort                string      `json:",omitempty"`
}

func GenTemplate(ips []string, name string, protocol string, port string) (Template, error) {
	var s []SecurityGroupIngress

	for _, i := range ips {
		s = append(s, SecurityGroupIngress{
			IpProtocol: protocol,
			FromPort:   port,
			ToPort:     port,
			CidrIp:     i,
		})
	}

	var p = make(map[string]interface{})
	p["GroupDescription"] = fmt.Sprintf("allow %s/%s connections from specified %s CIDR ranges", protocol, port, name)
	p["SecurityGroupIngress"] = s

	t := Template{
		AWSTemplateFormatVersion: "2010-09-09",

		Description: fmt.Sprintf("Auto-generated %s security group", name),
		Resources: map[string]Resource{
			"ServerSecurityGroup": Resource{
				Type:       "AWS::EC2::SecurityGroup",
				Properties: p,
			},
		},
	}

	return t, nil
}

//  vim: set ts=4 sw=4 tw=0 et:
