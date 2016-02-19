package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/quiffman/gen-secgroup-cfn"
)

type agent struct {
	IpAddress string   `json:"ip_address"`
	Name      string   `json:"name"`
	Location  location `json:"location"`
}
type location struct {
	Country string `json:"country"`
	State   string `json:"state"`
	Name    string `json:"name"`
	City    string `json:"city"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var apiKey, name, vpcId, protocol, port string
	flag.StringVar(&apiKey, "api-key", "", "API Key with permissions to query the service.")
	flag.StringVar(&name, "name", "Webmon", "Name to use for this auto-generated security group.")
	flag.StringVar(&vpcId, "vpc", "", "The VPC ID to place this auto-generated security group in.")
	flag.StringVar(&protocol, "protocol", "tcp", "The IP protocol name (tcp, udp, icmp) or number that these rules should apply to.")
	flag.StringVar(&port, "port", "80", "The port number or port range to allow.")
	flag.Parse()

	var list []agent
	err := apiGet(apiKey, "https://webmon.com/api/v1/targets/agents/", url.Values{}, &list)
	check(err)

	var ips []string
	for _, i := range list {
		if net.ParseIP(i.IpAddress) != nil {
			ips = append(ips, fmt.Sprintf("%s/32", i.IpAddress))
		}
	}

	t, err := cfn.GenTemplate(ips, name, vpcId, protocol, port)
	check(err)

	//b, err := json.MarshalIndent(t, "", "  ")
	b, err := json.Marshal(t)
	check(err)
	os.Stdout.Write(b)
}

// apiGet issues a GET request to the Webmon API and decodes the response JSON to data.
func apiGet(apiKey string, urlStr string, form url.Values, data interface{}) error {
	req, err := http.NewRequest("GET", urlStr, nil)
	req.SetBasicAuth(apiKey, "")
	//d, err := httputil.DumpRequestOut(req, true)
	//fmt.Printf("%s\n\n", d)
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Webmon API.
func decodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != 200 {
		return newApiError(resp)
	}
	//d, _ := httputil.DumpResponse(resp, true)
	//fmt.Printf("%s\n\n", d)
	return json.NewDecoder(resp.Body).Decode(data)
}

type ApiError struct {
	StatusCode int
	Header     http.Header
	Body       string
	URL        *url.URL
}

func newApiError(resp *http.Response) *ApiError {
	// TODO don't ignore this error
	// TODO don't use ReadAll
	p, _ := ioutil.ReadAll(resp.Body)

	return &ApiError{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       string(p),
		URL:        resp.Request.URL,
	}
}

// ApiError supports the error interface
func (aerr ApiError) Error() string {
	return fmt.Sprintf("Get %s returned status %d, %s", aerr.URL, aerr.StatusCode, aerr.Body)
}

//  vim: set ts=4 sw=4 tw=0 et:
