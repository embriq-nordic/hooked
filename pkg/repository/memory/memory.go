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
	if p.ID != "" {
		// Entry should exist. Update.
		pp, exists := m.participants[p.ID]
		if !exists {
			return nil, participant.ErrNotExist
		}

		if p.Name != "" {
			pp.Name = p.Name
		}

		if p.Email != "" {
			pp.Email = p.Email
		}

		if p.Phone != "" {
			pp.Phone = p.Phone
		}

		if p.Org != "" {
			pp.Org = p.Org
		}

		// TODO: This makes it impossible to set a score to 0. But always accepting zero might wipe the score when updating aother field.
		if p.Score != 0 {
			pp.Score = p.Score
		}

		if p.Comment != "" {
			pp.Comment = p.Comment
		}

		pp.Updated = time.Now()

		return pp, nil

	} else {
		// Entry doesn't exist. Insert.
		p.ID = uuid.New().String()
		p.Created = time.Now()
		p.Updated = time.Now()

		m.participants[p.ID] = &p

		return &p, nil
	}
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
