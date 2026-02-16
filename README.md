# TUI Search

A terminal UI for Search Engine APIs.

## Setup

1. Clone the project

```sh
git clone git@github.com:joelboersma/tui-search.git
```

2. Set up one of the available search engines

3. Build and run the project
```sh
go build && ./tui-search
```

## Search Engines

### Brave Search API (Recommended)

The Brave Search API has simple setup, allows searching the entire internet, and can be used for free.

#### Setup
1. Log into the [Brave Search API Dashboard](https://api-dashboard.search.brave.com/app/dashboard), or create an account if you don't have one already.
2. Create an API key for TUI Search.
3. Set the `BRAVE_API_KEY` environment variable (either in a `.env` file or in the shell) to your API key.

### Google Programmable Search

Google Programmable Search can be an effective tool for searching sites that you specify when creating the engine. However, Google is phasing out generally-available support for full web search, and as a result all new engines must be configured using the "Sites to search" feature ([source](https://programmablesearchengine.googleblog.com/2026/01/updates-to-our-web-search-products.html)).

If you want to search the full web and want simple setup, use Brave Search instead.

#### Setup
1. [Create a programmable search engine](https://developers.google.com/custom-search/docs/tutorial/creatingcse) and a related Google Cloud project.
2. Define the following environment variables, (either in a `.env` file or in the shell):
    - `GOOGLE_API_KEY`: The Google Cloud project's API key
    - `GOOGLE_CUSTOM_SEARCH_CONTEXT`: The search engine's ID (context)
