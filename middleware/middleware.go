package middleware

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"../models"
)

var phoneNumbers []models.Number

// SendSMS sends update to subscribers
func SendSMS(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO maybe remove when deploying (but behind API gateway soo...?)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // TODO maybe remove when deploying

	var number models.Number
	err := json.NewDecoder(r.Body).Decode(&number)

	if err != nil {
		panic(err)
	}

	phoneNumbers = append(phoneNumbers, number)

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	notifications := [2]string{"Test Notification 1",
		"Test Notification 2"}

	for i := 0; i < len(phoneNumbers); i++ {
		msgData := url.Values{}
		msgData.Set("To", phoneNumbers[i].Number)
		msgData.Set("From", "12018013744")
		msgData.Set("Body", notifications[rand.Intn(len(notifications))])
		msgDataReader := *strings.NewReader(msgData.Encode())

		client := &http.Client{}
		req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, _ := client.Do(req)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var data map[string]interface{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)
			if err == nil {
				log.Println(data["sid"])
			} else {
				log.Println(resp.Status)
			}
		} else {
			log.Println(resp.Status)
		}
	}
}
