package crawler

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/hiago-balbino/web-crawler/internal/core/pager"
	"github.com/hiago-balbino/web-crawler/internal/pkg/logger"
	"github.com/hiago-balbino/web-crawler/internal/pkg/metrics"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/html"
)

var log = logger.GetLogger()

const (
	linkTag    = "a"
	hrefProp   = "href"
	patternURI = `((http|https):\/\/)`
)

type CrawlerService struct {
	pagerService pager.PagerUsecase
	database     CrawlerDatabase
}

func NewCrawlerService(pagerService pager.PagerUsecase, database CrawlerDatabase) CrawlerService {
	return CrawlerService{pagerService: pagerService, database: database}
}

func (p CrawlerService) Craw(ctx context.Context, uri string, depth uint) ([]string, error) {
	start := time.Now().UTC()
	defer func() {
		metrics.DeltaTimeToProcessLinks.Observe(time.Since(start).Seconds())
	}()

	if links, err := p.database.Find(ctx, uri, depth); err == nil && len(links) > 0 {
		log.Info("returning data from database")

		return links, nil
	}

	links := make([]string, 0)
	ch := make(chan *linkAddress)
	fetched := sync.Map{}

	fetch := func(wg *sync.WaitGroup, uri string) {
		defer wg.Done()

		node, err := p.pagerService.GetNode(uri)
		uris := extractAddresses([]string{}, node)

		ch <- &linkAddress{uri, uris, err}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go fetch(&wg, uri)
	fetched.Store(uri, true)

	for fetching := uint(1); fetching <= depth; fetching++ {
		linkAddress := <-ch
		if linkAddress.err != nil {
			log.Error("error to get uri node", zap.Field{Type: zapcore.StringType, String: linkAddress.err.Error()})
			metrics.LinksErrorCounter.Inc()

			return nil, linkAddress.err
		}

		if len(linkAddress.uris) == 0 {
			break
		}

		for _, uri := range linkAddress.uris {
			metrics.LinksCounter.Inc()

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
		log.Error("error inserting data into database", zap.Field{Type: zapcore.StringType, String: err.Error()})
	}

	return links, nil
}

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
