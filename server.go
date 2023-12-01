package main

import (
	"fmt"
	"net/http"
)

type Message struct {
	Topic string
	Data  []byte
}

type Config struct {
	ListenAddr        string
	QueueProducerFunc QueueProducerFunc
}

type Server struct {
	*Config
	topics map[string]Queue
}

func NewServer(cfg *Config) (*Server, error) {
	return &Server{Config: cfg,
		topics: make(map[string]Queue)}, nil
}

func (s *Server) Start() {
	http.ListenAndServe(s.ListenAddr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}

func (s *Server) CreateTopic(name string) error {
	if _, exists := s.topics[name]; !exists {
		s.topics[name] = s.QueueProducerFunc()
		return nil
	}
	return fmt.Errorf("Error: Created Topic with name " + name + " already exists.")
}
