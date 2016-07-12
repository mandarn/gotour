package crawl

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return

	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return

	}

	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		Crawl(u, depth-1, fetcher)

	}

}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func CrawlUtil(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	var mut sync.Mutex
	var urlCache = make(map[string]bool)
	var wg sync.WaitGroup
	var CrawlerUtil func(url string, depth int, fetcher Fetcher)

	CrawlerUtil = func(url string, depth int, fetcher Fetcher) {
		defer wg.Done()

		if depth <= 0 {
			return

		}

		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return

		}

		mut.Lock()
		urlCache[url] = true
		mut.Unlock()

		fmt.Printf("found: %s %q\n", url, body)
		for _, u := range urls {
			mut.Lock()
			_, ok := urlCache[u]
			mut.Unlock()

			if !ok {
				wg.Add(1)
				go CrawlerUtil(u, depth-1, fetcher)

			}

		}

	}

	wg.Add(1)
	go CrawlerUtil(url, depth, fetcher)
	wg.Wait()

}

func TestCrawl() {
	fmt.Println("Normal crawl function")
	Crawl("http://golang.org/", 4, fetcher)
	fmt.Println("CrawlUtil function")
	CrawlUtil("http://golang.org/", 4, fetcher)

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
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
