package participant

import "time"

// Repository defines interface for persisting and retrieving participant info.
type Repository interface {
	Save(participant *Participant) (*Participant, error)
	Get(id string) (*Participant, error)
	GetAll() ([]*Participant, error)
}

// Participant represents a participants object.
type Participant struct {
	Id      string    `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Email   string    `json:"email,omitempty"`
	Phone   string    `json:"phone,omitempty"`
	Org     string    `json:"org,omitempty"`
	Score   int       `json:"score"`
	Comment string    `json:"comment,omitempty"`
	Created time.Time `json:"created" dynamodbav:",unixtime"`
	Updated time.Time `json:"updated" dynamodbav:",unixtime"`
}
