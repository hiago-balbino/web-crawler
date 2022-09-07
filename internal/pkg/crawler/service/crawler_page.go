package crawler

import (
	"context"
	"log"
	"regexp"
	"sync"

	"github.com/hiago-balbino/web-crawler/internal/core/crawler"
	"github.com/hiago-balbino/web-crawler/internal/core/pager"
	"golang.org/x/net/html"
)

const (
	linkTag    = "a"
	hrefProp   = "href"
	patternURI = `((http|https):\/\/)`
)

// CrawlerPage is a implementation to handle with web crawler.
type CrawlerPage struct {
	provider pager.PagerService
	database crawler.CrawlerDatabase
}

// NewCrawlerPage is a constructor to create a new instance of CrawlerPage.
func NewCrawlerPage(pager pager.PagerService, database crawler.CrawlerDatabase) CrawlerPage {
	return CrawlerPage{provider: pager, database: database}
}

// Craw execute the call to craw pages concurrently and will respect depth param.
func (p CrawlerPage) Craw(ctx context.Context, uri string, depth int) ([]string, error) {
	if links, err := p.database.Find(ctx, uri, depth); err == nil && len(links) > 0 {
		return links, nil
	}

	links := make([]string, 0)
	ch := make(chan *dataResult)
	fetched := sync.Map{}

	fetch := func(wg *sync.WaitGroup, uri string) {
		defer wg.Done()

		node, err := p.provider.GetNode(uri)
		uris := extractAddresses([]string{}, node)

		ch <- &dataResult{uri, uris, err}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go fetch(&wg, uri)
	fetched.Store(uri, true)

	for fetching := 1; fetching <= depth; fetching++ {
		result := <-ch
		if result.err != nil {
			return nil, result.err
		}

		if len(result.uris) == 0 {
			break
		}

		for _, uri := range result.uris {
			if _, found := fetched.Load(uri); !found {
				wg.Add(1)

				go fetch(&wg, uri)
				fetched.Store(uri, true)
				links = append(links, uri)
			}
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	if err := p.database.Insert(ctx, uri, depth, links); err != nil {
		log.Default().Println(err.Error())
	}

	return links, nil
}

// extractAddresses recursively extracts the addresses of the HTML node.
func extractAddresses(links []string, node *html.Node) []string {
	if node == nil {
		return links
	}

	if node.Type == html.ElementNode && node.Data == linkTag {
		compile := regexp.MustCompile(patternURI)

		for _, attr := range node.Attr {
			if attr.Key == hrefProp && compile.Match([]byte(attr.Val)) {
				links = append(links, attr.Val)
			}
		}
	}

	for next := node.FirstChild; next != nil; next = next.NextSibling {
		links = extractAddresses(links, next)
	}

	return links
}
