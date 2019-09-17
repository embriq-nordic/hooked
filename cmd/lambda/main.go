package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rejlersembriq/hooked/pkg/lambdahandler"
	"github.com/rejlersembriq/hooked/pkg/participant"
)

func main() {
	lambda.Start(lambdahandler.Handler{
		Handler: &participant.Handler{},
	}.Handle)
}
