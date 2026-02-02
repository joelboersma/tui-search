package main

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const resultsPerPage = 10
const maxResults = 100

var s SearchService

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func InitSearchService() {
	if fileExists(".env") {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}

	// If Brave is defined, use it
	braveKey := os.Getenv("BRAVE_API_KEY")
	if braveKey != "" {
		s = &BraveSearchService{}
		s.Init()
		return
	}

	// If Google is defined, use it
	googleKey := os.Getenv("GOOGLE_API_KEY")
	googleCx := os.Getenv("GOOGLE_CUSTOM_SEARCH_CONTEXT")
	if googleKey != "" && googleCx != "" {
		s = &GoogleSearchService{}
		s.Init()
		return
	}

	if s == nil {
		log.Fatal("No environment variables set for search. For Brave, must set BRAVE_API_KEY. For Google, must set GOOGLE_API_KEY and GOOGLE_CUSTOM_SEARCH_CONTEXT.")
	}
}

type SearchResult struct {
	Title      string
	URL        string
	DisplayURL string
}

type SearchService interface {
	Init()
	NewSearch(query string) []SearchResult
	NextPage(query string) []SearchResult
	PrevPage(query string) []SearchResult
	HasNextPage() bool
	HasPrevPage() bool
}
