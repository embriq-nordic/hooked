package server

import (
	"fmt"
	"github.com/rejlersembriq/hooked/pkg/participant"
	"github.com/rejlersembriq/hooked/pkg/router"
	"net/http"
)

// Server handles incomming http requests.
type Server struct {
	router          *router.Router
	participantRepo *participant.Repository
}

// New returns a new Server with routes initialized.
func New(r *router.Router, pr *participant.Repository) *Server {
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
	s.router.GET("/participant/:id", s.participantGET())
	s.router.DELETE("/participant/:id", s.participantDELETE())
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(res, req)
}

func (s *Server) participantsGET() http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "You called GET on /participants")
	})
}

func (s *Server) participantPOST() http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "You called POST on /participant")
	})
}

func (s *Server) participantGET() http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		param, exist := router.GetParam(req.Context(), "id")
		if !exist {
			http.Error(res, "Error getting parameter", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(res, fmt.Sprintf("You called GET on /participant with id: %s", param))
	})
}

func (s *Server) participantDELETE() http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		param, exist := router.GetParam(req.Context(), "id")
		if !exist {
			http.Error(res, "Error getting parameter", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(res, fmt.Sprintf("You called DELETE on /participant with id: %s", param))
	})
}
