package eskeeper

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
)

type esclient struct {
	client  *elasticsearch.Client
	verbose bool
}

func newEsClient(urls []string, user, pass string) (*esclient, error) {
	conf := elasticsearch.Config{
		Addresses: urls,
		Username:  user,
		Password:  pass,
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
