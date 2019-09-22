package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rejlersembriq/hooked/pkg/lambdahandler"
	"github.com/rejlersembriq/hooked/pkg/repository/dynamo"
	"github.com/rejlersembriq/hooked/pkg/router"
	"github.com/rejlersembriq/hooked/pkg/server"
	"os"
)

// Environment variable names
const (
	tableName = "TABLE_NAME"
	awsRegion = "REGION"
)

var rtr *router.Router
var dyna *dynamo.Dynamo

func init() {
	rtr = router.New()

	table := os.Getenv(tableName)
	region := os.Getenv(awsRegion)

	conf, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	conf.Region = region

	dyna = dynamo.New(dynamodb.New(conf), table)
}

func main() {
	lambda.Start(lambdahandler.Handler{
		Handler: server.New(rtr, dyna),
	}.Handle)
}
