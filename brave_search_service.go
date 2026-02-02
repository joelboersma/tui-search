package main

import (
	"log"
	"os"
)

type BraveSearchService struct {
	apiKey string // BRAVE_API_KEY
}

func (s *BraveSearchService) Init() {
	s.apiKey = os.Getenv("BRAVE_API_KEY")
	if s.apiKey == "" {
		log.Fatal("Must define environment variable BRAVE_API_KEY")
	}
}
