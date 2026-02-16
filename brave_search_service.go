package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BraveSearchResponse struct {
	Web BraveSearchWebResponse `json:"web"`
}
type BraveSearchWebResponse struct {
	Results []BraveSearchWebResult `json:"results"`
}
type BraveSearchWebResult struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type BraveSearchService struct {
	apiKey  string // BRAVE_API_KEY
	curPage int
}

func (s *BraveSearchService) Init() {
	s.apiKey = os.Getenv("BRAVE_API_KEY")
	if s.apiKey == "" {
		log.Fatal("Must define environment variable BRAVE_API_KEY")
	}
}

func (s *BraveSearchService) search(query string) []BraveSearchWebResult {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.search.brave.com/res/v1/web/search",
		nil,
	)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Subscription-Token", s.apiKey)

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("count", strconv.Itoa(resultsPerPage))
	q.Add("offset", strconv.Itoa(s.curPage-1))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		app.Stop()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf(
				"Response failed with status code: %d.\nBody parsing failed with error: %s",
				resp.StatusCode,
				err,
			)
		}
		log.Fatalf("Response failed with...\nStatus code: %d\nBody: %s\n", resp.StatusCode, body)
	}

	var parsedResponse BraveSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResponse)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	return parsedResponse.Web.Results
}

func (s *BraveSearchService) NewSearch(query string) []SearchResult {
	s.curPage = 1
	braveResults := s.search(query)
	return s.toSearchResults(braveResults)
}

func (s *BraveSearchService) NextPage(query string) []SearchResult {
	s.curPage += 1
	braveResults := s.search(query)
	return s.toSearchResults(braveResults)
}

func (s *BraveSearchService) PrevPage(query string) []SearchResult {
	s.curPage -= 1
	braveResults := s.search(query)
	return s.toSearchResults(braveResults)
}

func (s *BraveSearchService) HasNextPage() bool {
	return resultsPerPage*s.curPage < maxResults
}

func (s *BraveSearchService) HasPrevPage() bool {
	return s.curPage > 1
}

func (s *BraveSearchService) toSearchResults(braveResults []BraveSearchWebResult) []SearchResult {
	results := []SearchResult{}
	for _, item := range braveResults {
		displayUrl := item.Url
		displayUrl = strings.TrimPrefix(displayUrl, "http://")
		displayUrl = strings.TrimPrefix(displayUrl, "https://")
		displayUrl = strings.Split(displayUrl, "/")[0]

		result := SearchResult{
			item.Title, item.Url, displayUrl,
		}
		results = append(results, result)
	}
	return results
}

// curl "https://api.search.brave.com/res/v1/web/search?q=brave+search" \
//   -H "Accept: application/json" \
//   -H "X-Subscription-Token: <YOUR_API_KEY>"
