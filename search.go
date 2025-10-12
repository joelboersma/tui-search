package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

const resultsPerPage = 10

var (
	svc *customsearch.Service

	apiKey string // GOOGLE_API_KEY
	cx     string // GOOGLE_CUSTOM_SEARCH_CONTEXT
)

func InitSearchService() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	apiKey = os.Getenv("GOOGLE_API_KEY")
	cx = os.Getenv("GOOGLE_CUSTOM_SEARCH_CONTEXT")

	ctx := context.Background()
	svc, err = customsearch.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
}

func Search(query string, start int64) *customsearch.Search {
	if start+resultsPerPage > 100 {
		log.Fatal("cannot have more than 100 results per query")
	}

	resp, err := svc.Cse.List().Cx(cx).Q(query).Start(start).Num(resultsPerPage).Do()
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
