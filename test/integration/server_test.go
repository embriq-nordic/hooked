// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"
)

const (
	apiURL = "http://localhost:8081"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 1 * time.Minute,
	}
}

func TestIntegration_PostGetPutAndDelete(t *testing.T) {
	var p *participant.Participant
	p = PostTest(t)
	p = GetTest(t, p)
	p = PutTest(t, p)
	p = GetTest(t, p)
	p = PutTestPartialUpdate(t, p)
	p = GetTest(t, p)
	DeleteTest(t, p)
}

func TestIntegration_GET_NotExist(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, apiURL+"/participant/nonExisting", nil)
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestIntegration_PUT_NotExist(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, apiURL+"/participant/nonExisting", bytes.NewBufferString("{}"))
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestIntegration_DELETE_NotExist(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, apiURL+"/participant/nonExisting", nil)
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

// Helpers
func PostTest(t *testing.T) *participant.Participant {
	pReq := &participant.Participant{
		Name:    aws.String("Participant1"),
		Email:   aws.String("participant1@participant.com"),
		Phone:   aws.String("11111111"),
		Org:     aws.String("Org1"),
		Score:   aws.Int(1),
		Comment: aws.String("Comment 1"),
	}

	payload, err := json.Marshal(pReq)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, apiURL+"/participant", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	pRes := doReq(t, req)

	assert.NotNil(t, pRes.ID)
	assert.Equal(t, *pReq.Name, *pRes.Name)
	assert.Equal(t, *pReq.Email, *pRes.Email)
	assert.Equal(t, *pReq.Phone, *pRes.Phone)
	assert.Equal(t, *pReq.Org, *pRes.Org)
	assert.Equal(t, *pReq.Score, *pRes.Score)
	assert.Equal(t, *pReq.Comment, *pRes.Comment)

	return pRes
}

func GetTest(t *testing.T, p *participant.Participant) *participant.Participant {
	req, err := http.NewRequest(http.MethodGet, apiURL+"/participant/"+*p.ID, nil)
	assert.NoError(t, err)

	pRes := doReq(t, req)

	assert.Equal(t, *p.ID, *pRes.ID)
	assert.Equal(t, *p.Name, *pRes.Name)
	assert.Equal(t, *p.Email, *pRes.Email)
	assert.Equal(t, *p.Phone, *pRes.Phone)
	assert.Equal(t, *p.Org, *pRes.Org)
	assert.Equal(t, *p.Score, *pRes.Score)
	assert.Equal(t, *p.Comment, *pRes.Comment)
	assert.Equal(t, *p.Created, *pRes.Created)
	assert.Equal(t, *p.Updated, *pRes.Updated)

	return pRes
}

func PutTest(t *testing.T, p *participant.Participant) *participant.Participant {
	pReq := &participant.Participant{
		Name:    aws.String("Participant2"),
		Email:   aws.String("participant2@participant.com"),
		Phone:   aws.String("22222222"),
		Org:     aws.String("Org2"),
		Score:   aws.Int(2),
		Comment: aws.String("Comment 2"),
	}

	payload, err := json.Marshal(pReq)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, apiURL+"/participant/"+*p.ID, bytes.NewBuffer(payload))
	assert.NoError(t, err)

	pRes := doReq(t, req)

	assert.Equal(t, *pReq.Name, *pRes.Name)
	assert.Equal(t, *pReq.Email, *pRes.Email)
	assert.Equal(t, *pReq.Phone, *pRes.Phone)
	assert.Equal(t, *pReq.Org, *pRes.Org)
	assert.Equal(t, *pReq.Score, *pRes.Score)
	assert.Equal(t, *pReq.Comment, *pRes.Comment)
	assert.Equal(t, *p.Updated, *pRes.Created)

	return pRes
}

func PutTestPartialUpdate(t *testing.T, p *participant.Participant) *participant.Participant {
	pReq := &participant.Participant{
		Score:   aws.Int(3),
		Comment: aws.String("Comment 3"),
	}

	payload, err := json.Marshal(pReq)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, apiURL+"/participant/"+*p.ID, bytes.NewBuffer(payload))
	assert.NoError(t, err)

	pRes := doReq(t, req)

	assert.Equal(t, *p.ID, *pRes.ID)
	assert.Equal(t, *p.Name, *pRes.Name)
	assert.Equal(t, *p.Email, *pRes.Email)
	assert.Equal(t, *p.Phone, *pRes.Phone)
	assert.Equal(t, *p.Org, *pRes.Org)
	assert.Equal(t, *pReq.Score, *pRes.Score)
	assert.Equal(t, *pReq.Comment, *pRes.Comment)
	assert.True(t, p.Updated.After(*pRes.Created))

	return pRes
}

func DeleteTest(t *testing.T, p *participant.Participant) {
	req, err := http.NewRequest(http.MethodDelete, apiURL+"/participant/"+*p.ID, nil)
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Should not exist after DELETE
	req, err = http.NewRequest(http.MethodGet, apiURL+"/participant/"+*p.ID, nil)
	assert.NoError(t, err)

	res, err = client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func doReq(t *testing.T, req *http.Request) *participant.Participant {
	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Fatalf("Error during %s request. Status: %d", req.Method, res.StatusCode)
	}

	var pRes participant.Participant
	err = json.NewDecoder(res.Body).Decode(&pRes)
	assert.NoError(t, err)

	return &pRes
}
