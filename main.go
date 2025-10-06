package main

import (
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
			response := Search(query)
			renderResultsView(response)
		}
	})

	inputField.SetBorder(true).SetTitle("Search").SetBorderPadding(1, 1, 3, 3)
	app.SetRoot(inputField, true)
}

func renderResultsView(searchResponse *customsearch.Search) {
	results := searchResponse.Items

	list := tview.NewList()
	for index, result := range results {
		key := (index + 1) % 10 // 0-9
		shortcut := rune(key + '0')
		list.AddItem(result.Title, result.DisplayLink, shortcut, func() {
			OpenURL(result.Link)
		})
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	list.SetBorder(true).SetTitle("Results")
	app.SetRoot(list, true)
}

func main() {
	InitSearchService()

	app = tview.NewApplication()

	renderSearchView()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
