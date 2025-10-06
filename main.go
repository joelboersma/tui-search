package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Link struct {
	Title string
	Url   string
}

func (l *Link) open() {
	if err := OpenURL(l.Url); err != nil {
		panic(err)
	}
}

func renderSearchView(app *tview.Application, onSubmit func()) {
	inputField := tview.NewInputField()
	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			onSubmit()
		case tcell.KeyEscape:
			app.Stop()
		}
	})

	inputField.SetBorder(true).SetTitle("Search")
	app.SetRoot(inputField, true)
}

func renderResultsView(app *tview.Application, links *[]Link) {
	var firstTenLinks []Link
	if len(*links) > 10 {
		firstTenLinks = (*links)[:10]
	} else {
		firstTenLinks = *links
	}

	list := tview.NewList()
	for index, link := range firstTenLinks {
		key := (index + 1) % 10 // 0-9
		shortcut := rune(key + '0')
		list.AddItem(link.Title, link.Url, shortcut, func() { link.open() })
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	list.SetBorder(true).SetTitle("Results")
	app.SetRoot(list, true)
}

func main() {
	// sample data
	links := []Link{
		{
			"A (not so) short laptop recommendation guide - 2025 ...",
			"https://www.reddit.com/r/gamedev/comments/1hr463f/a_not_so_short_laptop_recommendation_guide_2025/",
		},
		{
			"The Best Laptops We've Tested (October 2025)",
			"https://www.pcmag.com/picks/the-best-laptops",
		},
		{
			"The best laptops in 2025 based on our testing and reviews",
			"https://www.laptopmag.com/reviews/best-laptops-1",
		},
	}

	app := tview.NewApplication()

	onSearchSubmit := func() { renderResultsView(app, &links) }
	renderSearchView(app, onSearchSubmit)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
