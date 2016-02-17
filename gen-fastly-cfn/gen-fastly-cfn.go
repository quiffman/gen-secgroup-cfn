package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/quiffman/gen-secgroup-cfn"
	"github.com/sethvargo/go-fastly"
)

func main() {
	var name, protocol, port string
	flag.StringVar(&name, "name", "Fastly", "Name to use for this auto-generated security group.")
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

	t, err := cfn.GenTemplate(ips, name, protocol, port)

	//b, err := json.MarshalIndent(t, "", "  ")
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}

//  vim: set ts=4 sw=4 tw=0 et:
