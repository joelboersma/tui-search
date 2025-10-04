package main

import (
	"github.com/rivo/tview"
)

type Link struct {
	Title string
	Url   string
}

func main() {
	app := tview.NewApplication()

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

	list := tview.NewList()
	for index, link := range links {
		index := index + 1
		if index == 10 {
			index = 0
		}
		shortcut := rune(index + '0')
		list.AddItem(link.Title, link.Url, shortcut, func() {})
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}
}
