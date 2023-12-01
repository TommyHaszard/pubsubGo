package main

import (
	"fmt"
	"sync"
)

type QueueProducerFunc func() Queue

type Queue interface {
	Push([]byte) (int, error)
	Grab(int) ([]byte, error)
}

type ByteQueue struct {
	name string
	mu   sync.RWMutex
	data [][]byte
}

func NewByteQueue(name string) *ByteQueue {
	return &ByteQueue{
		name: name,
		data: make([][]byte, 0),
	}
}

func (s *ByteQueue) Push(b []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, b)
	return len(s.data) - 1, nil
}

func (s *ByteQueue) Grab(offset int) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if 0 > offset {
		return nil, fmt.Errorf("offset cannot be smaller then 0")
	}
	if len(s.data) < offset {
		return nil, fmt.Errorf("offset (%d) too high", offset)
	}
	return s.data[offset], nil
}
