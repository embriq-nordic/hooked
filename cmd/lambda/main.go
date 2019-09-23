package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rejlersembriq/hooked/pkg/lambdahandler"
	"github.com/rejlersembriq/hooked/pkg/repository/dynamo"
	"github.com/rejlersembriq/hooked/pkg/router"
	"github.com/rejlersembriq/hooked/pkg/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
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
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig = config

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	rtr = router.New()

	table := os.Getenv(tableName)
	region := os.Getenv(awsRegion)

	conf, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("unable to load SDK config, %s", err.Error())
	}
	conf.Region = region

	dyna = dynamo.New(dynamodb.New(conf), table)
}

func main() {
	lambda.Start(lambdahandler.Handler{
		Handler: server.New(rtr, dyna),
	}.Handle)
}
