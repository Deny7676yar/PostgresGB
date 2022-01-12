package requester

import (
	"context"
	"github.com/Deny7676yar/Go_level2/GB_BP/page"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Requester interface {
	Get(ctx context.Context, url string) (page.Page, error)
}

type requester struct {
	timeout time.Duration
}

func NewRequester(timeout time.Duration) requester {
	return requester{timeout: timeout}
}

func (r requester) Get(ctx context.Context, url string) (page.Page, error) {
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

		page, err := page.NewPage(body.Body)
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

