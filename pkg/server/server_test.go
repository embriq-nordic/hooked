package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"github.com/rejlersembriq/hooked/pkg/router"
	"github.com/rejlersembriq/hooked/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_ServeHTTP_NotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/invalid", nil)
	res := httptest.NewRecorder()

	mock := &test.RepoMock{}

	rtr := router.New()
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_GETParticipants(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/participants", nil)
	res := httptest.NewRecorder()

	mock := &test.RepoMock{
		GetAllHandler: func() ([]*participant.Participant, participant.Error) {
			return []*participant.Participant{}, nil
		},
	}

	rtr := router.New()
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_GETparticipants_Error(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/participants", nil)
	res := httptest.NewRecorder()

	mock := &test.RepoMock{
		GetAllHandler: func() ([]*participant.Participant, participant.Error) {
			return nil, errors.New("SomeError")
		},
	}

	rtr := router.New()
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_POSTParticipant(t *testing.T) {
	p := &participant.Participant{
		ID:      "ignoreId",
		Name:    "Test Testson",
		Email:   "test@testson.com",
		Phone:   "12345678",
		Org:     "TestOrg",
		Score:   2,
		Comment: "Test comment.",
	}

	payloadBytes, _ := json.Marshal(p)

	req, _ := http.NewRequest(http.MethodPost, "/participant/", bytes.NewBuffer(payloadBytes))
	res := httptest.NewRecorder()

	mock := &test.RepoMock{
		SaveHandler: func(p participant.Participant) (*participant.Participant, participant.Error) {
			assert.Equal(t, "", p.ID)

			p.ID = "someId"
			p.Created = time.Now()
			p.Updated = time.Now()
			return &p, nil
		},
	}

	rtr := router.New()
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_POSTParticipant_Error(t *testing.T) {
	p := &participant.Participant{
		ID:      "ignoreId",
		Name:    "Test Testson",
		Email:   "test@testson.com",
		Phone:   "12345678",
		Org:     "TestOrg",
		Score:   2,
		Comment: "Test comment.",
	}

	payloadBytes, _ := json.Marshal(p)

	req, _ := http.NewRequest(http.MethodPost, "/participant", bytes.NewBuffer(payloadBytes))
	res := httptest.NewRecorder()

	mock := &test.RepoMock{
		SaveHandler: func(p participant.Participant) (*participant.Participant, participant.Error) {
			return &p, errors.New("SomeError")
		},
	}

	rtr := router.New()
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_PUTParticipant(t *testing.T) {
	p := &participant.Participant{
		ID:      "ignoreId",
		Name:    "Test Testson",
		Email:   "test@testson.com",
		Phone:   "12345678",
		Org:     "TestOrg",
		Score:   2,
		Comment: "Test comment.",
	}

	payloadBytes, _ := json.Marshal(p)

	req, _ := http.NewRequest(http.MethodPut, "/participant/someId", bytes.NewBuffer(payloadBytes))
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		SaveHandler: func(p participant.Participant) (*participant.Participant, participant.Error) {
			assert.Equal(t, "someId", p.ID)

			p.Created = time.Now()
			p.Updated = time.Now()
			return &p, nil
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_PUTParticipant_Error(t *testing.T) {
	p := &participant.Participant{
		ID:      "ignoreId",
		Name:    "Test Testson",
		Email:   "test@testson.com",
		Phone:   "12345678",
		Org:     "TestOrg",
		Score:   2,
		Comment: "Test comment.",
	}

	payloadBytes, _ := json.Marshal(p)

	req, _ := http.NewRequest(http.MethodPut, "/participant/someId", bytes.NewBuffer(payloadBytes))
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		SaveHandler: func(p participant.Participant) (*participant.Participant, participant.Error) {
			assert.Equal(t, "someId", p.ID)
			return &p, errors.New("SomeError")
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_PUTParticipant_NotFound(t *testing.T) {
	p := &participant.Participant{
		ID:      "ignoreId",
		Name:    "Test Testson",
		Email:   "test@testson.com",
		Phone:   "12345678",
		Org:     "TestOrg",
		Score:   2,
		Comment: "Test comment.",
	}

	payloadBytes, _ := json.Marshal(p)

	req, _ := http.NewRequest(http.MethodPut, "/participant/someId", bytes.NewBuffer(payloadBytes))
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		SaveHandler: func(p participant.Participant) (*participant.Participant, participant.Error) {
			assert.Equal(t, "someId", p.ID)
			return nil, participant.ErrNotExist
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_GETParticipant(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/participant/someId", nil)
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		GetHandler: func(id string) (*participant.Participant, participant.Error) {
			return &participant.Participant{}, nil
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_GETParticipant_Error(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/participant/someId", nil)
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		GetHandler: func(id string) (*participant.Participant, participant.Error) {
			return nil, errors.New("SomeError")
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_GETParticipant_NotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/participant/someId", nil)
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		GetHandler: func(id string) (*participant.Participant, participant.Error) {
			return nil, participant.ErrNotExist
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_DELETEParticipant(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "/participant/someId", nil)
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		DeleteHandler: func(id string) participant.Error {
			return nil
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_DELETEParticipant_Error(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "/participant/someId", nil)
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		DeleteHandler: func(id string) participant.Error {
			return errors.New("SomeError")
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestServer_ServeHTTP_DELETEParticipant_NotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "/participant/someId", nil)
	res := httptest.NewRecorder()

	rtr := router.New()
	mock := &test.RepoMock{
		DeleteHandler: func(id string) participant.Error {
			return participant.ErrNotExist
		},
	}
	srvr := New(rtr, mock)

	srvr.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
}
