package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/customsearch/v1"
)

var (
	app *tview.Application
)

func renderSearchView() {
	inputField := tview.NewInputField()
	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			app.Stop()
		case tcell.KeyEnter:
			query := inputField.GetText()
			if query == "" {
				return
			}

			// hide cursor
			app.SetFocus(nil)

			pageNumber := 1
			showLoading(pageNumber)

			go func() {
				response := NewSearch(query)
				renderResultsView(response, query, 1)
			}()
		}
	})

	inputField.SetBorder(true).SetTitle("Search").SetBorderPadding(1, 1, 3, 3)
	app.SetRoot(inputField, true)
}

func renderResultsView(searchResponse *customsearch.Search, query string, pageNumber int) {
	results := searchResponse.Items

	// Results
	list := tview.NewList()
	for index, result := range results {
		key := (index + 1) % resultsPerPage
		shortcut := rune(key + '0')
		list.AddItem(result.Title, result.DisplayLink, shortcut, func() {
			OpenURL(result.Link)
		})
	}

	// Next page
	if HasNextPage(searchResponse) {
		list.AddItem("Next", "Next page of results", 'n', func() {
			newPage := pageNumber + 1
			showLoading(newPage)

			go func() {
				response := NextPage(query, searchResponse)
				renderResultsView(response, query, newPage)
			}()
		})
	}

	// Previous page
	if HasPrevPage(searchResponse) {
		list.AddItem("Previous", "Previous page of results", 'b', func() {
			newPage := pageNumber - 1
			showLoading(newPage)

			go func() {
				response := PrevPage(query, searchResponse)
				renderResultsView(response, query, newPage)
			}()
		})
	}

	// New query
	list.AddItem("New Search", "Start a new search", 's', renderSearchView)

	// Quit
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	title := fmt.Sprint("Results - Page ", pageNumber)
	list.SetBorder(true).SetTitle(title)
	app.SetRoot(list, true)
}

func showLoading(pageNumber int) {
	textView := tview.NewTextView().
		SetText("Loading...").
		SetTextAlign(tview.AlignCenter)

	title := fmt.Sprint("Results - Page ", pageNumber)

	// Vertically center the column flex in a row flex
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(textView, 3, 1, false).
		AddItem(nil, 0, 1, false).
		SetTitle(title).
		SetBorder(true)

	app.SetRoot(flex, true)
}

func main() {
	InitSearchService()

	app = tview.NewApplication()

	renderSearchView()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
