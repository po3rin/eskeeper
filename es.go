package eskeeper

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
)

type esclient struct {
	client  *elasticsearch.Client
	verbose bool
}

func newEsClient(urls []string, user, pass string) (*esclient, error) {
	retryBackoff := backoff.NewExponentialBackOff()
	retryBackoff.InitialInterval = time.Second

	conf := elasticsearch.Config{
		Addresses:     urls,
		Username:      user,
		Password:      pass,
		RetryOnStatus: []int{408, 429, 502, 503, 504},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			d := retryBackoff.NextBackOff()
			return d
		},
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &esclient{
		client: es,
	}, nil
}

func (e *esclient) log(msg string) {
	if e.verbose {
		fmt.Println(msg)
	}
}

func (e *esclient) logf(format string, a ...interface{}) {
	if e.verbose {
		fmt.Printf(format, a...)
	}
}
