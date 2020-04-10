package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Response from API
type Response struct {
	StatusCode float64 `json:"statusCode"`
	Body       string  `json:"body"`
}

// Handler function for lambda
func Handler(ctx context.Context) (Response, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		os.Getenv("MONGODB_URI"),
	))
	if err != nil {
		log.Fatal("Connection error:", err)
		return Response{
			StatusCode: 500,
		}, err
	}

	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Ping error:", err)
		return Response{
			StatusCode: 500,
		}, err
	}

	log.Println("Successfully reached MongoDB")

	// Success
	return Response{
		StatusCode: 200,
		Body:       "Successfully pinged PhoneNumbers",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
