package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	. "github.com/awslabs/aws-cfn-go-template"
	"github.com/sethvargo/go-fastly"
)

type SecurityGroupIngress struct {
	GroupId               interface{} `json:",omitempty"`
	SourceSecurityGroupId interface{} `json:",omitempty"`
	CidrIp                interface{} `json:",omitempty"`
	IpProtocol            string      `json:",omitempty"`
	FromPort              string      `json:",omitempty"`
	ToPort                string      `json:",omitempty"`
}

func main() {
	var name, protocol, port string
	flag.StringVar(&name, "name", "", "Name to use for this auto-generated security group.")
	flag.StringVar(&protocol, "protocol", "tcp", "The IP protocol name (tcp, udp, icmp) or number that these rules should apply to.")
	flag.StringVar(&port, "port", "80", "The port number or port range to allow.")
	flag.Parse()

	client, err := fastly.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	ips, err := client.IPs()
	if err != nil {
		log.Fatal(err)
	}

	t, err := GenTemplate(ips, name, protocol, port)

	//b, err := json.MarshalIndent(t, "", "  ")
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
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
