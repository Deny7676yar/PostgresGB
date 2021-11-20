package main

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type CrawlResult struct {
	Err   error
	Title string
	Url   string
}

type Page interface {
	GetTitle() string
	GetLinks() []string
}

type page struct {
	doc *goquery.Document
}

func NewPage(raw io.Reader) (Page, error) {
	doc, err := goquery.NewDocumentFromReader(raw)
	if err != nil {
		log.Error(err)
	}
	return &page{doc: doc}, nil
}

func (p *page) GetTitle() string {

	str := p.doc.Find("title").First().Text()

	log.WithFields(log.Fields{
		"GetTitle": str,
	}).Info("get title")

	return str
}

func (p *page) GetLinks() []string {
	var urls []string
	p.doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			urls = append(urls, url)
			log.WithFields(log.Fields{
				"Get body:": urls,
			}).Info("all urls")
		}
	})
	return urls
}

type Requester interface {
	Get(ctx context.Context, url string) (Page, error)
}

type requester struct {
	timeout time.Duration
}

func NewRequester(timeout time.Duration) requester {
	return requester{timeout: timeout}
}

func (r requester) Get(ctx context.Context, url string) (Page, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		cl := &http.Client{
			Timeout: r.timeout,
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		body, err := cl.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"Get body:": err,
			}).Errorf("Do Not body")
			return nil, err
		}
		defer body.Body.Close()

		page, err := NewPage(body.Body)
		if err != nil {
			log.WithFields(log.Fields{
				"NewPage:": err,
			}).Panicf("No Page")
			return nil, err
		}
		return page, nil
	}
	return nil, nil
}

//Crawler - интерфейс (контракт) краулера
type Crawler interface {
	Scan(ctx context.Context, wg *sync.WaitGroup, url string, depth int)
	ChanResult() <-chan CrawlResult
}



type crawler struct {
	r       Requester
	res     chan CrawlResult
	visited map[string]struct{}
	mu           sync.RWMutex
	searchDepth     int
}

func NewCrawler(r Requester) *crawler {
	return &crawler{
		r:       r,
		res:     make(chan CrawlResult),
		visited: make(map[string]struct{}),
		mu:      sync.RWMutex{},
	}
}

func (c *crawler) Scan(ctx context.Context, wg *sync.WaitGroup, url string, depth int) {
	//if depth <= 0 { //Проверяем то, что есть запас по глубине
	//	return
	//}

	defer wg.Done()
	if depth >= c.searchDepth {
		return
	}
	c.mu.RLock()
	_, ok := c.visited[url] //Проверяем, что мы ещё не смотрели эту страницу
	c.mu.RUnlock()
	if ok {
		return
	}
	select {
	case <-ctx.Done(): //Если контекст завершен - прекращаем выполнение
		return
	default:
		page, err := c.r.Get(ctx, url) //Запрашиваем страницу через Requester
		if err != nil {
			c.res <- CrawlResult{Err: err} //Записываем ошибку в канал
			return
		}
		c.mu.Lock()
		c.visited[url] = struct{}{} //Помечаем страницу просмотренной
		c.mu.Unlock()
		c.res <- CrawlResult{ //Отправляем результаты в канал
			Title: page.GetTitle(),
			Url:   url,
		}
		for _, link := range page.GetLinks() {
			go c.Scan(ctx, wg, link, depth-1) //На все полученные ссылки запускаем новую рутину сборки
		}
	}
}

func (c *crawler) ChanResult() <-chan CrawlResult {
	return c.res
}

//Config - структура для конфигурации
type Config struct {
	MaxDepth   int
	MaxResults int
	MaxErrors  int
	Url        string
	Timeout    int //in seconds
}

func init(){
	//Используем JSON формат для вывода
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}



func SearchDepthCrawler(maxDepth int) *crawler {
	return &crawler{
		searchDepth: maxDepth,
	}
}

func SigDepth(ctx context.Context, c *crawler, d int) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)
	for {
		select {
		case <-ctx.Done():
			return
		case <-sigChan:
			log.WithFields(log.Fields{
				"SIGUSR1": <-sigChan,
			}).Info("Depth += 2")
			c.InitDepth(d)
		}
	}
}

func (c *crawler) InitDepth(dep int) {
	c.mu.Lock()
	c.searchDepth += dep
	c.mu.Unlock()
}

func ProcessResult(ctx context.Context, cancel func(), cr Crawler, cfg Config){
	var maxResult, maxErrors = cfg.MaxResults, cfg.MaxErrors
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-cr.ChanResult():
			if msg.Err != nil {
				maxErrors--
				log.Printf("crawler result return err: %s\n", msg.Err.Error())
				if maxErrors <= 0 {
					cancel()
					return
				}
			} else {
				maxResult--
				log.Printf("crawler result: [url: %s] Title: %s\n", msg.Url, msg.Title)
				if maxResult <= 0 {
					cancel()
					return
				}
			}
		}
	}
}
