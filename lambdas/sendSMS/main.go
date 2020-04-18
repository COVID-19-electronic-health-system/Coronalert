package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var notifications = [15]string{
	"Did you know? COVID-19 was first detected in Wuhan City, Hubei Province, China.",
	"Did you know? COVID-19 is not the same as typical coronaviruses commonly circulated amongst humans.",
	"Did you know? Social distancing (avoiding large crowds) is the undeniable best way to prevent the spread of COVID-19.",
	"Did you know? In COVID-19, 'CO' stands for 'corona', 'VI' for 'virus', and 'D' for disease.",
	"Did you know? The official name of the virus is SARS-CoV-2 (Severe Acute Respiratory Syndrome Coronavirus 2). The original SARS-CoV emerged in 2002 and was responsible for the SARS outbreak. While the death rate of SARS was higher than COVID-19, it infected and killed far fewer people. There have been no cases of SARS since 2003.",
	"Did you know? While SARS-CoV-2 originated in a wild animal market in Wuhan, China, the risk of catching a new coronavirus from animals remains very low. However, care should always be taken when handling raw meat, milk, or animal organs. Always ensure foods are properly cooked before consuming.",
	"Did you know? There is only one documented case of a dog being infected with SARS-CoV-2, so there is little evidence that pets can spread the disease. However, you should still wash your hands after handling pets, especially after walking your dog. Dogs could still contact contaminated surfaces and carry the virus externally.",
	"Did you know? It is safe to receive packages from areas where COVID-19 has been reported. The likelihood of an infected person contaminating commercial goods is low and the risk of catching the virus that causes COVID-19 from a package that has been moved, travelled, and exposed to different conditions and temperature is also low.",
	"Did you know? In order to fight a virus, your body’s immune system must be able to recognize it and destroy it. Viruses cannot be “killed” with medication because they are not alive.",
	"Did you know? Viruses are not living organisms. Unlike bacteria, viruses are unable to replicate and live on their own. They must infect a host to survive and reproduce.",
	"Did you know? The coronavirus is roughly 120nm in diameter. For reference, you could fit about 2,500 virus particles end to end on a single grain of salt!",
	"Did you know? Around 1 out of every 6 people who gets COVID-19 becomes seriously ill and develops difficulty breathing. Older people, and those with underlying medical problems like high blood pressure, heart problems or diabetes, are more likely to develop serious illness. People with fever, cough and difficulty breathing should seek medical attention.",
	"Did you know? Studies to date suggest that COVID-19 is mainly transmitted through contact with respiratory droplets rather than through the air.",
	"Did you know? The incubation period for COVID-19 ranges from 1 to 14 days, meaning it can take this long to start showing symptoms after initial infection. Some studies have shown that you are the most contagious during the early stages of the disease. This is why social distancing is so important to prevent asymptomatic transmission.",
	"If you experience symptoms of COVID-19, do not immediately head to the ER. Use CoronaTracker to monitor your symptoms and make a more informed decision",
}

var mongoURI = os.Getenv("MONGODB_URI")
var twilioSID = os.Getenv("TWILIO_ACCOUNT_SID")
var twilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
var twilioPhoneNumber = os.Getenv("TWILIO_PHONE_NUMBER")

var twilioURL string

type phoneNumber struct {
	Number string `bson:"phoneNumber"`
}

func init() {
	mongoURI = decrypt(mongoURI)
	twilioSID = decrypt(twilioSID)
	twilioAuthToken = decrypt(twilioAuthToken)
	twilioURL = "https://api.twilio.com/2010-04-01/Accounts/" + twilioSID + "/Messages.json"
}

func decrypt(encrypted string) string {
	kmsClient := kms.New(session.New())
	decodedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		panic(err)
	}
	input := &kms.DecryptInput{
		CiphertextBlob: decodedBytes,
	}
	response, err := kmsClient.Decrypt(input)
	if err != nil {
		panic(err)
	}
	// Plaintext is a byte array, so convert to string
	return string(response.Plaintext[:])
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

	// For different seed on every execution
	rand.Seed(time.Now().UnixNano())

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
