package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoURI = os.Getenv("MONGODB_URI")
var twilioSID = os.Getenv("TWILIO_ACCOUNT_SID")
var twilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
var twilioPhoneNumber = os.Getenv("TWILIO_PHONE_NUMBER")

var twilioURL = "https://api.twilio.com/2010-04-01/Accounts/" + twilioSID + "/Messages.json"

var notifications = [4]string{
	"Did you know? COVID-19 was first detected in Wuhan City, Hubei Province, China.",
	"Did you know? COVID-19 is not the same as typical coronaviruses commonly circulated amongst humans",
	"Did you know? Social distancing (avoiding large crowds) is the undeniable best way to prevent the spread of COVID-19",
	"If you experience symptoms of COVID-19, do not immediately head to the ER. Use CoronaTracker to monitor your symptoms and make a more informed decision",
}

type phoneNumber struct {
	Number string `bson:"phoneNumber"`
}

func sendNotification(phoneNumber string) error {
	msgData := url.Values{}
	msgData.Set("To", phoneNumber)
	msgData.Set("From", twilioPhoneNumber)
	msgData.Set("Body", notifications[rand.Intn(len(notifications))])
	msgDataReader := *strings.NewReader(msgData.Encode())

	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", twilioURL, &msgDataReader)
	req.SetBasicAuth(twilioSID, twilioAuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := httpClient.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Failed to send to %s | Status code %d", phoneNumber, resp.StatusCode)
	}

	return nil
}

// Handler queries MongoDB for phone numbers and sends notifications to them
func Handler(ctx context.Context) (string, error) {
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Println("Connection error", err)
		return "", err
	}

	// Test connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("Ping error", err)
		return "", err
	}

	phoneNumbersCollection := client.Database("Coronalert").Collection("PhoneNumbers")

	// Get all records
	cur, err := phoneNumbersCollection.Find(ctx, bson.D{})
	if err != nil {
		return "", err
	}

	defer cur.Close(ctx)

	// For each record, send a notification
	for cur.Next(ctx) {
		var phoneNumber phoneNumber

		err := cur.Decode(&phoneNumber)
		if err != nil {
			log.Println(err)
		}

		err = sendNotification(phoneNumber.Number)
		if err != nil {
			log.Println(err)
		}
	}

	if err = cur.Err(); err != nil {
		log.Println("Cursor error", err)
		return "", err
	}

	return "success", nil
}

func main() {
	lambda.Start(Handler)
}
