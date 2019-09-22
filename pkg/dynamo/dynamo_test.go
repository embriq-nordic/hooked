package dynamo

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"net/http"
	"testing"
)

// Mock
type dynamodbMock struct {
	dynamodbiface.ClientAPI

	getItemRequestHandler    func(*dynamodb.GetItemInput) dynamodb.GetItemRequest
	updateItemRequestHandler func(*dynamodb.UpdateItemInput) dynamodb.UpdateItemRequest
	deleteItemRequestHandler func(*dynamodb.DeleteItemInput) dynamodb.DeleteItemRequest
}

func (d dynamodbMock) GetItemRequest(input *dynamodb.GetItemInput) dynamodb.GetItemRequest {
	return d.getItemRequestHandler(input)
}

func (d dynamodbMock) UpdateItemRequest(input *dynamodb.UpdateItemInput) dynamodb.UpdateItemRequest {
	return d.updateItemRequestHandler(input)
}

func (d dynamodbMock) DeleteItemRequest(input *dynamodb.DeleteItemInput) dynamodb.DeleteItemRequest {
	return d.deleteItemRequestHandler(input)
}

// Tests
func TestDynamo_Get_NotExist(t *testing.T) {
	mock := dynamodbMock{
		getItemRequestHandler: func(input *dynamodb.GetItemInput) dynamodb.GetItemRequest {

			return dynamodb.GetItemRequest{
				Request: &aws.Request{
					Data:        &dynamodb.GetItemOutput{},
					HTTPRequest: &http.Request{},
				},
			}
		},
	}

	repo := New(mock, "test-table")

	if _, err := repo.Get(""); !errors.Is(err, participant.ErrNotExist) {
		t.Errorf("Got unexpected error %v", err)
	}
}

func TestDynamo_Delete_NotExist(t *testing.T) {
	mock := dynamodbMock{
		deleteItemRequestHandler: func(input *dynamodb.DeleteItemInput) dynamodb.DeleteItemRequest {

			return dynamodb.DeleteItemRequest{
				Request: &aws.Request{
					HTTPRequest: &http.Request{},
					Error:       awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "", nil),
				},
			}
		},
	}

	repo := New(mock, "test-table")

	if err := repo.Delete(""); !errors.Is(err, participant.ErrNotExist) {
		t.Errorf("Got unexpected error %v", err)
	}
}
