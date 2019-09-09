package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rejlersembriq/hooked/pkg/handler"
)

func main() {
	lambda.Start(handler.Handler)
}
