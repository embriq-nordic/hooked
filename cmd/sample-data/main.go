package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rejlersembriq/hooked/pkg/dynamo"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"log"
	"math/rand"
	"strconv"
)

const (
	profile   = "larwef"
	region    = endpoints.EuWest1RegionID
	tableName = "hooked-participants"
)

const (
	noToGenerate = 10
	maxScore     = 100
)

const (
	compA comapany = iota
	compB
	compC
	compD
	compE

	noOfCompanies = 5
)

const (
	commentA comment = iota
	commentB
	commentC
	commentD

	noOfComments = 4
)

type comapany uint8

func (c comapany) String() string {
	switch c {
	case compA:
		return "CompanyA"
	case compB:
		return "CompanyB"
	case compC:
		return "CompanyC"
	case compD:
		return "CompanyD"
	case compE:
		return "CompanyE"
	default:
		panic(fmt.Sprintf("Invalid value for company: %d", int8(c)))
	}
}

type comment uint8

func (c comment) String() string {
	switch c {
	case commentA:
		return "Some comment."
	case commentB:
		return "Some other comment."
	case commentC:
		return ""
	case commentD:
		return "This is a slightly longer comment than the others."
	default:
		panic(fmt.Sprintf("Invalid value for comment: %d", int8(c)))
	}
}

func getPhoneNo(someInt int, len int) string {
	var phoneNo string

	for i := 0; i < len; i++ {
		phoneNo += strconv.Itoa(someInt % 10)
		someInt = someInt / 10
	}

	return phoneNo
}

func main() {
	conf, err := external.LoadDefaultAWSConfig(external.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("Message getting AWS config: %v", err)
	}
	conf.Region = region

	ddb := dynamodb.New(conf)

	repo := dynamo.New(ddb, tableName)

	populateWithTestData(repo)
	//deleteParticipant(repo, "03560e71-3904-4cab-a9f3-aa8c5be74a87")
	//deleteAll(repo)
}

func populateWithTestData(repo participant.Repository) {
	rand.Seed(42) // "To get the same randomness each time"

	for i := 0; i < noToGenerate; i++ {
		r := rand.Int()

		p := &participant.Participant{
			Email:   fmt.Sprintf("Participant%d@%s.com", i, comapany(r%noOfCompanies).String()),
			Name:    fmt.Sprintf("Participant%d", i),
			Phone:   getPhoneNo(r, 8),
			Org:     comapany(r % noOfCompanies).String(),
			Comment: comment(r % noOfComments).String(),
			Score:   r % maxScore,
		}

		saved, err := repo.Save(p)
		if err != nil {
			log.Fatalf("Error saving participant: %v", err)
		}

		bytes, _ := json.MarshalIndent(saved, "", "    ")
		fmt.Printf("Saved:\n%s\n", string(bytes))
	}
}

func deleteParticipant(repo participant.Repository, id string) {
	if err := repo.Delete(id); err != nil {
		log.Fatalf("Error deleting entry: %v", err)
	}
}

func deleteAll(repo participant.Repository) {
	participants, err := repo.GetAll()
	if err != nil {
		log.Fatalf("Error getting participants: %v", err)
	}

	for _, p := range participants {
		deleteParticipant(repo, p.ID)
	}
}
