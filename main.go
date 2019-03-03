package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Format and return a given header's values
func processHeader(resp *http.Response, targetHeader string) map[string]*string {
	header := resp.Header.Get(targetHeader)
	parts := strings.SplitN(header, " ", 2)
	parts = strings.Split(parts[1], ", ")
	fmt.Println("Parts: ", parts) // Print www-authenticate values
	headerVals := make(map[string]*string)

	for _, part := range parts {
		fmt.Println("Part: ", part) // Print individual values, including nonce
		vals := strings.SplitN(part, "=", 2)
		key := vals[0]
		val := strings.Trim(vals[1], "\",")
		headerVals[key] = &val
	}

	return headerVals
}

// Make an http request
func newReq(verb string, reqURL string) *http.Response {
	req, err := http.NewRequest(verb, reqURL, strings.NewReader(""))
	if err != nil {
		log.Fatalln(err)
	}

	// "Do" the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	return resp
}

func tokenRequest(code string) {
	// Construct URL to retrieve a nonce
	nonceReqURL := "https://api.meethue.com/oauth2/token?code=" + code + "&grant_type=authorization_code"
	// POST and expect a 401 unauthorized response with a nonce
	nonceResp := newReq("POST", nonceReqURL)
	// Map the nonce response values to retrieve the nonce
	wwwAuthMap := processHeader(nonceResp, "www-authenticate")

}

func main() {
	code := "abc123"

	tokenRequest(code)
}
