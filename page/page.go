package page

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
)

type Page interface {
	GetTitle() string
	GetLinks() []string
}

type Collector struct {
	page Page
}

func NewCollector(page Page) *Collector{
	return &Collector{page: page}
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
