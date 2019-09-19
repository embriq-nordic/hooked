package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rejlersembriq/hooked/pkg/lambdahandler"
	"github.com/rejlersembriq/hooked/pkg/router"
	"github.com/rejlersembriq/hooked/pkg/server"
)

func main() {
	lambda.Start(lambdahandler.Handler{
		Handler: server.New(router.New(), nil),
	}.Handle)
}
