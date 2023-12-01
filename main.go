package main

import (
	"log"
)

func QueueFunc(url string) QueueProducerFunc {
	return func() Queue {
		return NewByteQueue(url)
	}
}

func main() {
	cfg := &Config{
		ListenAddr:        ":3000",
		QueueProducerFunc: QueueFunc("Name"),
	}
	s, err := NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}
