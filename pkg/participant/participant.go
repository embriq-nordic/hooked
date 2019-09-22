package participant

import (
	"errors"
	"time"
)

// Error defines error type for participant related actions.
// Not strictly neccessary, but might com in handy later.
type Error error

// ErrNotExist should be returned when the requested participant is not found in the repository.
var ErrNotExist Error = errors.New("participant doesn't exist")

// Repository defines interface for persisting and retrieving participant info.
type Repository interface {
	Save(participant Participant) (*Participant, Error)
	Get(id string) (*Participant, Error)
	GetAll() ([]*Participant, Error)
	Delete(id string) Error
}

// Participant represents a participants object.
type Participant struct {
	ID      *string    `json:"id,omitempty"`
	Name    *string    `json:"name,omitempty"`
	Email   *string    `json:"email,omitempty"`
	Phone   *string    `json:"phone,omitempty"`
	Org     *string    `json:"org,omitempty"`
	Score   *int       `json:"score"`
	Comment *string    `json:"comment,omitempty"`
	Created *time.Time `json:"created" dynamodbav:",unixtime"`
	Updated *time.Time `json:"updated" dynamodbav:",unixtime"`
}
