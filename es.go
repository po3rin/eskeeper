package eskeeper

import (
	"context"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

type esclient struct {
	client *elasticsearch.Client
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

func (c *esclient) createIndex(ctx context.Context, conf Config) error {
	i := c.client.Indices
	for indexName, mapping := range conf.Index {
		_, err := i.Create(
			indexName,
			i.Create.WithContext(ctx),
			i.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *esclient) syncAlias(ctx context.Context, conf Config) error {
	return nil
}
