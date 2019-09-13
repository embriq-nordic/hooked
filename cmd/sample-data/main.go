package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rejlersembriq/hooked/pkg/dynamo"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"log"
)

const (
	profile   = "larwef"
	region    = endpoints.EuWest1RegionID
	tableName = "hooked-participants"
)

func main() {
	conf, err := external.LoadDefaultAWSConfig(external.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("Message getting AWS config: %v", err)
	}
	conf.Region = region

	ddb := dynamodb.New(conf)

	repo := dynamo.New(ddb, tableName)

	p := &participant.Participant{
		//Id: "8c9c7a76-3d12-4f2c-8930-220e31033017",
		Email: "larwef@gmail.com",
		Name:  "Lars Wefald",
		Phone: "12345678",
		Org:   "Rejlers Embriq",
		Comment: "test",
		Score: 0,
	}

	if err := repo.Save(p); err != nil {
		log.Fatalf("Error saving participant: %v", err)
	}
}
