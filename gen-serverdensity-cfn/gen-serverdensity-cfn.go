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
	"strings"

	"github.com/quiffman/gen-secgroup-cfn"
	"golang.org/x/net/html"
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
	var name, vpcId, protocol, port, codes string
	flag.StringVar(&name, "name", "ServerDensity", "Name to use for this auto-generated security group.")
	flag.StringVar(&vpcId, "vpc", "", "The VPC ID to place this auto-generated security group in.")
	flag.StringVar(&protocol, "protocol", "tcp", "The IP protocol name (tcp, udp, icmp) or number that these rules should apply to.")
	flag.StringVar(&port, "port", "80", "The port number or port range to allow.")
	flag.StringVar(&codes, "codes", "", "The comma separated list of ServerDensity location codes to add rules for, or all if left empty.")
	flag.Parse()

	ips, err := scrapePage("https://support.serverdensity.com/hc/en-us/articles/201091476-Monitoring-node-locations-and-IP-addresses", codes)
	check(err)

	t, err := cfn.GenTemplate(ips, name, vpcId, protocol, port)
	check(err)

	//b, err := json.MarshalIndent(t, "", "  ")
	b, err := json.Marshal(t)
	check(err)
	os.Stdout.Write(b)
}

func scrapePage(url string, codes string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, newApiError(resp)
	}

	var ips []string

	z := html.NewTokenizer(resp.Body)
	defer resp.Body.Close()

	var code string
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return ips, nil
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "code" {
				z.Next()

				code = z.Token().Data
			}

			if t.Data == "li" {
				//fmt.Println(t)
				for _, a := range t.Attr {
					//fmt.Println(a)
					if a.Key == "class" && a.Val == "ip" {
						z.Next()

						i := z.Token().Data
						if code != "" && strings.Contains(codes, code) {
							if net.ParseIP(i) != nil {
								ips = append(ips, fmt.Sprintf("%s/32", i))
							}
						}
					}
				}
			}
		}
	}
	return ips, nil
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

//  vim: set ts=4 sw=4 tw=0 :
