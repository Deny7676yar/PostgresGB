package main

import (
	"context"
	"github.com/Deny7676yar/Go_level2/GB_BP/requester"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type CrawlResult struct {
	Err   error
	Title string
	Url   string
}




//Crawler - интерфейс (контракт) краулера
type Crawler interface {
	Scan(ctx context.Context, wg *sync.WaitGroup, url string, depth int)
	ChanResult() <-chan CrawlResult
}



type crawler struct {
	r       requester.Requester
	res     chan CrawlResult
	visited map[string]struct{}
	mu           sync.RWMutex
	searchDepth     int
}

func NewCrawler(r requester.Requester) *crawler {
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
