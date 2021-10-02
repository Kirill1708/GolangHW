package main

import (
  "fmt"
)

type Fetcher interface {
  // Fetch returns the body of URL and
  // a slice of URLs found on that page.
  Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
type result struct{
  url, body string
  urls[]string
  err error
  depth int

}

func Crawl(url string, depth int, fetcher Fetcher) {
  results := make(chan *result)
  fetched := make(map[string]bool)
  fetch := func(url string, depth int){
    body, urls, err := fetcher.Fetch(url)
    results <- &result{url, body, urls, err, depth}
  }
  go fetch(url, depth)
  fetched[url] = true
  for fetching := 1; fetching > 0; fetching--{
    res:= <- results
    if res.err !=nil{
      fmt.Println(res.err)
      continue
    }
    fmt.Printf("found: %s %q\n", res.url, res.body)
    if res.depth > 0 {
      for _, u := range res.urls{
        if !fetched[u]{
          fetching++
          go fetch(u,res.depth-1)
          fetched[u]=true
        }
      }
    }
  }
  close(results)

}

func main() {
  Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
  body string
  urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
  if res, ok := f[url]; ok {
    return res.body, res.urls, nil
  }
  return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
  "https://golang.org/": &fakeResult{
    "The Go Programming Language",
    []string{
      "https://golang.org/pkg/",
      "https://golang.org/cmd/",
    },
  },
  "https://golang.org/pkg/": &fakeResult{
    "Packages",
    []string{
      "https://golang.org/",
      "https://golang.org/cmd/",
      "https://golang.org/pkg/fmt/",
      "https://golang.org/pkg/os/",
    },
  },
  "https://golang.org/pkg/fmt/": &fakeResult{
    "Package fmt",
    []string{
      "https://golang.org/",
      "https://golang.org/pkg/",
    },
  },
  "https://golang.org/pkg/os/": &fakeResult{
    "Package os",
    []string{
      "https://golang.org/",
      "https://golang.org/pkg/",
    },
  },
}