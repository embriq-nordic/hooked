package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	apiURL = "<url>"
)

const (
	noToGenerate = 100
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

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	populateWithTestData(client)
	//deleteParticipant(client, "0d42191f-0284-4681-bbbd-e4316f5b8857")
	//deleteAll(client)
}

func populateWithTestData(client *http.Client) {
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

		paylaod, _ := json.Marshal(p)

		req, err := http.NewRequest(http.MethodPost, apiURL+"/participant", bytes.NewBuffer(paylaod))

		res, err := client.Do(req)
		if (err != nil) || (res.StatusCode < 200 || res.StatusCode > 299) {
			log.Fatalf("Error during participant POST. Error: %v, Status: %d", err, res.StatusCode)
		}

		b, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("Saved:\n%s\n", string(b))

		res.Body.Close()
	}
}

func deleteParticipant(client *http.Client, id string) {
	req, err := http.NewRequest(http.MethodDelete, apiURL+"/participant/"+id, nil)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error during participant DELETE. Error: %v", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Fatalf("Error during participant DELETE. Status: %d", res.StatusCode)
	}
}

func deleteAll(client *http.Client) {
	req, err := http.NewRequest(http.MethodGet, apiURL+"/participants", nil)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error during participants GET. Error: %v", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Fatalf("Error during participants GET. Status: %d", res.StatusCode)
	}

	defer res.Body.Close()
	var participants []*participant.Participant
	if err := json.NewDecoder(res.Body).Decode(&participants); err != nil {
		log.Fatalf("Error getting participants response: %v", err)
	}

	for _, p := range participants {
		deleteParticipant(client, p.ID)
	}
}
