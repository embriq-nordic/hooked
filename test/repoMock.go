package test

import "github.com/rejlersembriq/hooked/pkg/participant"

// RepoMock is used to mock Participant Repository. Inject the desired behaviour.
type RepoMock struct {
	SaveHandler   func(participant participant.Participant) (*participant.Participant, participant.Error)
	GetHandler    func(id string) (*participant.Participant, participant.Error)
	GetAllHandler func() ([]*participant.Participant, participant.Error)
	DeleteHandler func(id string) participant.Error
}

// Save mocks participant.Repository Save.
func (r *RepoMock) Save(participant participant.Participant) (*participant.Participant, participant.Error) {
	return r.SaveHandler(participant)
}

// Get mocks participant.Repository Get.
func (r *RepoMock) Get(id string) (*participant.Participant, participant.Error) {
	return r.GetHandler(id)
}

// GetAll mocks participant.Repository GetAll.
func (r *RepoMock) GetAll() ([]*participant.Participant, participant.Error) {
	return r.GetAllHandler()
}

// Delete mocks participant.Repository Delete.
func (r *RepoMock) Delete(id string) participant.Error {
	return r.DeleteHandler(id)
}
