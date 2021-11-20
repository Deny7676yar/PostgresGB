package main

import (
	"context"
	_ "github.com/golang/mock/mockgen/model"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	//Обьявляем поле с информациеей о стартее
	log.WithFields(log.Fields{
		"Start crawler": time.Now(),
	}).Info()

	cfg := Config{
		MaxDepth:   3,
		MaxResults: 10,
		MaxErrors:  5,
		Url:        "https://telegram.org",
		Timeout:    10,
	}
	var cr Crawler
	var r Requester
	wg := &sync.WaitGroup{}
	wg.Add(1)

	r = NewRequester(time.Duration(cfg.Timeout) * time.Second)
	log.WithFields(log.Fields{
		"New request": r,
	}).Debug()

	cr = NewCrawler(r)
	log.WithFields(log.Fields{
		"New Crawler": cr,
	}).Debug()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Timeout))//общий таймаут в секундах

	go cr.Scan(ctx, wg, cfg.Url, cfg.MaxDepth) //Запускаем краулер в отдельной рутине

	go ProcessResult(ctx, cancel, cr, cfg) //Обрабатываем результаты в отдельной рутине

	crawler := SearchDepthCrawler(cfg.MaxDepth)
	go SigDepth(ctx, crawler, 2)

	sigCh := make(chan os.Signal)        //Создаем канал для приема сигналов
	signal.Notify(sigCh, syscall.SIGINT) //Подписываемся на сигнал SIGINT

	for {
		select {
		case <-ctx.Done(): //Если всё завершили - выходим
			return
		case <-sigCh:
			log.WithFields(log.Fields{
				"SIGINT": <-sigCh,
			}).Info("cencel context")
			cancel() //Если пришёл сигнал SigInt - завершаем контекст
		}
	}

	wg.Wait()

	log.WithFields(log.Fields{
		"Wait ": "cancel crawler",
	}).Info()

}