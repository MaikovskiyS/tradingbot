package server

import (
	"log"
	"net/http"
	"time"
	"trrader/internal/traidingview"
)

type server struct {
	server *http.Server
}

func New(tv *traidingview.TraidingView) *server {
	return &server{

		server: &http.Server{
			Addr:         ":8080",
			Handler:      tv.Router,
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
		},
	}
}
func (s *server) Run() {
	// go func() {
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}
	//	}()
}
