package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func sendSMS(toUser string, fromUser string, msg string) bool {
	accountSid := "AC6861dbd95a4542d1be1c0a219d7cecec"
	authToken := ""
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// check if user exists in our directory
	if _, ok := directory[toUser]; !ok {
		fmt.Printf("User %s not found\n", toUser)
		return false
	}

	// build our data for our message
	v := url.Values{}
	v.Set("To", directory[toUser])
	v.Set("From", "+14159961749")
	v.Set("Body", msg+"\n- AMBot ("+fromUser+")")
	//v.Set("MediaUrl", "https://www.shitpostbot.com/img/sourceimages/hedgey-is-a-ok-57f2d71092a97.jpeg")

	rb := *strings.NewReader(v.Encode())

	// create and set parameters for the http client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make the request
	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
	return true
}
