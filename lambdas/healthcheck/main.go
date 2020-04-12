package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// Response from API
type Response struct {
	StatusCode float64 `json:"statusCode"`
	Body       string  `json:"body"`
}

// Handler function for lambda
func Handler(context.Context) (Response, error) {
	return Response{
		StatusCode: 200,
		Body:       "Healthy",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
