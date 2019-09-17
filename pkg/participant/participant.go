package participant

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Repository defines interface for persisting and retrieving participant info.
type Repository interface {
	Save(participant *Participant) (*Participant, error)
	Get(id string) (*Participant, error)
	GetAll() ([]*Participant, error)
	Delete(id string) error
}

// Participant represents a participants object.
type Participant struct {
	ID      string    `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Email   string    `json:"email,omitempty"`
	Phone   string    `json:"phone,omitempty"`
	Org     string    `json:"org,omitempty"`
	Score   int       `json:"score"`
	Comment string    `json:"comment,omitempty"`
	Created time.Time `json:"created" dynamodbav:",unixtime"`
	Updated time.Time `json:"updated" dynamodbav:",unixtime"`
}

// Handler is the entrypoint for Participant requests.
type Handler struct {
	repo Repository
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Println("Participant handler invoked...")
	defer log.Println("Participant handler finished")

	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
	}

	if _, err := fmt.Fprint(res, fmt.Sprintf("Thanks for the %s request! The payload was: %s", req.Method, string(payload))); err != nil {
		log.Printf("Error writing response: %v\n", err)
	}

}
