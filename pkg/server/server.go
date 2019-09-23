package server

import (
	"encoding/json"
	"errors"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"github.com/rejlersembriq/hooked/pkg/router"
	"go.uber.org/zap"
	"net/http"
)

// Server handles incomming http requests.
type Server struct {
	router          *router.Router
	participantRepo participant.Repository
}

// New returns a new Server with routes initialized.
func New(r *router.Router, pr participant.Repository) *Server {
	srvr := &Server{
		router:          r,
		participantRepo: pr,
	}

	srvr.routes()

	return srvr
}

func (s *Server) routes() {
	s.router.GET("/participants", s.participantsGET())
	s.router.POST("/participant", s.participantPOST())
	s.router.PUT("/participant/:id", s.participantPUT())
	s.router.GET("/participant/:id", s.participantGET())
	s.router.DELETE("/participant/:id", s.participantDELETE())
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(res, req)
}

func (s *Server) participantsGET() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ps, err := s.participantRepo.GetAll()
		if err != nil {
			zap.L().Error("Error retrieveing resources.", zap.String("error", err.Error()))
			http.Error(res, "Error retrieving resources", http.StatusInternalServerError)
			return
		}

		sendJSON(&ps).ServeHTTP(res, req)
	}
}

func (s *Server) participantPOST() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var p participant.Participant
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			zap.L().Error("Error unmarshalling request.", zap.String("error", err.Error()))
			http.Error(res, "Error unmarshalling request", http.StatusInternalServerError)
			return
		}

		p.ID = nil
		saved, err := s.participantRepo.Save(p)
		if err != nil {
			zap.L().Error("Error persisting resource.", zap.String("error", err.Error()))
			http.Error(res, "Error persisting resource", http.StatusInternalServerError)
			return
		}

		sendJSON(&saved).ServeHTTP(res, req)
	}
}

func (s *Server) participantPUT() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id, exists := router.GetParam(req.Context(), "id")
		if !exists {
			http.Error(res, "Unable to get request parameter", http.StatusInternalServerError)
			return
		}

		var p participant.Participant
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			zap.L().Error("Error unmarshalling request.", zap.String("error", err.Error()))
			http.Error(res, "Error unmarshalling request", http.StatusInternalServerError)
			return
		}

		p.ID = &id
		saved, err := s.participantRepo.Save(p)
		if err != nil {
			if errors.Is(err, participant.ErrNotExist) {
				http.Error(res, "Resource not found", http.StatusNotFound)
				return
			}

			zap.L().Error("Error persisting resource.", zap.String("error", err.Error()))
			http.Error(res, "Error persisting resource", http.StatusInternalServerError)
			return
		}

		sendJSON(&saved).ServeHTTP(res, req)
	}
}

func (s *Server) participantGET() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id, exists := router.GetParam(req.Context(), "id")
		if !exists {
			http.Error(res, "Unable to get request parameter", http.StatusInternalServerError)
			return
		}

		p, err := s.participantRepo.Get(id)
		if err != nil {
			if errors.Is(err, participant.ErrNotExist) {
				http.Error(res, "Resource not found", http.StatusNotFound)
				return
			}

			zap.L().Error("Error persisting resource.", zap.String("error", err.Error()))
			http.Error(res, "Error retrieving resource", http.StatusInternalServerError)
			return
		}

		sendJSON(&p).ServeHTTP(res, req)
	}
}

func (s *Server) participantDELETE() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id, exists := router.GetParam(req.Context(), "id")
		if !exists {
			http.Error(res, "Unable to get request parameter", http.StatusInternalServerError)
			return
		}

		if err := s.participantRepo.Delete(id); err != nil {
			if errors.Is(err, participant.ErrNotExist) {
				http.Error(res, "Resource not found", http.StatusNotFound)
				return
			}

			zap.L().Error("rror deleting resource.", zap.String("id", id), zap.String("error", err.Error()))
			http.Error(res, "Error retrieving resource", http.StatusInternalServerError)
			return
		}

		sendString("Deleted").ServeHTTP(res, req)
	}
}

func sendJSON(v interface{}) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(v); err != nil {
			zap.L().Error("Error marshalling response.", zap.String("error", err.Error()))
			http.Error(res, "Error marshalling response", http.StatusInternalServerError)
		}
	}
}

func sendString(s string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := res.Write([]byte(s)); err != nil {
			zap.L().Error("Error sending string response.", zap.String("error", err.Error()))
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
