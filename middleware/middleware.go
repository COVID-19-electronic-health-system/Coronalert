package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/COVID-19-electronic-health-system/Coronalert/models"
)

var phoneNumbers []models.Number

// HealthCheck simply returns a response that you successfully hit the API
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Coronalert Web Server")
	w.WriteHeader(200)
}

// StartPolling starts a cycle to send texts every 180 minutes
func StartPolling() {
	log.Println("starting notification service...")
	for {
		time.Sleep(180 * time.Minute)
		go SendSMS()
	}
}

// Subscribe subscribes a user to the message list
func Subscribe(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("subscribed to SMS service...")
}

// Unsubscribe unsubscribes a user from the notifications list
func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var number models.Number
	err := json.NewDecoder(r.Body).Decode(&number)
	if err != nil {
		panic(err)
	}

	// TOOD make this more elegant, perhaps use a map instead of a slice
	// NOTE we will likely eventually persist to a database, so not a huge worry for now
	for i := 0; i < len(phoneNumbers); i++ {
		if phoneNumbers[i] == number {
			phoneNumbers = append(phoneNumbers[:i], phoneNumbers[i+1])
		}
	}

	fmt.Println("unsubscribed from SMS service...")
}

// SendSMS sends update to subscribers
func SendSMS() {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	log.Println("requesting URL at", urlStr)

	notifications := [4]string{"Did you know? COVID-19 was first detected in Wuhan City, Hubei Province, China.",
		"Did you know? COVID-19 is not the same as typical coronaviruses commonly circulated amongst humans",
		"Did you know? Social distancing (avoiding large crowds) is the undeniable best way to prevent the spread of COVID-19",
		"If you experience symptoms of COVID-19, do not immediately head to the ER. Use CoronaTracker to monitor your symptoms and make a more informed decision"}

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
