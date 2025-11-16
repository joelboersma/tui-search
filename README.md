# TUI Search

A terminal UI for Google's Programmable Search Engine.

## Setup

1. Clone the project

```sh
git clone git@github.com:joelboersma/tui-search.git
```

2. [Create a programmable search engine](https://developers.google.com/custom-search/docs/tutorial/creatingcse) and a related Google Cloud project.
3. Define the following environment variables, (either in a `.env` or in the shell):
    - `GOOGLE_API_KEY`: The Google Cloud project's API key
    - `GOOGLE_CUSTOM_SEARCH_CONTEXT`: The search engine's ID (context)

4. Build and run the project
```sh
go build && ./tui-search
```
