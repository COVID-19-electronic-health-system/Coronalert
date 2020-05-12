package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type phoneNumber struct {
	Number string `json:"phoneNumber"`
}

var mongoDBURI = os.Getenv("MONGODB_URI")

func init() {
	mongoDBURI = decrypt(mongoDBURI)
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

// Handler function for lambda
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//				START MONGODB SETUP					//
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		log.Println("Connection error:", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("Ping error:", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	coronalertDB := client.Database("Coronalert")
	phoneNumbersCollection := coronalertDB.Collection("PhoneNumbers")
	//				END MONGODB SETUP					//

	//				START SUBSCRIBE						//
	requestBody := phoneNumber{}

	err = json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Println("error in unmarshal")
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	document := bson.D{
		{Key: "_id", Value: requestBody.Number},
		{Key: "phoneNumber", Value: requestBody.Number},
		{Key: "lastIndex", Value: 0},
	}

	phoneNumber := phoneNumber{}
	err = phoneNumbersCollection.FindOne(ctx, document).Decode(&phoneNumber)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("error checking for phone number in collection")
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}
	if err == nil {
		log.Println("phone number already subscribed")
		return events.APIGatewayProxyResponse{
			StatusCode: 409,
		}, nil
	}

	_, err = phoneNumbersCollection.InsertOne(ctx, document)
	if err != nil {
		log.Println("error adding to collection")
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	response, err := json.Marshal(&requestBody)
	if err != nil {
		log.Println("error in marshal")
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}
	//				END SUBSCRIBE						//

	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
