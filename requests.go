package main

import (
	"bytes"
	"encoding/json"
	"github.com/gookit/color"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultClient = &http.Client{Timeout: 10 * time.Second}
)

type (
	DiscriminatorCheckBody struct {
		Username      string `json:"username"`
		Password      string `json:"password"`
		Discriminator string `json:"discriminator"`
	}

	RatelimitResponse struct {
		Global     bool    `json:"global"`
		Message    string  `json:"message"`
		RetryAfter float64 `json:"retry_after"`
	}
)

// The CheckDiscriminator function checks if the discriminator is valid.
// The function returns true if the discriminator is valid, false if not.
// The function takes in a discord token string, username, and discriminator as an argument.
func CheckDiscriminator(discordToken, username string, discriminator int) (bool, error) {
	var body DiscriminatorCheckBody

	url := "https://discordapp.com/api/v9/users/@me"
	method := "PATCH"
	body = DiscriminatorCheckBody{
		Username:      username,
		Password:      "",
		Discriminator: strconv.Itoa(discriminator),
	}
	// Convert the body to JSON
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return false, err
	}

	// Create a new request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return false, err
	}

	// Set the headers
	req.Header.Set("Authorization", discordToken)

	// Send the request
DoRequest:
	resp, err := DefaultClient.Do(DiscordCommonHeaders(req))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Convert the body to a string
	bodyString := string(bodyBytes)

	// Check if the response is a ratelimit response
	if strings.Contains(bodyString, "You are being rate limited.") {
		// Get the retry after time
		var ratelimitResponse RatelimitResponse
		err = json.Unmarshal(bodyBytes, &ratelimitResponse)
		if err != nil {
			return false, err
		}
		// Wait for the retry after time
		color.Red.Println("[Discord] You are being rate limited. Application Waiting for " + strconv.FormatFloat(ratelimitResponse.RetryAfter, 'f', 2, 64) + " seconds.")
		time.Sleep(time.Duration(ratelimitResponse.RetryAfter) * time.Second)
		goto DoRequest
	}
	// Check if the discriminator is valid
	if !strings.Contains(bodyString, "This username and tag are already taken. Please try another.") {
		return true, nil
	}

	return false, nil

}

// The CheckToken function checks if the token is valid.
func CheckToken(auth string) (int, error) {
	url := "https://discord.com/api/v9/users/@me/affinities/guilds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red.Println("Error: " + err.Error())
		return -1, err
	}
	req.Header.Set("authorization", auth)
	resp, err := DefaultClient.Do(DiscordCommonHeaders(req))
	if err != nil {
		color.Red.Println("Error: " + err.Error())
		return -1, err
	}

	return resp.StatusCode, nil

}

func DiscordCommonHeaders(req *http.Request) *http.Request {

	req.Header.Set("host", "discord.com")
	req.Header.Set("x-super-properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDA2Iiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTQyNzUxLCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ")
	req.Header.Set("x-discord-locale", "en-US")
	req.Header.Set("x-debug-options", "bugReporterEnabled")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9005 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://discord.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://discord.com/channels/@me")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	return req
}
