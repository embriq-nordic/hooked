package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

// Handler is the lambda handler function
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Lambda invoked")
	defer log.Println("Lambda finished")

	return events.APIGatewayProxyResponse{}, nil
}
