package memory

import (
	"github.com/google/uuid"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"time"
)

// Memory implements repository persisting participants in memory.
type Memory struct {
	participants map[string]*participant.Participant
}

// New returns a new memory repository.
func New() *Memory {
	return &Memory{
		participants: make(map[string]*participant.Participant),
	}
}

// Save persists a participant to memory.
func (m *Memory) Save(p participant.Participant) (*participant.Participant, participant.Error) {
	if p.ID != nil {
		// Entry should exist. Update.
		pp, exists := m.participants[*p.ID]
		if !exists {
			return nil, participant.ErrNotExist
		}

		if p.Name != nil {
			pp.Name = p.Name
		}

		if p.Email != nil {
			pp.Email = p.Email
		}

		if p.Phone != nil {
			pp.Phone = p.Phone
		}

		if p.Org != nil {
			pp.Org = p.Org
		}

		if p.Score != nil {
			pp.Score = p.Score
		}

		if p.Comment != nil {
			pp.Comment = p.Comment
		}

		now := time.Now()
		pp.Updated = &now

		return pp, nil

	}

	// Entry doesn't exist. Insert.
	id := uuid.New().String()
	p.ID = &id

	now := time.Now()
	p.Created = &now
	p.Updated = &now

	m.participants[*p.ID] = &p

	return &p, nil
}

// Get retrieves a participant from memory.
func (m *Memory) Get(id string) (*participant.Participant, participant.Error) {
	p, exists := m.participants[id]
	if !exists {
		return nil, participant.ErrNotExist
	}
	return p, nil
}

// GetAll retrieves all participants from memory.
func (m *Memory) GetAll() ([]*participant.Participant, participant.Error) {
	var ps []*participant.Participant

	for _, v := range m.participants {
		ps = append(ps, v)
	}

	return ps, nil
}

// Delete removes and entry matching the provided id.
func (m *Memory) Delete(id string) participant.Error {
	if _, exists := m.participants[id]; !exists {
		return participant.ErrNotExist
	}

	delete(m.participants, id)

	return nil
}
