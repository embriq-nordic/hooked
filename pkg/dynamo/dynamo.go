package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"time"
)

// Dynamo implements the participant Repository class to interact with DynamoDb.
type Dynamo struct {
	dynamoDb         dynamodbiface.ClientAPI
	participantTable string
}

// New returns a new dynamo repository.
func New(dynamoIface dynamodbiface.ClientAPI, tableName string) *Dynamo {
	return &Dynamo{
		dynamoDb:         dynamoIface,
		participantTable: tableName,
	}
}

// Save persists a participant to DynamoDb.
func (d *Dynamo) Save(p *participant.Participant) (*participant.Participant, error) {
	// If id is specified the object should exist in the table. Otherwise we expect it to not be present.
	condition := expression.ConditionBuilder{}
	if p.ID != "" {
		condition = expression.AttributeExists(expression.Name("id"))
	} else {
		p.ID = uuid.New().String()
		condition = expression.AttributeNotExists(expression.Name("id"))
	}

	update := expression.
		Set(expression.Name("created"), expression.IfNotExists(expression.Name("created"), expression.Value(time.Now().Unix()))).
		Set(expression.Name("updated"), expression.Value(time.Now().Unix()))

	// Split up to support partial updates and empty attributes.
	if p.Name != "" {
		update = update.Set(expression.Name("name"), expression.Value(p.Name))
	}

	if p.Email != "" {
		update = update.Set(expression.Name("email"), expression.Value(p.Email))
	}

	if p.Phone != "" {
		update = update.Set(expression.Name("phone"), expression.Value(p.Phone))
	}

	if p.Org != "" {
		update = update.Set(expression.Name("org"), expression.Value(p.Org))
	}

	if p.Score != 0 {
		update = update.Set(expression.Name("score"), expression.Value(p.Score))
	}

	if p.Comment != "" {
		update = update.Set(expression.Name("comment"), expression.Value(p.Comment))
	}

	exp, err := expression.NewBuilder().
		WithUpdate(update).
		WithCondition(condition).
		Build()
	if err != nil {
		return nil, err
	}

	res, err := d.dynamoDb.UpdateItemRequest(
		&dynamodb.UpdateItemInput{
			ConditionExpression:       exp.Condition(),
			ExpressionAttributeValues: exp.Values(),
			ExpressionAttributeNames:  exp.Names(),
			Key: map[string]dynamodb.AttributeValue{
				"id": {S: &p.ID},
			},
			ReturnValues:     dynamodb.ReturnValueAllNew,
			TableName:        &d.participantTable,
			UpdateExpression: exp.Update(),
		}).Send(context.Background())
	if err != nil {
		return nil, err
	}

	var savedParticipant participant.Participant
	if err := dynamodbattribute.UnmarshalMap(res.Attributes, &savedParticipant); err != nil {
		return nil, err
	}

	return &savedParticipant, err
}

// Get retrieves a participant from DynamoDb.
func (d *Dynamo) Get(id string) (*participant.Participant, error) {
	res, err := d.dynamoDb.GetItemRequest(
		&dynamodb.GetItemInput{
			Key: map[string]dynamodb.AttributeValue{
				"id": {S: &id},
			},
			TableName: &d.participantTable,
		}).Send(context.Background())

	if err != nil {
		return nil, err
	}

	var p participant.Participant
	if err := dynamodbattribute.UnmarshalMap(res.Item, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// GetAll retrieves all participants from DynamoDb.
func (d *Dynamo) GetAll() ([]*participant.Participant, error) {
	var result []*participant.Participant

	scanReq := d.dynamoDb.ScanRequest(&dynamodb.ScanInput{
		TableName: &d.participantTable,
	})

	paginator := dynamodb.NewScanPaginator(scanReq)

	// Paginator next returns false when finished or an error has occured
	for paginator.Next(context.Background()) {
		var recs []*participant.Participant
		page := paginator.CurrentPage()
		if err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &recs); err != nil {
			return nil, err
		}

		result = append(result, recs...)
	}

	if err := paginator.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete removes and entry matching the provided id.
func (d *Dynamo) Delete(id string) error {
	_, err := d.dynamoDb.DeleteItemRequest(
		&dynamodb.DeleteItemInput{
			ConditionExpression: aws.String("attribute_exists(id)"),
			Key: map[string]dynamodb.AttributeValue{
				"id": {S: &id},
			},
			TableName: &d.participantTable,
		}).Send(context.Background())

	return err
}
