package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

const resultsPerPage = 10
const maxResults = 100

var (
	svc *customsearch.Service

	apiKey string // GOOGLE_API_KEY
	cx     string // GOOGLE_CUSTOM_SEARCH_CONTEXT
)

func InitSearchService() {
	err := godotenv.Load()
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	apiKey = os.Getenv("GOOGLE_API_KEY")
	cx = os.Getenv("GOOGLE_CUSTOM_SEARCH_CONTEXT")

	ctx := context.Background()
	svc, err = customsearch.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
}

func Search(query string, start int64) *customsearch.Search {
	resp, err := svc.Cse.List().Cx(cx).Q(query).Start(start).Num(resultsPerPage).Do()
	if err != nil {
		app.Stop()

		if strings.Contains(err.Error(), "\"reason\": \"RATE_LIMIT_EXCEEDED\"") {
			log.Fatal("Google CustomSearch API Quota Exceeded")
		} else {
			// Unknown error
			log.Fatal(err)
		}
	}

	return resp
}

func HasNextPage(searchResponse *customsearch.Search) bool {
	return len(searchResponse.Queries.NextPage) > 0
}

func HasPrevPage(searchResponse *customsearch.Search) bool {
	if len(searchResponse.Queries.PreviousPage) == 0 {
		return false
	}

	startIndex := searchResponse.Queries.PreviousPage[0].StartIndex
	if resultsPerPage+startIndex > maxResults {
		// Cannot have more than maxResults results per query across all pages.
		return false
	}

	return true
}

func NextPage(query string, searchResponse *customsearch.Search) *customsearch.Search {
	startIndex := searchResponse.Queries.NextPage[0].StartIndex
	return Search(query, startIndex)
}

func PrevPage(query string, searchResponse *customsearch.Search) *customsearch.Search {
	startIndex := searchResponse.Queries.PreviousPage[0].StartIndex
	return Search(query, startIndex)
}
