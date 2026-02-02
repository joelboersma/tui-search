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

	// TODO: select between Google and other providers
	s = &GoogleSearchService{}
	s.Init()
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
