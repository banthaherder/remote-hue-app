package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type auth struct {
	AccessToken     string `json:"access_token"`
	AccessTokenExp  string `json:"access_token_expires_in"`
	RefreshToken    string `json:"refresh_token"`
	RefreshTokenExp string `json:"refresh_token_expires_in"`
	TokenType       string `json:"token_type"`
}

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

func hashIt(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func genDigestResp(wwwAuth map[string]*string, clientID string, clientSecret string) string {
	hash1 := clientID + ":" + *wwwAuth["realm"] + ":" + clientSecret
	hash1 = hashIt(hash1)

	hash2 := hashIt("POST:/oauth2/token")

	digestResp := hash1 + ":" + *wwwAuth["nonce"] + ":" + hash2
	digestResp = hashIt(digestResp)

	digestAuth := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%s"`,
		clientID, *wwwAuth["realm"], *wwwAuth["nonce"], "/oauth2/token", digestResp)

	return digestAuth
}

// Make an http request
func newReq(verb string, reqURL string, authHeader string) *http.Response {
	req, err := http.NewRequest(verb, reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// "Do" the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	return resp
}

func tokenRequest(authCode string, clientID string, clientSecret string) {
	// Construct URL to retrieve a nonce
	tokenReqURL := "https://api.meethue.com/oauth2/token?code=" + authCode + "&grant_type=authorization_code"
	// POST and expect a 401 unauthorized response with a nonce
	nonceResp := newReq("POST", tokenReqURL, "")
	// Map the nonce response values to retrieve the nonce
	wwwAuthMap := processHeader(nonceResp, "www-authenticate")
	if wwwAuthMap["nonce"] == nil || wwwAuthMap["realm"] == nil {
		log.Fatalln("www-Authorization header incomplete...")
	}
	// Generate the md5 handshake using the nonce and secrets
	authHeader := genDigestResp(wwwAuthMap, clientID, clientSecret)
	digestResp := newReq("POST", tokenReqURL, authHeader)

	body, readErr := ioutil.ReadAll(digestResp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	authResp := auth{}
	jsonErr := json.Unmarshal(body, &authResp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(authResp)

}

func getEnv(env string) string {
	if val, ok := os.LookupEnv(env); ok {
		return val
	}
	return ""
}

func main() {
	// Will be recieved by function from params
	code := "abc123"

	clientID := getEnv("HUE_CLIENT_ID")
	if clientID == "" {
		log.Fatalln("No client id present")
	}
	clientSecret := getEnv("HUE_CLIENT_SECRET")
	if clientID == "" {
		log.Fatalln("No client secret present")
	}

	tokenRequest(code, clientID, clientSecret)
}
