package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
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

type GoogleSearchService struct {
	svc    *customsearch.Service
	apiKey string // GOOGLE_API_KEY
	cx     string // GOOGLE_CUSTOM_SEARCH_CONTEXT

	curResponse *customsearch.Search
}

func (s *GoogleSearchService) Init() {
	s.apiKey = os.Getenv("GOOGLE_API_KEY")
	s.cx = os.Getenv("GOOGLE_CUSTOM_SEARCH_CONTEXT")

	if s.apiKey == "" || s.cx == "" {
		log.Fatal("Must define environment variables GOOGLE_API_KEY and GOOGLE_CUSTOM_SEARCH_CONTEXT")
	}

	var err error
	ctx := context.Background()
	s.svc, err = customsearch.NewService(ctx, option.WithAPIKey(s.apiKey))
	if err != nil {
		log.Fatal(err)
	}
}

func (s *GoogleSearchService) search(query string, start int64) {
	resp, err := s.svc.Cse.List().Cx(s.cx).Q(query).Start(start).Num(resultsPerPage).Do()
	if err != nil {
		app.Stop()

		if strings.Contains(err.Error(), "\"reason\": \"RATE_LIMIT_EXCEEDED\"") {
			log.Fatal("Google CustomSearch API Quota Exceeded")
		} else {
			// Unknown error
			log.Fatal(err)
		}
	}

	s.curResponse = resp
}

func (s *GoogleSearchService) NewSearch(query string) []SearchResult {
	s.search(query, 0)
	return s.toSearchResults()
}

func (s *GoogleSearchService) NextPage(query string) []SearchResult {
	startIndex := s.curResponse.Queries.NextPage[0].StartIndex
	s.search(query, startIndex)
	return s.toSearchResults()
}

func (s *GoogleSearchService) PrevPage(query string) []SearchResult {
	startIndex := s.curResponse.Queries.PreviousPage[0].StartIndex
	s.search(query, startIndex)
	return s.toSearchResults()
}

func (s *GoogleSearchService) HasNextPage() bool {
	return len(s.curResponse.Queries.NextPage) > 0
}

func (s *GoogleSearchService) HasPrevPage() bool {
	if s == nil || len(s.curResponse.Queries.PreviousPage) == 0 {
		return false
	}

	startIndex := s.curResponse.Queries.PreviousPage[0].StartIndex
	if resultsPerPage+startIndex > maxResults {
		// Cannot have more than maxResults results per query across all pages.
		return false
	}

	return true
}

func (s *GoogleSearchService) toSearchResults() []SearchResult {
	results := []SearchResult{}
	for _, item := range s.curResponse.Items {
		result := SearchResult{
			item.Title, item.Link, item.DisplayLink,
		}
		results = append(results, result)
	}
	return results
}
