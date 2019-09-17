package main

import (
	"github.com/rejlersembriq/hooked/pkg/participant"
	"log"
	"net/http"
	"time"
)

func main() {

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      &participant.Handler{},
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on %s\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error serving http: %v", err)
	}
}
